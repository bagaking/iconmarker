package assets

import (
	"bytes"
	"testing"

	"github.com/golang/freetype"
	"github.com/srwiley/oksvg"
)

func TestEmbeddedDefaultFontIsParseable(t *testing.T) {
	data, err := GetDefaultFont()
	if err != nil {
		t.Fatalf("GetDefaultFont() error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("default font is empty")
	}
	if _, err := freetype.ParseFont(data); err != nil {
		t.Fatalf("default font cannot be parsed: %v", err)
	}
}

func TestAllEmbeddedIconsAreReadableAndParseable(t *testing.T) {
	icons := AllIcons()
	if len(icons) != 22 {
		t.Fatalf("embedded icon count = %d, want 22", len(icons))
	}
	seen := make(map[IconType]bool, len(icons))
	for _, icon := range icons {
		if seen[icon] {
			t.Fatalf("duplicate icon %q", icon)
		}
		seen[icon] = true
		data, err := icon.Load()
		if err != nil {
			t.Fatalf("load %s: %v", icon, err)
		}
		if len(data) == 0 {
			t.Fatalf("icon %s is empty", icon)
		}
		if _, err := oksvg.ReadIconStream(bytes.NewReader(data)); err != nil {
			t.Fatalf("parse %s: %v", icon, err)
		}
	}
}

func TestIconLookupAndParsingBoundaries(t *testing.T) {
	icon, err := ParseIconType("heart.svg")
	if err != nil || icon != IconHeart {
		t.Fatalf("ParseIconType(heart.svg) = %q, %v", icon, err)
	}
	if _, err := ParseIconType("does-not-exist"); err == nil {
		t.Fatal("expected unknown icon error")
	}
	if _, err := GetSVGIcon("../AlibabaPuHuiTi-3-105-Heavy.ttf"); err == nil {
		t.Fatal("expected path traversal lookup to fail")
	}
}
