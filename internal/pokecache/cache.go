package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]CacheEntry
	mu      sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (e CacheEntry) reborn() CacheEntry {
	return CacheEntry{time.Now(), e.val}
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries: make(map[string]CacheEntry),
		mu:      sync.Mutex{},
	}

	go c.reapLoop(interval)

	return &c
}

func (c *Cache) Add(key string, val []byte) {
	entry := CacheEntry{time.Now(), val}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]

	if !ok {
		return nil, false
	}

	c.entries[key] = e.reborn()
	return e.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		c.mu.Lock()
		for key, val := range c.entries {
			if time.Since(val.createdAt) >= interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}

}
