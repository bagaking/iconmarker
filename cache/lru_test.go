package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type testItem struct {
	value string
	size  int
}

func (i *testItem) Size() int { return i.size }

func (i *testItem) Clone() Resource {
	if i == nil {
		return &testItem{}
	}
	return &testItem{value: i.value, size: i.size}
}

type testCacheItem int

func (i testCacheItem) Size() int { return int(i) }

type mutableResource struct {
	data []byte
}

func (r *mutableResource) Size() int {
	if r == nil {
		return 0
	}
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
	if r == nil {
		return 0
	}
	return len(r.data)
}

func (r *stressResource) Clone() Resource {
	if r == nil {
		return &stressResource{}
	}
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
	c := NewLRUCache(2)
	c.Put("a", &testItem{value: "a"})
	c.Put("b", &testItem{value: "b"})

	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected a to be present")
	}
	c.Put("c", &testItem{value: "c"})

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected b to be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected a to remain after being touched")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("expected c to be present")
	}
	if got := c.Size(); got != 2 {
		t.Fatalf("cache size = %d, want 2", got)
	}
}

func TestLRUCacheUpdateMovesItemToFront(t *testing.T) {
	c := NewLRUCache(2)
	c.Put("a", &testItem{value: "old"})
	c.Put("b", &testItem{value: "b"})
	c.Put("a", &testItem{value: "new"})
	c.Put("c", &testItem{value: "c"})

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected b to be evicted after updating a")
	}
	item, ok := c.Get("a")
	if !ok || item.(*testItem).value != "new" {
		t.Fatalf("updated item = %#v, found=%v", item, ok)
	}
}

func TestLRUCacheRemoveAndClear(t *testing.T) {
	c := NewLRUCache(2)
	c.Put("a", testCacheItem(1))
	c.Put("b", testCacheItem(2))

	c.Remove("a")
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected removed key a to be absent")
	}
	if got := c.Size(); got != 1 {
		t.Fatalf("cache size after remove = %d, want 1", got)
	}

	c.Clear()
	if got := c.Size(); got != 0 {
		t.Fatalf("cache size after clear = %d, want 0", got)
	}
	if _, ok := c.Get("b"); ok {
		t.Fatal("expected key b to be absent after clear")
	}
}

func TestLRUCacheZeroCapacityDoesNotRetainItems(t *testing.T) {
	c := NewLRUCache(0)
	if !c.Put("a", testCacheItem(1)) {
		t.Fatal("expected zero-capacity Put to preserve its success result")
	}
	if _, ok := c.Get("a"); ok {
		t.Fatal("expected zero-capacity cache to store no items")
	}
	if got := c.Size(); got != 0 {
		t.Fatalf("zero-capacity cache size = %d, want 0", got)
	}
}

func TestLRUCacheExpiresItemsAfterTTL(t *testing.T) {
	c := NewLRUCache(2)
	c.SetTTL(5 * time.Millisecond)
	c.Put("a", &testItem{value: "a"})
	time.Sleep(20 * time.Millisecond)

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected expired item to be a cache miss")
	}
	if got := c.Size(); got != 0 {
		t.Fatalf("expired item count = %d, want 0", got)
	}
}

func TestLRUCacheClonesResourcesOnPut(t *testing.T) {
	c := NewLRUCache(1)
	resource := &mutableResource{data: []byte("cached")}
	if !c.Put("resource", resource) {
		t.Fatal("expected Put to succeed")
	}
	resource.data[0] = 'm'

	item, ok := c.Get("resource")
	if !ok {
		t.Fatal("expected resource to be cached")
	}
	if got := item.(*mutableResource); string(got.data) != "cached" {
		t.Fatalf("cached resource changed with caller mutation: %q", got.data)
	}
}

func TestLRUCacheClonesResourcesOnGet(t *testing.T) {
	c := NewLRUCache(1)
	if !c.Put("resource", &mutableResource{data: []byte("cached")}) {
		t.Fatal("expected Put to succeed")
	}

	item, ok := c.Get("resource")
	if !ok {
		t.Fatal("expected resource to be cached")
	}
	item.(*mutableResource).data[0] = 'm'

	item, ok = c.Get("resource")
	if !ok {
		t.Fatal("expected resource to remain cached")
	}
	if got := item.(*mutableResource); string(got.data) != "cached" {
		t.Fatalf("cached resource changed with returned mutation: %q", got.data)
	}
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	c := NewLRUCache(32)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		worker := worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				key := fmt.Sprintf("%d-%d", worker, i%64)
				c.Put(key, &testItem{value: key, size: 1})
				c.Get(key)
				if i%7 == 0 {
					c.Remove(key)
				}
			}
		}()
	}
	wg.Wait()
	if got := c.Size(); got < 0 || got > 32 {
		t.Fatalf("invalid final cache size: %d", got)
	}
}

func TestLRUCacheConcurrentAccessStress(t *testing.T) {
	const (
		capacity   = 16
		workers    = 12
		iterations = 4000
		keyCount   = 32
	)

	c := NewLRUCache(capacity)
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
					if !c.Put(key, resource) {
						t.Errorf("LRUCache.Put(%q, resource) = false, want true", key)
					}
					resource.mutate(byte(iteration))
				case 3, 4, 5:
					item, ok := c.Get(key)
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
					c.Remove(key)
				default:
					c.Clear()
				}

				if size := c.Size(); size < 0 || size > capacity {
					t.Errorf("LRUCache.Size() = %d, want between 0 and %d", size, capacity)
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}
