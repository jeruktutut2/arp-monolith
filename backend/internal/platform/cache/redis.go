package cache

import (
	"context"
	"errors"
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

// GetClient returns the underlying *redis.Client.
func (r *redisClient) GetClient() *redis.Client {
	return r.client
}

// Set stores a byte slice in Redis with an expiration duration.
func (r *redisClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis set %s: %w", key, err)
	}
	return nil
}

// Get retrieves a byte slice from Redis by key.
// If the key does not exist, it returns nil, nil (cache miss).
func (r *redisClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get %s: %w", key, err)
	}
	return data, nil
}

// Delete removes one or more keys from Redis.
func (r *redisClient) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete: %w", err)
	}
	return nil
}

// Close closes the Redis client connection.
func (r *redisClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// Ping verifies the Redis connection is alive.
func (r *redisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
