package filter

import (
	"image"
	"image/color"
	"testing"
)

func TestApplyFiltersPreservesPixelsForNonZeroBounds(t *testing.T) {
	srcBounds := image.Rect(10, 20, 12, 21)
	src := image.NewRGBA(srcBounds)
	src.Set(10, 20, color.RGBA{R: 10, G: 20, B: 30, A: 255})
	src.Set(11, 20, color.RGBA{R: 250, G: 240, B: 230, A: 128})

	got, err := NewFilterManager().ApplyFilters(src, []string{"invert"}, []FilterOption{InvertOption{}})
	if err != nil {
		t.Fatalf("ApplyFilters(non-zero bounds image) error = %v, want nil", err)
	}
	if got.Bounds() != srcBounds {
		t.Fatalf("ApplyFilters(non-zero bounds image).Bounds() = %v, want %v", got.Bounds(), srcBounds)
	}

	tests := []struct {
		point image.Point
		want  color.RGBA
	}{
		{
			point: image.Pt(10, 20),
			want:  color.RGBA{R: 245, G: 235, B: 225, A: 255},
		},
		{
			point: image.Pt(11, 20),
			want:  color.RGBA{R: 5, G: 15, B: 25, A: 128},
		},
	}

	for _, tt := range tests {
		gotColor := color.RGBAModel.Convert(got.At(tt.point.X, tt.point.Y)).(color.RGBA)
		if gotColor != tt.want {
			t.Errorf("ApplyFilters(non-zero bounds image).At(%v) = %v, want %v", tt.point, gotColor, tt.want)
		}
	}
}
