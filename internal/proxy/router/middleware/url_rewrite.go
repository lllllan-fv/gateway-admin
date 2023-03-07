package middleware

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		for _, item := range strings.Split(serviceDetail.URLRewrite, ",") {
			//fmt.Println("item rewrite",item)
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}

			regexp, err := regexp.Compile(items[0])
			if err != nil {
				//fmt.Println("regexp.Compile err",err)
				continue
			}

			//fmt.Println("before rewrite",c.Request.URL.Path)
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			//fmt.Println("after rewrite",c.Request.URL.Path)
		}

		c.Next()
	}
}
