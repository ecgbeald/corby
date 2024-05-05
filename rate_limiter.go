package main

import (
	"context"
	"log"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	cli *redis.Client
	ctx context.Context
	*redisrate.Limiter
}

func NewRedisRateLimiter(str string) *RedisRateLimiter {
	cli := redis.NewClient(&redis.Options{
		Addr: str,
	})
	ctx := context.Background()
	pong, err := cli.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Cannot ping to redis:", err)
	}
	log.Printf("Successfully connected to redis: %s", pong)

	return &RedisRateLimiter{cli, ctx, redisrate.NewLimiter(cli)}
}
