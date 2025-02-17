package lock

import (
	"context"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries     = 5
	baseRetryDelay = 100 * time.Millisecond
	maxRetryDelay  = 5 * time.Second
)

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}

// Acquire now includes retry logic with exponential backoff and jitter.
func (r *RedisLock) Acquire(ctx context.Context, key string) (bool, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		ok, err := r.client.SetNX(ctx, key, "locked", 5*time.Minute).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}

		// Compute exponential backoff with jitter
		backoff := baseRetryDelay * (1 << attempt)               // Exponential growth
		jitter := time.Duration(rand.Int63n(int64(backoff / 2))) // Add randomness
		sleepDuration := backoff + jitter

		if sleepDuration > maxRetryDelay {
			sleepDuration = maxRetryDelay
		}

		time.Sleep(sleepDuration)
	}
	return false, nil
}

func (r *RedisLock) Release(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
