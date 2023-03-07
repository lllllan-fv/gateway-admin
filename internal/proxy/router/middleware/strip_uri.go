package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func HTTPStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		if serviceDetail.RuleType == consts.HTTPRuleTypePrefixURL && serviceDetail.NeedStripURI == 1 {
			//fmt.Println("c.Request.URL.Path",c.Request.URL.Path)
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.Rule, "", 1)
			//fmt.Println("c.Request.URL.Path",c.Request.URL.Path)
		}
		//http://127.0.0.1:8080/test_http_string/abbb
		//http://127.0.0.1:2004/abbb

		c.Next()
	}
}
