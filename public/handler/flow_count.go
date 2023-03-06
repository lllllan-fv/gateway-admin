package handler

import (
	"sync"
	"time"

	"github.com/lllllan-fv/gateway-admin/public/redis"
)

var flowCounterHandler *FlowCounter

func GetFlowCounterHandler() *FlowCounter {
	return flowCounterHandler
}

type FlowCounter struct {
	RedisFlowCountMap   map[string]*redis.RedisFlowCountService
	RedisFlowCountSlice []*redis.RedisFlowCountService
	Locker              sync.RWMutex
}

func init() {
	flowCounterHandler = &FlowCounter{
		RedisFlowCountMap:   map[string]*redis.RedisFlowCountService{},
		RedisFlowCountSlice: []*redis.RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func (counter *FlowCounter) GetCounter(serverName string) (*redis.RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCountSlice {
		if item.AppID == serverName {
			return item, nil
		}
	}

	newCounter := redis.NewRedisFlowCountService(serverName, 1*time.Second)
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serverName] = newCounter
	return newCounter, nil
}
