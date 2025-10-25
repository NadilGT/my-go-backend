package utils

import (
	"sync"
	"time"
)

// CacheItem holds cached data with expiration
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// SimpleCache is a thread-safe in-memory cache
type SimpleCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// NewSimpleCache creates a new cache instance
func NewSimpleCache() *SimpleCache {
	cache := &SimpleCache{
		items: make(map[string]*CacheItem),
	}
	// Start cleanup goroutine
	go cache.startCleanup()
	return cache
}

// Set stores a value in cache with TTL (time to live)
func (c *SimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value from cache
// Returns (value, true) if found and not expired, (nil, false) otherwise
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}

// Delete removes a key from cache
func (c *SimpleCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from cache
func (c *SimpleCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// startCleanup runs a background goroutine to remove expired items
func (c *SimpleCache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes expired items from cache
func (c *SimpleCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.ExpiresAt) {
			delete(c.items, key)
		}
	}
}

// Global cache instance for metrics
var MetricsCache = NewSimpleCache()
