package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/controller"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/router/middleware"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func InitHttpRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Recovery())

	router.GET("/ping", func(c *gin.Context) { resp.Success(c, "pong") })
	router.POST("/token", controller.Tokens)

	router.Use(
		middleware.HTTPAccessModeMiddleware(),
		middleware.HTTPFlowCountMiddleware(),
		middleware.HTTPFlowLimitMiddleware(),
		middleware.HTTPJwtAuthTokenMiddleware(),
		middleware.HTTPJwtFlowCountMiddleware(),
		middleware.HTTPJwtFlowLimitMiddleware(),
		middleware.HTTPWhiteListMiddleware(),
		middleware.HTTPBlackListMiddleware(),
		middleware.HTTPHeaderTransferMiddleware(),
		middleware.HTTPStripUriMiddleware(),
		middleware.HTTPUrlRewriteMiddleware(),
		middleware.HTTPReverseProxyMiddleware(),
	)

	return router
}
