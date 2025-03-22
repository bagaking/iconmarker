// Package cache provides unified caching infrastructure for Icon Marker
package cache

// CacheItem represents an item that can be stored in cache. Values passed to
// Cache.Put must not be nil or typed nil.
type CacheItem interface {
	// Size returns estimated memory size of the item in bytes
	Size() int
}

// Resource represents a cacheable resource with cloning capability. Resource
// values passed to Cache.Put must not be nil or typed nil.
type Resource interface {
	CacheItem

	// Clone returns an isolated copy of mutable resource state to avoid
	// concurrent modification; immutable resources may return themselves. Clone
	// must be cheap, must not call back into cache operations, must not return
	// nil, and must be safe to call while the cache holds its internal lock.
	Clone() Resource
}

// Cache defines the common interface for all cache implementations
type Cache interface {
	// Get retrieves an item from cache by key
	Get(key string) (CacheItem, bool)

	// Put adds or updates a non-nil item in cache
	Put(key string, item CacheItem) bool

	// Remove removes an item from cache
	Remove(key string)

	// Clear removes all items from cache
	Clear()

	// Size returns the total number of items in cache
	Size() int
}
