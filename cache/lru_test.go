package cache

import "testing"

type testCacheItem int

func (i testCacheItem) Size() int {
	return int(i)
}

type mutableResource struct {
	data []byte
}

func (r *mutableResource) Size() int {
	return len(r.data)
}

func (r *mutableResource) Clone() Resource {
	if r == nil {
		return (*mutableResource)(nil)
	}

	data := make([]byte, len(r.data))
	copy(data, r.data)
	return &mutableResource{data: data}
}

func TestLRUCacheEvictsLeastRecentlyUsedItem(t *testing.T) {
	cache := NewLRUCache(2)

	if !cache.Put("a", testCacheItem(1)) {
		t.Fatal("expected first put to succeed")
	}
	cache.Put("b", testCacheItem(2))

	if _, ok := cache.Get("a"); !ok {
		t.Fatal("expected get to find key a")
	}

	cache.Put("c", testCacheItem(3))

	if _, ok := cache.Get("b"); ok {
		t.Fatal("expected key b to be evicted")
	}
	if _, ok := cache.Get("a"); !ok {
		t.Fatal("expected recently used key a to remain")
	}
	if _, ok := cache.Get("c"); !ok {
		t.Fatal("expected key c to be present")
	}
	if size := cache.Size(); size != 2 {
		t.Fatalf("expected size 2, got %d", size)
	}
}

func TestLRUCacheUpdatesExistingItem(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Put("a", testCacheItem(1))
	cache.Put("b", testCacheItem(2))

	if !cache.Put("a", testCacheItem(10)) {
		t.Fatal("expected update put to succeed")
	}
	cache.Put("c", testCacheItem(3))

	item, ok := cache.Get("a")
	if !ok {
		t.Fatal("expected updated key a to remain")
	}
	if got := item.(testCacheItem); got != 10 {
		t.Fatalf("expected updated value 10, got %d", got)
	}
	if _, ok := cache.Get("b"); ok {
		t.Fatal("expected key b to be evicted after key a update")
	}
}

func TestLRUCacheRemoveAndClear(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Put("a", testCacheItem(1))
	cache.Put("b", testCacheItem(2))

	cache.Remove("a")

	if _, ok := cache.Get("a"); ok {
		t.Fatal("expected removed key a to be absent")
	}
	if size := cache.Size(); size != 1 {
		t.Fatalf("expected size 1 after remove, got %d", size)
	}

	cache.Clear()

	if size := cache.Size(); size != 0 {
		t.Fatalf("expected size 0 after clear, got %d", size)
	}
	if _, ok := cache.Get("b"); ok {
		t.Fatal("expected key b to be absent after clear")
	}
}

func TestLRUCacheRejectsZeroCapacityPuts(t *testing.T) {
	cache := NewLRUCache(0)

	if cache.Put("a", testCacheItem(1)) {
		t.Fatal("expected zero-capacity put to fail")
	}
	if _, ok := cache.Get("a"); ok {
		t.Fatal("expected zero-capacity cache to store no items")
	}
	if size := cache.Size(); size != 0 {
		t.Fatalf("expected size 0, got %d", size)
	}
}

func TestLRUCacheClonesResourcesOnPut(t *testing.T) {
	cache := NewLRUCache(1)
	resource := &mutableResource{data: []byte("cached")}

	if !cache.Put("resource", resource) {
		t.Fatal("expected put to succeed")
	}
	resource.data[0] = 'm'

	item, ok := cache.Get("resource")
	if !ok {
		t.Fatal("expected resource to be cached")
	}
	got := item.(*mutableResource)
	if string(got.data) != "cached" {
		t.Fatalf("expected cached resource to be isolated from caller mutation, got %q", got.data)
	}
}

func TestLRUCacheClonesResourcesOnGet(t *testing.T) {
	cache := NewLRUCache(1)
	if !cache.Put("resource", &mutableResource{data: []byte("cached")}) {
		t.Fatal("expected put to succeed")
	}

	item, ok := cache.Get("resource")
	if !ok {
		t.Fatal("expected resource to be cached")
	}
	item.(*mutableResource).data[0] = 'm'

	item, ok = cache.Get("resource")
	if !ok {
		t.Fatal("expected resource to remain cached")
	}
	got := item.(*mutableResource)
	if string(got.data) != "cached" {
		t.Fatalf("expected cached resource to be isolated from returned item mutation, got %q", got.data)
	}
}
