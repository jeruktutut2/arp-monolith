package cache

import (
	"context"
	"time"

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
