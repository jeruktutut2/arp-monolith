package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockRedis is a thread-safe in-memory mock implementation of the Redis interface.
type MockRedis struct {
	mu   sync.RWMutex
	data map[string]mockItem

	Client *redis.Client

	// Programmable errors for testing error paths
	SetErr    error
	GetErr    error
	DeleteErr error
	CloseErr  error
	PingErr   error

	// Optional custom function overrides
	SetFn    func(ctx context.Context, key string, value []byte, expiration time.Duration) error
	GetFn    func(ctx context.Context, key string) ([]byte, error)
	DeleteFn func(ctx context.Context, keys ...string) error
	PingFn   func(ctx context.Context) error
	CloseFn  func() error
}

type mockItem struct {
	value     []byte
	expiresAt time.Time
}

var _ Redis = (*MockRedis)(nil)

// NewMockRedis initializes a new in-memory MockRedis instance.
func NewMockRedis() *MockRedis {
	return &MockRedis{
		data: make(map[string]mockItem),
	}
}

// GetClient returns the mock *redis.Client (nil by default unless set).
func (m *MockRedis) GetClient() *redis.Client {
	return m.Client
}

// Set stores a key-value byte slice with TTL expiration in memory.
func (m *MockRedis) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if m.SetFn != nil {
		return m.SetFn(ctx, key, value, expiration)
	}
	if m.SetErr != nil {
		return m.SetErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var expTime time.Time
	if expiration > 0 {
		expTime = time.Now().Add(expiration)
	}

	valCopy := make([]byte, len(value))
	copy(valCopy, value)

	m.data[key] = mockItem{
		value:     valCopy,
		expiresAt: expTime,
	}
	return nil
}

// Get retrieves the byte slice by key. Returns nil, nil on cache miss or expiration.
func (m *MockRedis) Get(ctx context.Context, key string) ([]byte, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, key)
	}
	if m.GetErr != nil {
		return nil, m.GetErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.data[key]
	if !exists {
		return nil, nil
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return nil, nil
	}

	valCopy := make([]byte, len(item.value))
	copy(valCopy, item.value)
	return valCopy, nil
}

// Delete removes one or more keys from the in-memory cache.
func (m *MockRedis) Delete(ctx context.Context, keys ...string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, keys...)
	}
	if m.DeleteErr != nil {
		return m.DeleteErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

// Close closes the mock and clears memory.
func (m *MockRedis) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	if m.CloseErr != nil {
		return m.CloseErr
	}
	m.Flush()
	return nil
}

// Ping verifies the mock is responsive.
func (m *MockRedis) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}
	if m.PingErr != nil {
		return m.PingErr
	}
	return nil
}

// Flush removes all keys stored in the mock cache.
func (m *MockRedis) Flush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]mockItem)
}
