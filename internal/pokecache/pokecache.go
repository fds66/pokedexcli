package pokecache

import (
	"fmt"
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
	//newCache.reapLoop()
	return &newCache
}

func (c *Cache) Add(key string, value []byte) bool {
	fmt.Println("Adding new data to the cache")
	c.Mu.Lock()
	defer c.Mu.Unlock()
	newEntry := CacheEntry{
		Val:       value,
		CreatedAt: time.Now(),
	}

	c.CacheEntry[key] = newEntry
	fmt.Printf("Added %s to the cache", key)
	return true

}

func (c *Cache) Get(key string) ([]byte, bool) {
	fmt.Println("Retrieving data from the cache")
	entry, exists := c.CacheEntry[key]
	if exists {
		return entry.Val, true
	} else {
		return nil, false
	}
}

//func (c *Cache) reapLoop() {
//return
//}
