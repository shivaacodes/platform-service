package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()


type Client struct {
	rdb *redis.Client // field
}

func NewClient(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr : addr,
		Password : "", // set via env in production
		DB: 0,
	})
	return &Client{rdb: rdb}
}

func (c *Client)Get(key string) (string,error) {
	return c.rdb.Get(ctx,key).Result()
}

func (c *Client)Set(key string, value interface{}, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Client)Close() error {
	return c.rdb.Close()
}


