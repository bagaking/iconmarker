package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

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
		return &mutableResource{}
	}

	data := make([]byte, len(r.data))
	copy(data, r.data)
	return &mutableResource{data: data}
}

type stressResource struct {
	data []byte
}

func newStressResource(worker, iteration int) *stressResource {
	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(worker + iteration + i)
	}
	return &stressResource{data: data}
}

func (r *stressResource) Size() int {
	return len(r.data)
}

func (r *stressResource) Clone() Resource {
	data := make([]byte, len(r.data))
	copy(data, r.data)
	return &stressResource{data: data}
}

func (r *stressResource) mutate(seed byte) {
	for i := range r.data {
		r.data[i] ^= seed + byte(i)
	}
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

func TestLRUCacheConcurrentAccessStress(t *testing.T) {
	const (
		capacity   = 16
		workers    = 12
		iterations = 4000
		keyCount   = 32
	)

	cache := NewLRUCache(capacity)
	start := make(chan struct{})

	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		worker := worker
		wg.Add(1)
		go func() {
			defer wg.Done()

			rng := rand.New(rand.NewSource(int64(worker + 1)))
			<-start

			for iteration := 0; iteration < iterations; iteration++ {
				key := fmt.Sprintf("key-%02d", rng.Intn(keyCount))

				switch rng.Intn(8) {
				case 0, 1, 2:
					resource := newStressResource(worker, iteration)
					if ok := cache.Put(key, resource); !ok {
						t.Errorf("LRUCache.Put(%q, resource) = false, want true", key)
					}
					resource.mutate(byte(iteration))
				case 3, 4, 5:
					item, ok := cache.Get(key)
					if !ok {
						continue
					}
					resource, ok := item.(*stressResource)
					if !ok {
						t.Errorf("LRUCache.Get(%q) item type = %T, want *stressResource", key, item)
						continue
					}
					resource.mutate(byte(worker + iteration))
				case 6:
					cache.Remove(key)
				default:
					cache.Clear()
				}

				if size := cache.Size(); size < 0 || size > capacity {
					t.Errorf("LRUCache.Size() = %d, want between 0 and %d", size, capacity)
				}
			}
		}()
	}

	close(start)
	wg.Wait()
}
