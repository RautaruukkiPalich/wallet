package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	URI string
}

type Cache struct {
	client *redis.Client
}

func New(ctx context.Context, cfg Config) (*Cache, error) {
	opts, err := redis.ParseURL(cfg.URI)
	if err != nil {
		return nil, err
	}
	c := redis.NewClient(opts)

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{client: c}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return c.SetWithTTL(ctx, key, value, 0)
}

func (c *Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
