// Package cache provides unified caching infrastructure for Icon Marker
package cache

// CacheItem represents an item that can be stored in cache
type CacheItem interface {
	// Size returns estimated memory size of the item in bytes
	Size() int
}

// Resource represents a cacheable resource with cloning capability
type Resource interface {
	CacheItem
	// Clone returns a deep copy of the resource to avoid concurrent modification
	Clone() Resource
}

// Cache defines the common interface for all cache implementations
type Cache interface {
	// Get retrieves an item from cache by key
	Get(key string) (CacheItem, bool)

	// Put adds or updates an item in cache
	Put(key string, item CacheItem) bool

	// Remove removes an item from cache
	Remove(key string)

	// Clear removes all items from cache
	Clear()

	// Size returns the total number of items in cache
	Size() int
}
