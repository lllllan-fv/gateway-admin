package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"              // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	_ "github.com/GoAdminGroup/themes/adminlte"                   // ui theme
	"github.com/lllllan-fv/gateway-admin/internal/admin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy"
)

func main() {
	admin.Run()
	proxy.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	admin.Stop()
	proxy.Stop()
}
