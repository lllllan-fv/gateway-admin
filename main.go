package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"              // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	_ "github.com/GoAdminGroup/themes/adminlte"                   // ui theme
	"github.com/lllllan-fv/gateway-admin/conf"
	"github.com/lllllan-fv/gateway-admin/internal/admin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy"
	"github.com/lllllan-fv/gateway-admin/public/redis"
)

func main() {
	// init module
	conf.Init()
	redis.Init()
	log.Println("init module")

	// run server
	admin.Run()
	proxy.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// stop server
	admin.Stop()
	proxy.Stop()
}
