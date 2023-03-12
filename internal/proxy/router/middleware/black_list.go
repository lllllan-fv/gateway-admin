package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"github.com/lllllan-fv/gateway-admin/public/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func HTTPBlackListMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		whileIpList := []string{}
		if serviceDetail.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.WhiteList, ",")
		}

		blackIpList := []string{}
		if serviceDetail.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.BlackList, ",")
		}

		if serviceDetail.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if utils.InStringSlice(blackIpList, c.ClientIP()) {
				resp.Error(c, 3001, fmt.Errorf("%s in black ip list", c.ClientIP()))
				return
			}
		}

		c.Next()
	}
}

func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		whileIpList := []string{}
		if serviceDetail.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.WhiteList, ",")
		}

		blackIpList := []string{}
		if serviceDetail.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.BlackList, ",")
		}

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}
		if serviceDetail.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if utils.InStringSlice(blackIpList, clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func GrpcBlackListMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		whileIpList := []string{}
		if serviceDetail.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.WhiteList, ",")
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		blackIpList := []string{}

		if serviceDetail.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.BlackList, ",")
		}

		if serviceDetail.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if utils.InStringSlice(blackIpList, clientIP) {
				return fmt.Errorf("%s in black ip list", clientIP)
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
