package uredis

import (
	"context"
	"errors"
	"git.umu.work/be/goframework/accelerator/cache"
	"github.com/go-redis/redis/v8"
)

type UCacheNameType string

const (
	UCacheNameAICache UCacheNameType = "ai-cache"
)

type RedisError struct {
	error
}

func (e RedisError) Error() string { return e.Error() }

func (e RedisError) IsNil() bool {
	return errors.Is(e, redis.Nil)
}

func GetRedisCmd(ctx context.Context, cacheName UCacheNameType) (redis.Cmdable, *RedisError) {
	redisClient, err := cache.GetRedis(ctx, string(cacheName))
	if err != nil {
		return nil, &RedisError{err}
	}
	return redisClient, nil
}
