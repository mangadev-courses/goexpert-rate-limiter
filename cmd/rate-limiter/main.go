package main

import (
	"context"
	"fmt"

	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/limiter"
	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/redis"
)

func main() {
	ctx := context.Background()

	redisClient, err := redis.NewRedisClient(ctx)
	if err != nil {
		panic(err)
	}

	limiter := limiter.New(redisClient)
	res, err := limiter.Allow(ctx, "123", "456")
	if err != nil {
		panic(err)
	}

	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
}
