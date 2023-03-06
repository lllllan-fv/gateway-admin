package proxy

import "github.com/lllllan-fv/gateway-admin/internal/proxy/server"

func Run() {
	server.HttpServerRun()
}

func Stop() {
	server.HttpServerStop()
}
