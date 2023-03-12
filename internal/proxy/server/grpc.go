package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/router/middleware"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	load_balance "github.com/lllllan-fv/gateway-admin/internal/proxy/service/load_balance"
)

var grpcServerList = []*warpGrpcServer{}

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	serviceList := dao.ListService(consts.GrpcLoadType)
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *models.GatewayServiceInfo) {
			addr := fmt.Sprintf(":%d", serviceDetail.Port)
			rb, err := handler.GetLoadBalancerHandler().GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf("[INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("[INFO] GrpcListen %v err:%v\n", addr, err)
			}

			grpcHandler := NewGrpcLoadBalanceHandler(rb)
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					middleware.GrpcFlowCountMiddleware(serviceDetail),
					middleware.GrpcFlowLimitMiddleware(serviceDetail),
					middleware.GrpcJwtAuthTokenMiddleware(serviceDetail),
					middleware.GrpcJwtFlowCountMiddleware(serviceDetail),
				// grpc_proxy_middleware.GrpcJwtFlowLimitMiddleware(serviceDetail),
				// grpc_proxy_middleware.GrpcWhiteListMiddleware(serviceDetail),
				// grpc_proxy_middleware.GrpcBlackListMiddleware(serviceDetail),
				// grpc_proxy_middleware.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)

			grpcServerList = append(grpcServerList, &warpGrpcServer{
				Addr:   addr,
				Server: s,
			})

			log.Printf("[INFO] GRPC proxy run %v\n", addr)
			if err := s.Serve(lis); err != nil {
				log.Fatalf("[INFO] GRPC proxy run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}

func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
