package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	Entries		map[string]cacheEntry 
	Interval 	time.Duration
	Mu			sync.Mutex
}

type cacheEntry struct {
	createdAt	time.Time
	val			[]byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Interval: interval,
		Entries: make(map[string]cacheEntry),
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	cEntry := cacheEntry{
		createdAt: time.Now(),
		val: value,
	}
	c.Entries[key] = cEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	cEntry, ok := c.Entries[key]
	if !ok {
		var b []byte
		return b, false
	}
	return cEntry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()
	for range ticker.C {
		c.Mu.Lock()
		for key, value := range c.Entries {
			if time.Now().Sub(value.createdAt) > c.Interval {
				delete(c.Entries, key)
			}	
		}
		c.Mu.Unlock()
	}
}

