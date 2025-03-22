package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestResourceManagerPrefixesKeysAndClearsAllCaches(t *testing.T) {
	rm := NewResourceManager(2, 2, 2)
	item := &testItem{value: "svg", size: 3}
	rm.PutResource("svg", "same", rm.GetSVGCache(), item)
	rm.PutResource("font", "same", rm.GetFontCache(), item)

	if got, ok := rm.GetResource("svg", "same", rm.GetSVGCache()); !ok || got.(*testItem).value != "svg" {
		t.Fatalf("SVG resource lookup = %#v, found=%v", got, ok)
	}
	if _, ok := rm.GetResource("font", "same", rm.GetFontCache()); !ok {
		t.Fatal("expected font cache entry")
	}

	rm.ClearAll()
	if rm.GetSVGCache().Size() != 0 || rm.GetFontCache().Size() != 0 || rm.GetImageCache().Size() != 0 {
		t.Fatal("ClearAll left entries in a cache")
	}
}

func TestResourceManagerTTLIsAppliedToAllCaches(t *testing.T) {
	rm := NewResourceManager(2, 2, 2)
	rm.SetTTL(5 * time.Millisecond)
	item := &testItem{value: "value"}
	rm.PutResource("svg", "key", rm.GetSVGCache(), item)
	rm.PutResource("font", "key", rm.GetFontCache(), item)
	rm.PutResource("image", "key", rm.GetImageCache(), item)
	time.Sleep(20 * time.Millisecond)

	if _, ok := rm.GetResource("svg", "key", rm.GetSVGCache()); ok {
		t.Fatal("expected SVG item to expire")
	}
	if _, ok := rm.GetResource("font", "key", rm.GetFontCache()); ok {
		t.Fatal("expected font item to expire")
	}
	if _, ok := rm.GetResource("image", "key", rm.GetImageCache()); ok {
		t.Fatal("expected image item to expire")
	}
}

func TestResourceManagerGeneratesStableContentKeys(t *testing.T) {
	rm := NewResourceManager(1, 1, 1)
	data := []byte("same content")
	if got, want := rm.GenerateKeyFromData(data), rm.GenerateKeyFromData(bytes.Clone(data)); got != want {
		t.Fatalf("same content keys differ: %q vs %q", got, want)
	}
	if got := rm.GenerateKeyFromData([]byte("different")); got == rm.GenerateKeyFromData(data) {
		t.Fatalf("different content unexpectedly shared key %q", got)
	}
}
