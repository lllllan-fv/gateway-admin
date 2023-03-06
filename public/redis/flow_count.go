package redis

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/lllllan-fv/gateway-admin/conf"
	"github.com/lllllan-fv/gateway-admin/public/consts"
)

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) //获取数据
			atomic.StoreInt64(&reqCounter.TickerCount, 0)            //重置数据

			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)
			rdb.IncrBy(context.Background(), dayKey, tickerCount)
			rdb.Expire(context.Background(), dayKey, 86400*2)
			rdb.IncrBy(context.Background(), hourKey, tickerCount)
			rdb.Expire(context.Background(), hourKey, 86400*2)

			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("reqCounter.GetDayData err", err)
				continue
			}
			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}
			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (o *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(conf.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", consts.RedisFlowDayKey, dayStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(conf.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", consts.RedisFlowHourKey, hourStr, o.AppID)
}

func (o *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return rdb.Get(context.Background(), o.GetHourKey(t)).Int64()
}

func (o *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return rdb.Get(context.Background(), o.GetDayKey(t)).Int64()
}

//原子增加
func (o *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
