package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lllllan-fv/gateway-admin/conf"
)

var rdb *redis.Client

func GetRedisClient() *redis.Client {
	return rdb
}

func Init() {
	cfg := conf.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       int(cfg.Redis.DB),
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}
}
