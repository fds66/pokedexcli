package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntry map[string]CacheEntry
	Mu         sync.Mutex
	Interval   time.Duration
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(time time.Duration) *Cache {
	newCacheEntry := make(map[string]CacheEntry)
	newCache := Cache{
		CacheEntry: newCacheEntry,
		Interval:   time,
	}
	go newCache.reapLoop()
	return &newCache
}

func (c *Cache) Add(key string, value []byte) bool {
	//fmt.Println("Adding new data to the cache")
	c.Mu.Lock()
	defer c.Mu.Unlock()
	newEntry := CacheEntry{
		Val:       value,
		CreatedAt: time.Now(),
	}

	c.CacheEntry[key] = newEntry
	//fmt.Printf("Added %s to the cache", key)
	return true

}

func (c *Cache) Get(key string) ([]byte, bool) {
	//fmt.Println("Retrieving data from the cache")
	entry, exists := c.CacheEntry[key]
	if exists {
		return entry.Val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	// each time interval passes it will remove entries older than interval
	ticker := time.NewTicker(c.Interval)

	go func() {
		for range ticker.C {
			for entry := range c.CacheEntry {

				if time.Since(c.CacheEntry[entry].CreatedAt) > c.Interval {

					c.Mu.Lock()
					delete(c.CacheEntry, entry)
					c.Mu.Unlock()
				}
			}
		}
	}()
}
