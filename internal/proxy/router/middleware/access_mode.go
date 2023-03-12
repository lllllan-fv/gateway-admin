package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

// Match access mode, based on request information
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := HTTPAccessMode(c)
		if err != nil {
			resp.Error(c, 1001, err)
			return
		}

		fmt.Printf("service: %v\n", service)
		c.Set("service", service)
		c.Next()
	}
}

func HTTPAccessMode(c *gin.Context) (*models.GatewayServiceInfo, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path

	for _, service := range dao.ListService(consts.HttpLoadType) {
		if service.RuleType == consts.DomainHTTPRuleType {
			if service.Rule == host {
				return service, nil
			}
		} else {
			if strings.HasPrefix(path, service.Rule) {
				return service, nil
			}
		}
	}

	return nil, errors.New("not matched service")
}
