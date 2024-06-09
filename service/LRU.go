package service

import (
	"container/list"
	"sync"
	"time"

	"lru/resources"
)

// CacheEntry represents an entry in the cache
type CacheEntry struct {
	key        string
	value      interface{}
	expiration time.Time
}

// LRUCache represents the LRU cache
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	lruCache := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
	// Start the cleanup goroutine
	go lruCache.cleanupExpiredEntries()
	return lruCache
}

// Get retrieves the value associated with the given key
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		entry := element.Value.(*CacheEntry)
		if time.Now().Before(entry.expiration) {
			return entry.value, true
		}
		c.delete(key)
	}
	return nil, false
}

// Set sets the value for the given key with an expiration time
func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var expire time.Time
	if expiration > 0 {
		expire = time.Now().Add(expiration * time.Second)
	} else {
		expire = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	}
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		entry := element.Value.(*CacheEntry)
		entry.value = value
		entry.expiration = expire
	} else {
		if c.list.Len() >= c.capacity {
			c.evict()
		}
		entry := &CacheEntry{
			key:        key,
			value:      value,
			expiration: expire,
		}
		element := c.list.PushFront(entry)
		c.cache[key] = element
	}
}

// Delete removes the entry associated with the given key
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(key)
}

// delete removes the entry associated with the given key without locking
func (c *LRUCache) delete(key string) {
	if element, ok := c.cache[key]; ok {
		c.list.Remove(element)
		delete(c.cache, key)
	}
}

// evict removes the least recently used entry from the cache
func (c *LRUCache) evict() {
	if c.list.Len() > c.capacity {
		element := c.list.Back()
		if element != nil {
			c.list.Remove(element)
			entry := element.Value.(*CacheEntry)
			delete(c.cache, entry.key)
		}
	}
}

// cleanupExpiredEntries periodically removes expired entries from the cache
func (c *LRUCache) cleanupExpiredEntries() {
	for {
		time.Sleep(time.Second) // adjust the interval as needed
		c.mu.Lock()
		now := time.Now()
		for element := c.list.Back(); element != nil; element = element.Prev() {
			entry := element.Value.(*CacheEntry)
			if now.After(entry.expiration) {
				c.delete(entry.key)
			} else {
				break
			}
		}
		c.mu.Unlock()
	}
}

func (c *LRUCache) GetAll() []resources.GetAllResponse {
	var items []resources.GetAllResponse
	var item resources.GetAllResponse
	for element := c.list.Front(); element != nil; element = element.Next() {

		entry := element.Value.(*CacheEntry)
		if time.Now().Before(entry.expiration) {
			item.Key = entry.key
			item.Value = entry.value
			item.Expiration = entry.expiration
			items = append(items, item)
		}
	}

	return items
}
