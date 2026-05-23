package changelog

import (
	"sync"
	"time"
)

// CacheEntry holds a cached fetch result with an expiry timestamp.
type CacheEntry struct {
	Content   string
	FetchedAt time.Time
}

// Cache is a simple in-memory cache for fetched changelog content.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	TTL     time.Duration
}

// NewCache creates a new Cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &Cache{
		entries: make(map[string]CacheEntry),
		TTL:     ttl,
	}
}

// Get returns the cached content for key if present and not expired.
// The second return value indicates whether the entry was found and valid.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}
	if time.Since(entry.FetchedAt) > c.TTL {
		return "", false
	}
	return entry.Content, true
}

// Set stores content for key with the current timestamp.
func (c *Cache) Set(key, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Content:   content,
		FetchedAt: time.Now(),
	}
}

// Invalidate removes the entry for key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of entries currently in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
