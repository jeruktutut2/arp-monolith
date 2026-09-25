package eventbus

import (
	"context"
	"fmt"
	"sync"
)

// HandlerFunc is the subscriber callback
type HandlerFunc func(ctx context.Context, payload interface{}) error

// EventBus interface for publishing and subscribing to domain events
type EventBus interface {
	Publish(ctx context.Context, topic string, payload interface{}) error
	Subscribe(topic string, handler HandlerFunc)
}

type inMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]HandlerFunc
}

// NewInMemoryEventBus creates an in-memory asynchronous event bus
func NewInMemoryEventBus() EventBus {
	return &inMemoryEventBus{
		subscribers: make(map[string][]HandlerFunc),
	}
}

func (b *inMemoryEventBus) Subscribe(topic string, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

func (b *inMemoryEventBus) Publish(ctx context.Context, topic string, payload interface{}) error {
	b.mu.RLock()
	handlers, exists := b.subscribers[topic]
	b.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		return nil
	}

	for _, handler := range handlers {
		// Run handler in background goroutine with detached context copy to decouple execution
		go func(h HandlerFunc) {
			if err := h(ctx, payload); err != nil {
				fmt.Printf("⚠️ EventBus handler error on topic [%s]: %v\n", topic, err)
			}
		}(handler)
	}

	return nil
}
