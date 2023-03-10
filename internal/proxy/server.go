package proxy

import "github.com/lllllan-fv/gateway-admin/internal/proxy/server"

func Run() {
	server.HttpServerRun()
	server.HttpsServerRun()
	server.TcpServerRun()
}

func Stop() {
	server.HttpServerStop()
	server.HttpsServerStop()
	server.TcpServerStop()
}
