package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
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
		ctx := c.Context()
		res, _ := r.limiter.Allow(ctx, fmt.Sprintf(RateRequest, "username"), redis_rate.Limit{
			Rate:   10,
			Burst:  10,
			Period: time.Second,
		})
		if res.Allowed <= 0 {
			return c.SendStatus(fiber.StatusTooManyRequests)
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
