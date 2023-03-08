package handler

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
)

var transportorHandler *Transportor

func GetTransportorHandler() *Transportor {
	return transportorHandler
}

type Transportor struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

func init() {
	transportorHandler = &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

func (t *Transportor) GetTrans(service *models.GatewayServiceInfo) (*http.Transport, error) {
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.ServiceName {
			return transItem.Trans, nil
		}
	}

	// todo 优化点5
	if service.UpstreamConnectTimeout == 0 {
		service.UpstreamConnectTimeout = 30
	}
	if service.UpstreamMaxIdle == 0 {
		service.UpstreamMaxIdle = 100
	}
	if service.UpstreamIdleTimeout == 0 {
		service.UpstreamIdleTimeout = 90
	}
	if service.UpstreamHeaderTimeout == 0 {
		service.UpstreamHeaderTimeout = 30
	}
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(service.UpstreamConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          service.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.UpstreamIdleTimeout) * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: time.Duration(service.UpstreamHeaderTimeout) * time.Second,
	}

	// save to map and slice
	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.ServiceName,
	}

	t.TransportSlice = append(t.TransportSlice, transItem)
	t.Locker.Lock()
	defer t.Locker.Unlock()

	t.TransportMap[service.ServiceName] = transItem
	return trans, nil
}
