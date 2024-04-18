package db

import (
	"github.com/redis/go-redis/v9"
	"trpc.group/trpc-go/trpc-database/goredis"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/config"
)

var (
	RedisClient redis.UniversalClient
)

func InitRedisClient() redis.UniversalClient {
	// 扩展接口
	client, err := goredis.New(config.RedisServiceName)
	if err != nil {
		log.Errorf("Redis initialization failed, err=[%v].", err)
	}
	RedisClient = client
	log.Info("Redis initialization completed.")
	return RedisClient
}
