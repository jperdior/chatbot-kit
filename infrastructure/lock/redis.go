package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{client: client}
}

func (r *RedisLock) Acquire(ctx context.Context, key string) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, "locked", 5*time.Minute).Result()
	return ok, err
}

func (r *RedisLock) Release(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
