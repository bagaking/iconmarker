package filter

import (
	"image"
	"image/color"
	"reflect"
	"sync"
	"testing"
)

func onePixel(c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.SetRGBA(0, 0, c)
	return img
}

func TestGrayscalePreservesAlpha(t *testing.T) {
	img := onePixel(color.RGBA{R: 100, G: 150, B: 200, A: 73})
	if err := NewGrayscaleFilter().Apply(img, GrayscaleOption{PreserveAlpha: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if got, want := img.RGBAAt(0, 0), (color.RGBA{R: 140, G: 140, B: 140, A: 73}); got != want {
		t.Fatalf("pixel = %#v, want %#v", got, want)
	}
}

func TestOpacityValidatesAndScalesAlpha(t *testing.T) {
	if err := NewOpacityFilter().Apply(onePixel(color.RGBA{R: 10, G: 20, B: 30, A: 200}), OpacityOption{Opacity: 1.5}); err != ErrInvalidOpacity {
		t.Fatalf("invalid opacity error = %v, want %v", err, ErrInvalidOpacity)
	}

	img := onePixel(color.RGBA{R: 10, G: 20, B: 30, A: 200})
	if err := NewOpacityFilter().Apply(img, OpacityOption{Opacity: 0.5}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := img.RGBAAt(0, 0), (color.RGBA{R: 10, G: 20, B: 30, A: 100}); got != want {
		t.Fatalf("pixel = %#v, want %#v", got, want)
	}
}

func TestInvertCanInvertAlpha(t *testing.T) {
	img := onePixel(color.RGBA{R: 10, G: 20, B: 30, A: 40})
	if err := NewInvertFilter().Apply(img, InvertOption{InvertAlpha: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := img.RGBAAt(0, 0), (color.RGBA{R: 245, G: 235, B: 225, A: 215}); got != want {
		t.Fatalf("pixel = %#v, want %#v", got, want)
	}
}

func TestTintIntensityZeroPreservesOriginalColor(t *testing.T) {
	img := onePixel(color.RGBA{R: 200, G: 100, B: 50, A: 255})
	if err := NewTintFilter().Apply(img, TintOption{Color: [3]uint8{0, 0, 255}, Intensity: 0}); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := img.RGBAAt(0, 0), (color.RGBA{R: 200, G: 100, B: 50, A: 255}); got != want {
		t.Fatalf("pixel = %#v, want original %#v", got, want)
	}
}

func TestTintValidatesIntensity(t *testing.T) {
	err := NewTintFilter().Apply(onePixel(color.RGBA{A: 255}), TintOption{Intensity: -0.1})
	if err != ErrInvalidIntensity {
		t.Fatalf("invalid intensity error = %v, want %v", err, ErrInvalidIntensity)
	}
}

func TestFilterManagerRegistersComposite(t *testing.T) {
	manager := NewFilterManager()
	img := onePixel(color.RGBA{R: 10, G: 20, B: 30, A: 255})
	option := CompositeOption{
		Filters:     []Filter{NewInvertFilter()},
		Options:     []FilterOption{InvertOption{InvertAlpha: false}},
		StopOnError: true,
	}
	if err := manager.Apply(img, "composite", option); err != nil {
		t.Fatalf("composite Apply: %v", err)
	}
	if got, want := img.RGBAAt(0, 0), (color.RGBA{R: 245, G: 235, B: 225, A: 255}); got != want {
		t.Fatalf("pixel = %#v, want %#v", got, want)
	}
}

func TestApplyFiltersPreservesSourceAndBoundsOrigin(t *testing.T) {
	src := image.NewRGBA(image.Rect(10, 20, 11, 21))
	src.SetRGBA(10, 20, color.RGBA{R: 10, G: 20, B: 30, A: 255})
	wantSource := src.RGBAAt(10, 20)

	gotImage, err := NewFilterManager().ApplyFilters(src, []string{"invert"}, []FilterOption{InvertOption{}})
	if err != nil {
		t.Fatalf("ApplyFilters: %v", err)
	}
	got := gotImage.(*image.RGBA).RGBAAt(10, 20)
	if want := (color.RGBA{R: 245, G: 235, B: 225, A: 255}); got != want {
		t.Fatalf("filtered pixel = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(src.RGBAAt(10, 20), wantSource) {
		t.Fatalf("source pixel changed: got %#v, want %#v", src.RGBAAt(10, 20), wantSource)
	}
}

func TestFilterManagerConcurrentRegistrationAndLookup(t *testing.T) {
	manager := NewFilterManager()
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				manager.Register("custom", NewInvertFilter())
				if _, ok := manager.Get("custom"); !ok {
					t.Errorf("worker %d: custom filter missing", worker)
				}
			}
		}(worker)
	}
	wg.Wait()
}

func TestFilterManagerRejectsNilSource(t *testing.T) {
	if _, err := NewFilterManager().ApplyFilters(nil, []string{"invert"}, nil); err != ErrNilImage {
		t.Fatalf("nil source error = %v, want %v", err, ErrNilImage)
	}
}
