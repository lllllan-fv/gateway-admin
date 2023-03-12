package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		//统计项 1 全站 2 服务 3 租户
		totalCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()

		serviceCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowServicePrefix + serviceDetail.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()

		c.Next()
	}
}

func GrpcFlowCountMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		totalCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowTotal)
		if err != nil {
			return err
		}
		totalCounter.Increase()

		serviceCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowServicePrefix + serviceDetail.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := h(srv, ss); err != nil {
			log.Printf("GrpcFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}

func GrpcJwtFlowCountMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		appInfos := md.Get("app")
		if len(appInfos) == 0 {
			if err := h(srv, ss); err != nil {
				log.Printf("RPC failed with error %v\n", err)
				return err
			}
			return nil
		}

		appInfo := &models.GatewayApp{}
		if err := json.Unmarshal([]byte(appInfos[0]), appInfo); err != nil {
			return err
		}

		appCounter, err := handler.GetFlowCounterHandler().GetCounter(consts.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()
		if appInfo.QPD > 0 && appCounter.TotalCount > appInfo.QPD {
			return fmt.Errorf("租户日请求量限流 limit:%v current:%v", appInfo.QPD, appCounter.TotalCount)
		}

		if err := h(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
