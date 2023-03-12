package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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

func TCPFlowLimitMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		if serviceDetail.ServiceFlowLimit != 0 {
			serviceLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName,
				float64(serviceDetail.ServiceFlowLimit),
			)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("service flow limit %v", serviceDetail.ServiceFlowLimit)))
				c.Abort()
				return
			}
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		if serviceDetail.ClientIPFlowLimit > 0 {
			clientLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName+"_"+clientIP,
				float64(serviceDetail.ClientIPFlowLimit),
			)
			if err != nil {
				c.conn.Write([]byte(err.Error()))
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				c.conn.Write([]byte(fmt.Sprintf("%v flow limit %v", clientIP, serviceDetail.ClientIPFlowLimit)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func GrpcFlowLimitMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, h grpc.StreamHandler) error {
		if serviceDetail.ServiceFlowLimit != 0 {
			serviceLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName,
				float64(serviceDetail.ServiceFlowLimit),
			)
			if err != nil {
				return err
			}
			if !serviceLimiter.Allow() {
				return fmt.Errorf("service flow limit %v", serviceDetail.ServiceFlowLimit)
			}
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if serviceDetail.ClientIPFlowLimit > 0 {
			clientLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowServicePrefix+serviceDetail.ServiceName+"_"+clientIP,
				float64(serviceDetail.ClientIPFlowLimit),
			)
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return fmt.Errorf("%v flow limit %v", clientIP, serviceDetail.ClientIPFlowLimit)
			}
		}

		if err := h(srv, ss); err != nil {
			log.Printf("GrpcFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}

func GrpcJwtFlowLimitMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if appInfo.QPS > 0 {
			clientLimiter, err := handler.GetFlowLimiterHandler().GetLimiter(
				consts.FlowAppPrefix+appInfo.AppID+"_"+clientIP,
				float64(appInfo.QPS),
			)
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return fmt.Errorf("%v flow limit %v", clientIP, appInfo.QPS)
			}
		}

		if err := h(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
