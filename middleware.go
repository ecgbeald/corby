package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimitMiddleware struct {
	limiter *RedisRateLimiter
}

func NewRateLimitMiddleware(str string) *RateLimitMiddleware {
	return &RateLimitMiddleware{NewRedisRateLimiter(str)}
}

const RateRequest = "rate_request_%s"

func (r *RateLimitMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ipaddr := c.IP()
		expiration := time.Duration(10 * time.Minute)
		pipe := r.limiter.cli.Pipeline()
		pipeCmd := []interface{}{
			pipe.SetNX(r.limiter.ctx, ipaddr, 0, expiration),
			pipe.Incr(r.limiter.ctx, ipaddr),
			pipe.TTL(r.limiter.ctx, ipaddr),
		}
		_, err := pipe.Exec(r.limiter.ctx)
		if err != nil {
			log.Fatal("pipe not executing", err)
		}
		ipCount := pipeCmd[0].(*redis.BoolCmd)
		ipCountInc := pipeCmd[1].(*redis.IntCmd)
		ttl := pipeCmd[2].(*redis.DurationCmd)
		fmt.Printf("%s\n %s\n %s\n", ipCount, ipCountInc, ttl)
		ctx := c.Context()
		res, _ := r.limiter.Allow(ctx, fmt.Sprintf(RateRequest, "username"), redis_rate.Limit{
			Rate:   10,
			Burst:  10,
			Period: time.Second,
		})
		if res.Allowed <= 0 {
			log.Printf("Reached IP rate limit on: %s on %s", c.IP(), c.Path())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too Many Requests on " + c.Path(),
			})
		}
		return c.Next()
	}
}

func main() {
	myFiberApp := fiber.New()
	myFiberApp.Use(NewRateLimitMiddleware("localhost:6379").Handler())
	myFiberApp.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{"ping": "pong"})
	})
	err := myFiberApp.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
