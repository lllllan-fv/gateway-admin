package middleware

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

func HTTPFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		if serviceDetail.ServiceFlowLimit != 0 {
			if serviceLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName,
				float64(serviceDetail.ServiceFlowLimit),
			); err != nil {
				resp.Error(c, 5001, err)
				return
			} else if !serviceLimiter.Allow() {
				resp.Error(c, 5002, fmt.Errorf("service flow limit %v", serviceDetail.ServiceFlowLimit))
				return
			}
		}

		if serviceDetail.ClientIPFlowLimit > 0 {
			if clientLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName+"_"+c.ClientIP(),
				float64(serviceDetail.ClientIPFlowLimit),
			); err != nil {
				resp.Error(c, 5003, err)
				return
			} else if !clientLimiter.Allow() {
				resp.Error(c, 5002, fmt.Errorf("%v flow limit %v", c.ClientIP(), serviceDetail.ClientIPFlowLimit))
				return
			}
		}

		fmt.Println("flow limit next")
		c.Next()
	}
}

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*models.GatewayApp)

		if appInfo.QPS > 0 {
			if clientLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.QPS),
			); err != nil {
				resp.Error(c, 5001, err)
				return
			} else if !clientLimiter.Allow() {
				resp.Error(c, 5002, fmt.Errorf("%v flow limit %v", c.ClientIP(), appInfo.QPS))
				return
			}
		}

		c.Next()
	}
}
