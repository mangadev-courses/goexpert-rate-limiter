package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

const (
	oneSecond = time.Second * 1
)

func NewRedisClient(ctx context.Context) (*Redis, error) {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_ = rdb.FlushDB(ctx).Err()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return &Redis{}, fmt.Errorf("Error connecting to Redis: %w", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return &Redis{
		rdb: rdb,
	}, nil
}

func (r *Redis) IncrementRequestCount(ctx context.Context, key string) (int, error) {
	count, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("Error incrementing key: %w", err)
	}

	if count == 1 {
		err = r.rdb.Expire(ctx, key, time.Second).Err()
		if err != nil {
			return 0, fmt.Errorf("Error setting expiration on key: %w", err)
		}
	}
	return int(count), nil
}

func (r *Redis) FreezeRequestCount(ctx context.Context, key string, seconds int) error {
	err := r.rdb.Expire(ctx, key, time.Second*time.Duration(seconds)).Err()
	if err != nil {
		return fmt.Errorf("Error expiring key: %w", err)
	}
	return nil
}

func (r *Redis) IsFrozen(ctx context.Context, key string) (bool, float64, error) {
	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, fmt.Errorf("error checking TTL for key: %w", err)
	}

	if ttl > oneSecond {
		return true, ttl.Seconds(), nil
	}

	return false, ttl.Seconds(), nil
}
