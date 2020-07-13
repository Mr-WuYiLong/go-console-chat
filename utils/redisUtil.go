package utils

import (
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

// 初始化redisPool
func init() {
	redisPool = RedisPool()
}

// RedisPool redis的池
func RedisPool() *redis.Pool {
	return &redis.Pool{
		MaxActive:   0,
		IdleTimeout: 300,
		MaxIdle:     8,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

// GetRedis 获得redis连接
func GetRedis() redis.Conn {
	return redisPool.Get()
}
