package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Redis interface {
	GetClient() *redis.Client
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, keys ...string) error
	Close() error
	Ping(ctx context.Context) error
}

type redisClient struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, redisURL string, password string) (Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	if password != "" {
		opts.Password = password
	}

	client := redis.NewClient(opts)

	// Instrument all Redis commands with OTel spans
	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, fmt.Errorf("instrument redis tracing: %w", err)
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &redisClient{client: client}, nil
}
