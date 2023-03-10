package router

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/router/middleware"
	tcp_server "github.com/lllllan-fv/gateway-admin/internal/proxy/server/tcp"
	load_balance "github.com/lllllan-fv/gateway-admin/internal/proxy/service/load_balance"
)

func InitTCPRouter(rb load_balance.LoadBalance) *middleware.TcpSliceRouterHandler {
	// 构建路由及设置中间件
	router := middleware.NewTcpSliceRouter()
	router.Group("/").Use(
		middleware.TCPFlowCountMiddleware(),
		middleware.TCPFlowLimitMiddleware(),
		middleware.TCPWhiteListMiddleware(),
		middleware.TCPBlackListMiddleware(),
	)

	// 构建回调 handler
	return middleware.NewTcpSliceRouterHandler(
		func(c *middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
			return NewTcpLoadBalanceReverseProxy(c, rb)
		},
		router,
	)
}

func NewTcpLoadBalanceReverseProxy(c *middleware.TcpSliceRouterContext, lb load_balance.LoadBalance) *TcpReverseProxy {
	return func() *TcpReverseProxy {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		return &TcpReverseProxy{
			ctx:             c.Ctx,
			Addr:            nextAddr,
			KeepAlivePeriod: time.Second,
			DialTimeout:     time.Second,
		}
	}()
}

//TCP反向代理
type TcpReverseProxy struct {
	ctx                  context.Context //单次请求单独设置
	Addr                 string
	KeepAlivePeriod      time.Duration //设置
	DialTimeout          time.Duration //设置超时时间
	DialContext          func(ctx context.Context, network, address string) (net.Conn, error)
	OnDialError          func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
}

func (dp *TcpReverseProxy) dialTimeout() time.Duration {
	if dp.DialTimeout > 0 {
		return dp.DialTimeout
	}
	return 10 * time.Second
}

func (dp *TcpReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if dp.DialContext != nil {
		return dp.DialContext
	}
	return (&net.Dialer{
		Timeout:   dp.DialTimeout,     //连接超时
		KeepAlive: dp.KeepAlivePeriod, //设置连接的检测时长
	}).DialContext
}

func (dp *TcpReverseProxy) keepAlivePeriod() time.Duration {
	if dp.KeepAlivePeriod != 0 {
		return dp.KeepAlivePeriod
	}
	return time.Minute
}

//传入上游 conn，在这里完成下游连接与数据交换
func (dp *TcpReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	//设置连接超时
	var cancel context.CancelFunc
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, dp.dialTimeout())
	}
	dst, err := dp.dialContext()(ctx, "tcp", dp.Addr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		dp.onDialError()(src, err)
		return
	}

	defer func() { go dst.Close() }() //记得退出下游连接

	//设置dst的 keepAlive 参数,在数据请求之前
	if ka := dp.keepAlivePeriod(); ka > 0 {
		if c, ok := dst.(*net.TCPConn); ok {
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(ka)
		}
	}
	errc := make(chan error, 1)
	go dp.proxyCopy(errc, src, dst)
	go dp.proxyCopy(errc, dst, src)
	<-errc
}

func (dp *TcpReverseProxy) onDialError() func(src net.Conn, dstDialErr error) {
	if dp.OnDialError != nil {
		return dp.OnDialError
	}
	return func(src net.Conn, dstDialErr error) {
		log.Printf("tcpproxy: for incoming conn %v, error dialing %q: %v", src.RemoteAddr().String(), dp.Addr, dstDialErr)
		src.Close()
	}
}

func (dp *TcpReverseProxy) proxyCopy(errc chan<- error, dst, src net.Conn) {
	_, err := io.Copy(dst, src)
	errc <- err
}
