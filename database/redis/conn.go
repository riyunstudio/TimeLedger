package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"timeLedger/configs"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	DB0 *redis.Client
}

// 初始化 Redis 連線
func Initialize(env *configs.Env) *Redis {
	var (
		once sync.Once
		db0  *redis.Client
	)

	once.Do(func() {
		db0 = newRedisClient(env.RedisHost, env.RedisPort, env.RedisPass, env.RedisPoolSize, 0) // 預設使用 DB 0
	})

	return &Redis{
		DB0: db0,
	}
}

// 建立指定 DB 的 Redis 連線
func newRedisClient(host, port, pass string, pool, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		PoolSize: pool,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("Redis connection failed: %v", err))
	}
	log.Printf("Redis connected (DB %d)", db)

	return rdb
}
