package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestFilterExamplesPropagateWriteErrors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	missingDir := filepath.Join(t.TempDir(), "does-not-exist")
	if err := compositeSingleFilter(img, missingDir); err == nil {
		t.Fatal("expected output write error to propagate")
	}
}

func TestFilterExamplesWriteAllOutputs(t *testing.T) {
	outputDir := t.TempDir()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 20), G: uint8(y * 20), B: 100, A: 255})
		}
	}
	if err := compositeSingleFilter(img, outputDir); err != nil {
		t.Fatalf("compositeSingleFilter() error = %v", err)
	}
	if err := sequentialFilters(img, outputDir); err != nil {
		t.Fatalf("sequentialFilters() error = %v", err)
	}
	if err := customCompositeFilter(img, outputDir); err != nil {
		t.Fatalf("customCompositeFilter() error = %v", err)
	}
	for _, name := range []string{"composite_filter.jpg", "sequential_filters.jpg", "custom_composite.jpg"} {
		if info, err := os.Stat(filepath.Join(outputDir, name)); err != nil || info.Size() == 0 {
			t.Fatalf("output %s unavailable: info=%v err=%v", name, info, err)
		}
	}
}
