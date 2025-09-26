package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

// NewCache creates a new cache with a configurable reap interval
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	}

	go c.reapLoop(interval)
	return c
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get gets an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

// reapLoop removes expired entries from the cache
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			cutoff := time.Now().Add(-interval)
			for key, entry := range c.cache {
				if entry.createdAt.Before(cutoff) {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
