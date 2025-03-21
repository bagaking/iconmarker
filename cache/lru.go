package cache

import (
	"sync"
)

// lruItem represents an item in the LRU cache
type lruItem struct {
	key   string
	value CacheItem
	prev  *lruItem
	next  *lruItem
}

// LRUCache implements an LRU (Least Recently Used) cache
type LRUCache struct {
	capacity int                 // Maximum number of items
	size     int                 // Current number of items
	items    map[string]*lruItem // Map for O(1) lookup
	head     *lruItem            // Most recently used item
	tail     *lruItem            // Least recently used item
	mu       sync.RWMutex        // For thread safety
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*lruItem),
	}
}

// Get retrieves an item from cache
func (c *LRUCache) Get(key string) (CacheItem, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	// Move item to front (most recently used)
	c.mu.Lock()
	c.moveToFront(item)
	c.mu.Unlock()

	return item.value, true
}

// Put adds or updates an item in cache
func (c *LRUCache) Put(key string, value CacheItem) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if item already exists
	if item, found := c.items[key]; found {
		item.value = value
		c.moveToFront(item)
		return true
	}

	// Create new item
	item := &lruItem{
		key:   key,
		value: value,
	}

	// Add to cache
	c.items[key] = item

	// If this is the first item
	if c.head == nil {
		c.head = item
		c.tail = item
	} else {
		// Add to front
		item.next = c.head
		c.head.prev = item
		c.head = item
	}

	c.size++

	// Evict if over capacity
	if c.size > c.capacity {
		c.evictLRU()
	}

	return true
}

// Remove removes an item from cache
func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.items[key]; found {
		c.removeItem(item)
		delete(c.items, key)
		c.size--
	}
}

// Clear removes all items from cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*lruItem)
	c.head = nil
	c.tail = nil
	c.size = 0
}

// Size returns the number of items in cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}

// moveToFront moves an item to the front of the list (most recently used)
func (c *LRUCache) moveToFront(item *lruItem) {
	// Already at front
	if item == c.head {
		return
	}

	// Remove from current position
	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
	if item == c.tail {
		c.tail = item.prev
	}

	// Move to front
	item.prev = nil
	item.next = c.head
	c.head.prev = item
	c.head = item
}

// evictLRU removes the least recently used item
func (c *LRUCache) evictLRU() {
	if c.tail == nil {
		return
	}

	// Remove from map
	delete(c.items, c.tail.key)

	// Remove from list
	c.removeItem(c.tail)
	c.size--
}

// removeItem removes an item from the linked list
func (c *LRUCache) removeItem(item *lruItem) {
	// Update neighbors
	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}

	// Update head/tail if needed
	if item == c.head {
		c.head = item.next
	}
	if item == c.tail {
		c.tail = item.prev
	}
}
