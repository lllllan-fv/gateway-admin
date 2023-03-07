package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/jwt"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

//jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		appMatched := false
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		fmt.Printf("token: %v\n", token)
		if token != "" {
			claims, err := jwt.Decode(token)
			if err != nil {
				resp.Error(c, 2002, err)
				return
			}

			appList := dao.ListApp()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.OpenAuth == 1 && !appMatched {
			resp.Error(c, 2003, errors.New("not match valid app"))
			return
		}

		c.Next()
	}
}
