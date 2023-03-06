package redis

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/lllllan-fv/gateway-admin/conf"
)

var rdb *redis.Client

func GetRedisClient() *redis.Client {
	return rdb
}

func init() {
	cfg := conf.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       int(cfg.Redis.DB),
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

}
