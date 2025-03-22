package core

import (
	"image"
	"image/color"
	"testing"
)

func TestIconMarkerApplyFilterPreservesNonZeroImageOrigin(t *testing.T) {
	img := image.NewRGBA(image.Rect(10, 20, 12, 22))
	img.SetRGBA(10, 20, color.RGBA{R: 10, G: 20, B: 30, A: 255})
	filtered, err := NewIconMarker().ApplyFilter(img, "invert", nil)
	if err != nil {
		t.Fatalf("ApplyFilter() error = %v", err)
	}
	if got := filtered.Bounds(); got != img.Bounds() {
		t.Fatalf("filtered bounds = %v, want %v", got, img.Bounds())
	}
	if got := color.NRGBAModel.Convert(filtered.At(10, 20)).(color.NRGBA); got.R != 245 || got.G != 235 || got.B != 225 {
		t.Fatalf("filtered pixel = %#v, want inverted source pixel", got)
	}
}
