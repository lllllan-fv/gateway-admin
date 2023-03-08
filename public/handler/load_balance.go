package handler

import (
	"fmt"
	"sync"

	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	load_balance "github.com/lllllan-fv/gateway-admin/internal/proxy/service/load_balance"
	"github.com/lllllan-fv/gateway-admin/public/consts"
)

var loadBalancerHandler *LoadBalancer

func GetLoadBalancerHandler() *LoadBalancer {
	return loadBalancerHandler
}

type LoadBalancer struct {
	LoadBanlanceMap   map[string]*LoadBalancerItem
	LoadBanlanceSlice []*LoadBalancerItem
	Locker            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBanlance load_balance.LoadBalance
	ServiceName  string
}

func init() {
	loadBalancerHandler = &LoadBalancer{
		LoadBanlanceMap:   map[string]*LoadBalancerItem{},
		LoadBanlanceSlice: []*LoadBalancerItem{},
		Locker:            sync.RWMutex{},
	}
}

func (lbr *LoadBalancer) GetLoadBalancer(service *models.GatewayServiceInfo) (load_balance.LoadBalance, error) {
	for _, lbrItem := range lbr.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.ServiceName {
			return lbrItem.LoadBanlance, nil
		}
	}

	schema := "http://"
	if service.NeedHTTPS == 1 {
		schema = "https://"
	}
	if service.LoadType == consts.TcpLoadType || service.LoadType == consts.GrpcLoadType {
		schema = ""
	}

	ipList := service.GetIPListByModel()
	weightList := service.GetWeightListByModel()
	ipConf := map[string]string{}
	for ipIndex, ipItem := range ipList {
		ipConf[ipItem] = weightList[ipIndex]
	}

	//fmt.Println("ipConf", ipConf)
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}
	lb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbType(service.RoundType), mConf)

	//save to map and slice
	lbItem := &LoadBalancerItem{
		LoadBanlance: lb,
		ServiceName:  service.ServiceName,
	}
	lbr.LoadBanlanceSlice = append(lbr.LoadBanlanceSlice, lbItem)

	lbr.Locker.Lock()
	defer lbr.Locker.Unlock()
	lbr.LoadBanlanceMap[service.ServiceName] = lbItem
	return lb, nil
}
