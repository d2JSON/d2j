package caching

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

var _ Cacher = (*redisCache)(nil)

// NewRedis is used to create an instance of redis cache.
func NewRedis(opts ConnectionOptions) *redisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Host,
		Password: opts.Password,
		DB:       opts.Database,
	})

	return &redisCache{client}
}

func (r redisCache) Write(ctx context.Context, options WriteOptions) error {
	err := r.client.Set(ctx, options.Key, options.Value, options.TTL).Err()
	if err != nil {
		return fmt.Errorf("set into redis: %w", err)
	}

	return nil
}

func (r redisCache) Read(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrResultIsNil
		}

		return "", fmt.Errorf("get from redis: %w", err)
	}

	return data, nil
}

func (r redisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("delete from redis: %w", err)
	}

	return nil
}

func (r redisCache) Close() error {
	if r.client != nil {
		err := r.client.Close()
		if err != nil {
			return fmt.Errorf("close redis connection: %w", err)
		}
	}

	return nil
}
