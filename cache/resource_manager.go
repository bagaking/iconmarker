package cache

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// ResourceManager manages different types of cached resources
type ResourceManager struct {
	svgCache    Cache
	fontCache   Cache
	imageCache  Cache
	ttlDuration time.Duration
	mu          sync.RWMutex
}

// NewResourceManager creates a new resource manager with specified cache sizes
func NewResourceManager(svgCacheSize, fontCacheSize, imageCacheSize int) *ResourceManager {
	return &ResourceManager{
		svgCache:    NewLRUCache(svgCacheSize),
		fontCache:   NewLRUCache(fontCacheSize),
		imageCache:  NewLRUCache(imageCacheSize),
		ttlDuration: 30 * time.Minute, // Default TTL
	}
}

// SetTTL sets the time-to-live duration for cached resources
func (rm *ResourceManager) SetTTL(duration time.Duration) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.ttlDuration = duration
}

// GetResource is a generic method to get a resource from the specified cache
func (rm *ResourceManager) GetResource(cacheType string, key string, cache Cache) (CacheItem, bool) {
	// Generate cache key with type prefix for better organization
	cacheKey := fmt.Sprintf("%s:%s", cacheType, key)

	// Try to get from cache
	item, found := cache.Get(cacheKey)
	return item, found
}

// PutResource is a generic method to store a resource in the specified cache
func (rm *ResourceManager) PutResource(cacheType string, key string, cache Cache, resource CacheItem) {
	// Generate cache key with type prefix
	cacheKey := fmt.Sprintf("%s:%s", cacheType, key)

	// Store in cache
	cache.Put(cacheKey, resource)
}

// GenerateKeyFromData generates a cache key from binary data using MD5 hash
func (rm *ResourceManager) GenerateKeyFromData(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// ClearAll clears all caches
func (rm *ResourceManager) ClearAll() {
	rm.svgCache.Clear()
	rm.fontCache.Clear()
	rm.imageCache.Clear()
}

// GetSVGCache returns the SVG cache
func (rm *ResourceManager) GetSVGCache() Cache {
	return rm.svgCache
}

// GetFontCache returns the font cache
func (rm *ResourceManager) GetFontCache() Cache {
	return rm.fontCache
}

// GetImageCache returns the image cache
func (rm *ResourceManager) GetImageCache() Cache {
	return rm.imageCache
}
