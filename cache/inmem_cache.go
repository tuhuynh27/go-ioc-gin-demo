package cache

import (
	"container/list"
	"context"
	"sync"
)

// LRUCache is an in-memory LRU cache implementation
type LRUCache struct {
	Component  struct{}
	Implements struct{} `implements:"Cache"`
	Qualifier  struct{} `value:"inmem"`

	capacity int
	mu       sync.Mutex
	store    map[string]*list.Element
	ll       *list.List
}

// entry is a key-value pair for the LRU cache
type entry struct {
	key   string
	value interface{}
}

// NewLRUCache creates a new LRUCache with the specified capacity
func NewLRUCache() *LRUCache {
	return &LRUCache{
		capacity: 100,
		store:    make(map[string]*list.Element, 100),
		ll:       list.New(),
	}
}

// Get retrieves a value from the LRU cache
func (c *LRUCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.store[key]; ok {
		c.ll.MoveToFront(elem) // Move accessed item to the front
		return elem.Value.(*entry).value, nil
	}
	return nil, nil
}

// Set stores a value in the LRU cache
func (c *LRUCache) Set(ctx context.Context, key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.store[key]; ok {
		c.ll.MoveToFront(elem) // Move existing item to the front
		elem.Value.(*entry).value = value
		return nil
	}

	// If cache is at capacity, remove the least recently used item
	if c.ll.Len() == c.capacity {
		backElem := c.ll.Back()
		if backElem != nil {
			c.ll.Remove(backElem)
			delete(c.store, backElem.Value.(*entry).key)
		}
	}

	// Add new item to the front of the list
	newEntry := &entry{key: key, value: value}
	newElem := c.ll.PushFront(newEntry)
	c.store[key] = newElem
	return nil
}
