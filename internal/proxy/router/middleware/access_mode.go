package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/service"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

// Match access mode, based on request information
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := service.HTTPAccessMode(c)
		if err != nil {
			resp.Error(c, 1001, err)
			return
		}

		fmt.Printf("service: %v\n", service)
		c.Set("service", service)
		c.Next()
	}
}
