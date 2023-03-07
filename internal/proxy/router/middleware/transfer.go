package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(models.GatewayServiceInfo)

		for _, item := range strings.Split(serviceDetail.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}

			if items[0] == "add" || items[0] == "edit" {
				c.Request.Header.Set(items[1], items[2])
			}

			if items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}

		c.Next()
	}
}
