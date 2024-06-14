package service

import (
	"lru/resources"
	"sync"
	"time"
)

// LRUCache represents the LRU cache
type LRUCache struct {
	capacity int
	cache    map[string]*Node
	dll      DoublyLinkedList
	mu       sync.Mutex
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	lruCache := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		dll:      NewDll(),
	}
	// Start the cleanup goroutine
	go lruCache.cleanupExpiredEntries()
	return lruCache
}

// Get retrieves the value associated with the given key
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.cache[key]; ok {
		c.dll.MoveToFront(node)
		if time.Now().Before(node.expiration) {
			return node.value, true
		}
		c.deleteNode(node)
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
	if node, ok := c.cache[key]; ok {
		node.value = value
		node.expiration = expire
		c.dll.MoveToFront(node)
	} else {
		if len(c.cache) >= c.capacity {
			c.evict()
		}
		node := &Node{
			key:        key,
			value:      value,
			expiration: expire,
		}
		c.dll.PushFront(node)
		c.cache[key] = node
	}
}

// Delete removes the entry associated with the given key
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.cache[key]; ok {
		c.deleteNode(node)
	}
}

// deleteNode removes a node from the linked list and cache
func (c *LRUCache) deleteNode(node *Node) {
	c.dll.Remove(node)
	delete(c.cache, node.key)
}

// evict removes the least recently used entry from the cache
func (c *LRUCache) evict() {
	node := c.dll.RemoveLast()
	if node != nil {
		delete(c.cache, node.key)
	}
}

// cleanupExpiredEntries periodically removes expired entries from the cache
func (c *LRUCache) cleanupExpiredEntries() {
	for {
		time.Sleep(time.Second)
		c.mu.Lock()
		now := time.Now()
		for node := c.dll.(*Dll).tail; node != nil; node = node.prev {
			if now.After(node.expiration) {
				c.deleteNode(node)
			} else {
				break
			}
		}
		c.mu.Unlock()
	}
}

// GetAll retrieves all non-expired entries from the cache
func (c *LRUCache) GetAll() []resources.GetAllResponse {
	c.mu.Lock()
	defer c.mu.Unlock()

	var items []resources.GetAllResponse
	for node := c.dll.(*Dll).head; node != nil; node = node.next {
		if time.Now().Before(node.expiration) {
			item := resources.GetAllResponse{
				Key:        node.key,
				Value:      node.value,
				Expiration: node.expiration,
			}
			items = append(items, item)
		}
	}
	return items
}
