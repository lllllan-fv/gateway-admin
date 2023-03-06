package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/router"
)

var httpPort = 8001
var httpServer *http.Server

func HttpServerRun() {
	r := router.InitHttpRouter()

	httpServer = &http.Server{
		Addr:           fmt.Sprint(":", httpPort),
		Handler:        r,
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(10) * time.Second,
		MaxHeaderBytes: 1 << uint(20),
	}
	log.Printf("[INFO] HTTP proxy run: %s\n", fmt.Sprint(":", httpPort))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] HTTP proxy %s err:%v\n", fmt.Sprint(":", httpPort), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] HTTP proxy stop:%v\n", err)
	}
	log.Printf("[INFO] HTTP proxy stop\n")
}
