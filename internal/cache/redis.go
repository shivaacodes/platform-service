package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Cache is the interface for the cache.
type Cache interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Close() error
	Ping(ctx context.Context) error
}

type client struct {
	rdb *redis.Client // field
}

func NewClient(addr, password string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // set via env in production
		DB:       0,
	})
	return &client{rdb: rdb}
}

func (c *client) Get(key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *client) Set(key string, value interface{}, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}
