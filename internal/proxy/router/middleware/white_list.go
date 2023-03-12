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

func TCPWhiteListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		iplist := []string{}
		if serviceDetail.WhiteList != "" {
			iplist = strings.Split(serviceDetail.WhiteList, ",")
		}
		if serviceDetail.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s not in white ip list", clientIP)))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func GrpcWhiteListMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		iplist := []string{}
		if serviceDetail.WhiteList != "" {
			iplist = strings.Split(serviceDetail.WhiteList, ",")
		}

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]

		if serviceDetail.OpenAuth == 1 && len(iplist) > 0 {
			if !utils.InStringSlice(iplist, clientIP) {
				return fmt.Errorf("%s not in white ip list", clientIP)
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
