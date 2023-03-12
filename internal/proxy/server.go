package proxy

import "github.com/lllllan-fv/gateway-admin/internal/proxy/server"

func Run() {
	server.HttpServerRun()
	server.HttpsServerRun()
	server.TcpServerRun()
	server.GrpcServerRun()
}

func Stop() {
	server.HttpServerStop()
	server.HttpsServerStop()
	server.TcpServerStop()
	server.GrpcServerStop()
}
