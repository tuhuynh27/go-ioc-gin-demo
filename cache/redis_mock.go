package cache

import (
	"context"
	"sync"
)

// RedisMock is a mock implementation of Cache interface for testing
type RedisMock struct {
	Component  struct{}
	Implements struct{} `implements:"Cache"`
	Qualifier  struct{} `value:"mock"`

	store sync.Map
}

// Get retrieves a value from the mock cache
func (m *RedisMock) Get(ctx context.Context, key string) (interface{}, error) {
	if value, ok := m.store.Load(key); ok {
		return value, nil
	}
	return nil, nil
}

// Set stores a value in the mock cache
func (m *RedisMock) Set(ctx context.Context, key string, value interface{}) error {
	m.store.Store(key, value)
	return nil
}
