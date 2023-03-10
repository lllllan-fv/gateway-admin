package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return tc, nil
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tcp_proxy context value " + k.name
}

type conn struct {
	server     *TcpServer
	cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) close() {
	c.rwc.Close()
}

func (c *conn) serve(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("tcp: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		c.close()
	}()

	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	if c.server.Handler == nil {
		panic("handler empty")
	}
	c.server.Handler.ServeTCP(ctx, c.rwc)
}

var (
	ErrServerClosed     = errors.New("tcp: Server closed")
	ErrAbortHandler     = errors.New("tcp: abort TCPHandler")
	ServerContextKey    = &contextKey{"tcp-server"}
	LocalAddrContextKey = &contextKey{"local-addr"}
)

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() {
	oc.closeErr = oc.Listener.Close()
}

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

type TcpServer struct {
	Addr    string
	Handler TCPHandler
	err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.Mutex
	inShutdown int32
	doneChan   chan struct{}
	l          *onceCloseListener
}

func (s *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

func (srv *TcpServer) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	if srv.doneChan == nil {
		srv.doneChan = make(chan struct{})
	}
	addr := srv.Addr
	if addr == "" {
		return errors.New("need addr")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (srv *TcpServer) Close() error {
	atomic.StoreInt32(&srv.inShutdown, 1)
	close(srv.doneChan) // 关闭channel
	srv.l.Close()       // 执行listener关闭
	return nil
}

func (srv *TcpServer) Serve(l net.Listener) error {
	srv.l = &onceCloseListener{Listener: l}
	defer srv.l.Close() // 执行listener关闭
	if srv.BaseCtx == nil {
		srv.BaseCtx = context.Background()
	}
	baseCtx := srv.BaseCtx
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			fmt.Printf("accept fail, err: %v\n", e)
			continue
		}
		c := srv.newConn(rw)
		go c.serve(ctx)
	}
	return nil
}

func (srv *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	// 设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}
	return c
}

func (s *TcpServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func ListenAndServe(addr string, handler TCPHandler) error {
	server := &TcpServer{Addr: addr, Handler: handler, doneChan: make(chan struct{})}
	return server.ListenAndServe()
}

var tcpServerList = []*TcpServer{}

type tcpHandler struct{}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpServerRun() {
	serviceList := dao.ListService(consts.TcpLoadType)
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *models.GatewayServiceInfo) {
			addr := fmt.Sprintf(":%d", serviceDetail.Port)
			_, err := handler.GetLoadBalancerHandler().GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			// 构建路由及设置中间件
			// router := tcp_proxy_middleware.NewTcpSliceRouter()
			// router.Group("/").Use(
			// 	tcp_proxy_middleware.TCPFlowCountMiddleware(),
			// 	tcp_proxy_middleware.TCPFlowLimitMiddleware(),
			// 	tcp_proxy_middleware.TCPWhiteListMiddleware(),
			// 	tcp_proxy_middleware.TCPBlackListMiddleware(),
			// )

			// 构建回调 handler
			// routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
			// 	func(c *tcp_proxy_middleware.TcpSliceRouterContext) TCPHandler {
			// 		return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
			// 	},
			// 	router,
			// )

			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &TcpServer{
				Addr:    addr,
				Handler: nil,
				BaseCtx: baseCtx,
			}

			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf("[INFO] TCP proxy run %v\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != ErrServerClosed {
				log.Fatalf("[INFO] TCP proxy run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf("[INFO] TCP proxy stop %v stopped\n", tcpServer.Addr)
	}
}
