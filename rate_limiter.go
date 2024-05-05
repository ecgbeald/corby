package main

import (
	"context"
	"log"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	*redisrate.Limiter
}

func NewRedisRateLimiter(str string) *RedisRateLimiter {
	cli := redis.NewClient(&redis.Options{
		Addr: str,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Cannot ping to redis:", err)
	}
	return &RedisRateLimiter{redisrate.NewLimiter(cli)}
}
