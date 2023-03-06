package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/router/middleware"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func InitHttpRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Recovery())

	router.GET("/ping", func(c *gin.Context) { resp.Success(c, "pong") })

	router.Use(
		middleware.HTTPAccessModeMiddleware(),
		middleware.HTTPFlowCountMiddleware(),
	)

	return router
}
