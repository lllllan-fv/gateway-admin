package router

import (
	"github.com/lllllan-fv/gateway-admin/internal/proxy/router/middleware"
	tcp_server "github.com/lllllan-fv/gateway-admin/internal/proxy/server/tcp"
)

func InitTCPRouter() *middleware.TcpSliceRouterHandler {
	// 构建路由及设置中间件
	router := middleware.NewTcpSliceRouter()
	router.Group("/").Use(
		middleware.TCPFlowCountMiddleware(),
		middleware.TCPFlowLimitMiddleware(),
		middleware.TCPWhiteListMiddleware(),
	// tcp_proxy_middleware.TCPBlackListMiddleware(),
	)

	// 构建回调 handler
	return middleware.NewTcpSliceRouterHandler(
		func(c *middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
			// return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
			return nil
		},
		router,
	)
}
