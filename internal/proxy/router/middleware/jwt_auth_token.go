package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/jwt"
	"github.com/lllllan-fv/gateway-admin/public/resp"
	"github.com/lllllan-fv/gateway-admin/public/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			resp.Error(c, 2001, errors.New("service not found"))
			return
		}
		serviceDetail := serverInterface.(*models.GatewayServiceInfo)

		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		appMatched := false
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		fmt.Printf("token: %v\n", token)
		if token != "" {
			claims, err := jwt.Decode(token)
			if err != nil {
				resp.Error(c, 2002, err)
				return
			}

			appList := dao.ListApp()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					c.Set("app", appInfo)
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.OpenAuth == 1 && !appMatched {
			resp.Error(c, 2003, errors.New("not match valid app"))
			return
		}

		c.Next()
	}
}

func GrpcJwtAuthTokenMiddleware(serviceDetail *models.GatewayServiceInfo) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		authToken := ""
		auths := md.Get("authorization")
		if len(auths) > 0 {
			authToken = auths[0]
		}

		token := strings.ReplaceAll(authToken, "Bearer ", "")
		appMatched := false
		if token != "" {
			claims, err := jwt.Decode(token)
			if err != nil {
				return err
			}

			appList := dao.ListApp()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					md.Set("app", utils.Obj2Json(appInfo))
					appMatched = true
					break
				}
			}
		}

		if serviceDetail.OpenAuth == 1 && !appMatched {
			return errors.New("not match valid app")
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcJwtAuthTokenMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
