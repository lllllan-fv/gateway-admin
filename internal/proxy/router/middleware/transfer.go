package middleware

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(models.GatewayServiceInfo)

		for _, item := range strings.Split(serviceDetail.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}

			if items[0] == "add" || items[0] == "edit" {
				c.Request.Header.Set(items[1], items[2])
			}

			if items[0] == "del" {
				c.Request.Header.Del(items[1])
			}
		}

		c.Next()
	}
}

func GrpcHeaderTransferMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			}
			if items[0] == "del" {
				delete(md, items[1])
			}
		}

		if err := ss.SetHeader(md); err != nil {
			return err
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
