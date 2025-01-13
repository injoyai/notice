package user

import (
	"context"
	"github.com/injoyai/base/maps"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache2 interface {
	Get(key string) (string, error)
	Set(key, vale string, expiration time.Duration) error
	Del(key string) error
}

type _redis struct {
	*redis.Client
}

func (this *_redis) Get(key string) (string, error) {
	return this.Client.Get(context.Background(), key).Result()
}

func (this *_redis) Set(key string, value string, expiration time.Duration) error {
	return this.Client.Set(context.Background(), key, value, expiration).Err()
}

func (this *_redis) Del(key string) error {
	return this.Client.Del(context.Background(), key).Err()
}

type _memory struct {
	*maps.Safe
}

func (this *_memory) Get(key string) (string, error) {
	return this.Safe.GetString(key), nil
}

func (this *_memory) Set(key string, value string, expiration time.Duration) error {
	this.Safe.Set(key, value, expiration)
	return nil
}

func (this *_memory) Del(key string) error {
	this.Safe.Del(key)
	return nil
}
