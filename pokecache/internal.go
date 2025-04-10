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
	entries  map[string]cacheEntry
	interval time.Duration
	mx       *sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mx:      &sync.Mutex{},
	}

	go cache.reapLoop(interval)

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.entries[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if resource, found := c.entries[key]; found {
		return resource.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mx.Lock()
		for key := range c.entries {
			if time.Since(c.entries[key].createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mx.Unlock()
	}
}
