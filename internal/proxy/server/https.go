package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/router"
)

var httpsPort = 8002
var httpsSrvHandler *http.Server

func HttpsServerRun() {
	r := router.InitHttpRouter()

	httpsSrvHandler = &http.Server{
		Addr:           fmt.Sprint(":", httpsPort),
		Handler:        r,
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(10) * time.Second,
		MaxHeaderBytes: 1 << uint(20),
	}
	log.Printf("[INFO] HTTPS proxy run: %s\n", fmt.Sprint(":", httpsPort))

	//todo 以下命令只在编译机有效，如果是交叉编译情况下需要单独设置路径
	//if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err!=http.ErrServerClosed {
	if err := httpsSrvHandler.ListenAndServeTLS("./cert_file/server.crt", "./cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] HTTP proxy %s err:%v\n", fmt.Sprint(":", httpsPort), err)
	}
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] HTTPS proxy stop::%v\n", err)
	}
	log.Printf("[INFO] HTTPS proxy stopped\n")
}
