package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/router"
	tcp_server "github.com/lllllan-fv/gateway-admin/internal/proxy/server/tcp"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/handler"
)

var tcpServerList = []*tcp_server.TcpServer{}

type tcpHandler struct{}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpServerRun() {
	serviceList := dao.ListService(consts.TcpLoadType)
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *models.GatewayServiceInfo) {
			addr := fmt.Sprintf(":%d", serviceDetail.Port)
			_, err := handler.GetLoadBalancerHandler().GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			type serviceKeyType string
			var serviceKey serviceKeyType = "service"
			baseCtx := context.WithValue(context.Background(), serviceKey, serviceDetail)
			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: router.InitTCPRouter(),
				BaseCtx: baseCtx,
			}

			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf("[INFO] TCP proxy run %v\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf("[INFO] TCP proxy run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf("[INFO] TCP proxy stop %v stopped\n", tcpServer.Addr)
	}
}
