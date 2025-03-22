package cache

import (
	"sync"
	"time"
)

// lruItem represents an item in the LRU cache
type lruItem struct {
	key       string
	value     CacheItem
	expiresAt time.Time
	prev      *lruItem
	next      *lruItem
}

// LRUCache implements an LRU (Least Recently Used) cache
type LRUCache struct {
	capacity int                 // Maximum number of items
	size     int                 // Current number of items
	items    map[string]*lruItem // Map for O(1) lookup
	head     *lruItem            // Most recently used item
	tail     *lruItem            // Least recently used item
	ttl      time.Duration       // Expiration duration; <= 0 means no expiration
	mu       sync.RWMutex        // For thread safety
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 0 {
		capacity = 0
	}
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*lruItem),
	}
}

// Get retrieves an item from cache
func (c *LRUCache) Get(key string) (CacheItem, bool) {
	// A cache hit changes recency, so lookup, expiration, and list mutation must
	// happen under one write lock.  The previous implementation released the
	// read lock before acquiring the write lock, which allowed Remove/Clear to
	// invalidate the item in between and introduced a race.
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if c.isExpired(item, time.Now()) {
		c.deleteItem(item)
		return nil, false
	}

	c.moveToFront(item)
	return cloneCacheItem(item.value), true
}

// Put adds or updates an item in cache
func (c *LRUCache) Put(key string, value CacheItem) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// A zero-capacity cache is a valid disabled cache.  Keep Put's historical
	// success return value while ensuring it never retains an item.
	if c.capacity <= 0 {
		return true
	}

	expiresAt := c.expiryFrom(time.Now())

	// Check if item already exists
	if item, found := c.items[key]; found {
		item.value = cloneCacheItem(value)
		item.expiresAt = expiresAt
		c.moveToFront(item)
		return true
	}

	// Create new item
	item := &lruItem{
		key:       key,
		value:     cloneCacheItem(value),
		expiresAt: expiresAt,
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
		c.deleteItem(item)
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
	c.mu.Lock()
	defer c.mu.Unlock()
	c.purgeExpired(time.Now())
	return c.size
}

// SetTTL configures the time-to-live for entries in the cache.  A duration of
// zero or less disables expiration.  Existing entries are rebased from the
// time of this call so changing the policy has deterministic behavior.
func (c *LRUCache) SetTTL(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ttl = duration
	now := time.Now()
	for _, item := range c.items {
		item.expiresAt = c.expiryFrom(now)
	}
}

// TTL reports the configured expiration duration.  It is primarily useful to
// introspect a cache in diagnostics and tests.
func (c *LRUCache) TTL() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ttl
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

	c.deleteItem(c.tail)
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
	item.prev = nil
	item.next = nil
}

func (c *LRUCache) expiryFrom(now time.Time) time.Time {
	if c.ttl <= 0 {
		return time.Time{}
	}
	return now.Add(c.ttl)
}

func (c *LRUCache) isExpired(item *lruItem, now time.Time) bool {
	return !item.expiresAt.IsZero() && !now.Before(item.expiresAt)
}

// deleteItem removes an item from both the linked list and the lookup map.
// Callers must hold c.mu.
func (c *LRUCache) deleteItem(item *lruItem) {
	if item == nil {
		return
	}
	c.removeItem(item)
	if _, found := c.items[item.key]; found {
		delete(c.items, item.key)
		c.size--
	}
}

// purgeExpired removes all entries whose TTL has elapsed.  Callers must hold
// c.mu.
func (c *LRUCache) purgeExpired(now time.Time) {
	for _, item := range c.items {
		if c.isExpired(item, now) {
			c.deleteItem(item)
		}
	}
}

func cloneCacheItem(item CacheItem) CacheItem {
	resource, ok := item.(Resource)
	if !ok {
		return item
	}
	return resource.Clone()
}
