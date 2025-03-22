package cache

import (
	"fmt"
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
	return &testItem{value: i.value, size: i.size}
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

func TestLRUCacheClearAndZeroCapacity(t *testing.T) {
	zero := NewLRUCache(0)
	zero.Put("a", &testItem{value: "a"})
	if got := zero.Size(); got != 0 {
		t.Fatalf("zero-capacity cache size = %d, want 0", got)
	}

	c := NewLRUCache(2)
	c.Put("a", &testItem{value: "a"})
	c.Put("b", &testItem{value: "b"})
	c.Clear()
	if got := c.Size(); got != 0 {
		t.Fatalf("cleared cache size = %d, want 0", got)
	}
}
