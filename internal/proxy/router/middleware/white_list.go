package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"github.com/lllllan-fv/gateway-admin/public/utils"
)

// 匹配接入方式 基于请求信息
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		iplist := []string{}
		if serviceDetail.WhiteList != "" {
			iplist = strings.Split(serviceDetail.WhiteList, ",")
		}

		if serviceDetail.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, c.ClientIP()) {
				resp.Error(c, 3001, fmt.Errorf("%s not in white ip list", c.ClientIP()))
				return
			}
		}

		c.Next()
	}
}
