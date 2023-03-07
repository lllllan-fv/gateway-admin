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

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serviceInterface.(*models.GatewayServiceInfo)

		// 统计项 1 全站 2 服务 3 租户
		h := handler.GetFlowCounterHandler()
		totalCounter, err := h.GetCounter(consts.FlowTotal)
		if err != nil {
			resp.Error(c, 4001, err)
			return
		}
		totalCounter.Increase()

		// dayCount, _ := totalCounter.GetDayData(time.Now())
		// fmt.Printf("totalCounter qps:%v,dayCount:%v", totalCounter.QPS, dayCount)
		serviceCounter, err := h.GetCounter(consts.FlowServicePrefix + serviceDetail.ServiceName)
		if err != nil {
			resp.Error(c, 4001, err)
			return
		}
		serviceCounter.Increase()

		// dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		// fmt.Printf("serviceCounter qps:%v,dayCount:%v", serviceCounter.QPS, dayServiceCount)
		c.Next()
	}
}

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*models.GatewayApp)

		appCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			resp.Error(c, 2002, err)
			return
		}
		appCounter.Increase()

		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			resp.Error(c, 2003, fmt.Errorf("租户日请求量限流 limit:%v current:%v", appInfo.QPD, appCounter.TotalCount))
			return
		}

		c.Next()
	}
}
