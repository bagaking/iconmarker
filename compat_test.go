package iconmarker_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	"github.com/bagaking/iconmarker"
)

func solidJPEG(t *testing.T, width, height int, c color.Color) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode background: %v", err)
	}
	return buf.Bytes()
}

func countDarkPixels(img image.Image, threshold uint8) int {
	count := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if uint8(r>>8) < threshold || uint8(g>>8) < threshold || uint8(b>>8) < threshold {
				count++
			}
		}
	}
	return count
}

func TestCreateImgUsesEmbeddedFontAndJPEGBackground(t *testing.T) {
	background := solidJPEG(t, 160, 80, color.RGBA{R: 248, G: 248, B: 248, A: 255})

	img, err := iconmarker.CreateImg(nil, background,
		iconmarker.DrawTextOption{
			FontColor: color.Black,
			Text:      "TDD",
		}.SetStaticSize(24),
	)
	if err != nil {
		t.Fatalf("CreateImg with embedded font: %v", err)
	}
	if got, want := img.Bounds().Dx(), 160; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}
	if got, want := img.Bounds().Dy(), 80; got != want {
		t.Fatalf("height = %d, want %d", got, want)
	}
	if dark := countDarkPixels(img, 180); dark == 0 {
		t.Fatal("expected rendered text to change at least one background pixel")
	}
}

func TestCreateImgRejectsInvalidJPEG(t *testing.T) {
	_, err := iconmarker.CreateImg(nil, []byte("not a jpeg"),
		iconmarker.DrawTextOption{FontColor: color.Black, Text: "TDD"}.SetStaticSize(16),
	)
	if err == nil {
		t.Fatal("expected invalid JPEG to return an error")
	}
}

func TestDrawCenteredFontRejectsInvalidStaticSize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 80, 40))
	err := iconmarker.DrawCenteredFont(nil, img, iconmarker.DrawTextOption{
		FontColor: color.Black,
		Text:      "TDD",
	})
	if err == nil {
		t.Fatal("expected zero font size without adaptation to return an error")
	}
}

func TestCompatibilityApplyFilterDoesNotMutateSource(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 20, G: 40, B: 60, A: 255})

	filtered, err := iconmarker.ApplyFilter(src, "invert", nil)
	if err != nil {
		t.Fatalf("ApplyFilter: %v", err)
	}
	if got := src.RGBAAt(0, 0); got != (color.RGBA{R: 20, G: 40, B: 60, A: 255}) {
		t.Fatalf("source mutated: got %#v", got)
	}
	if got := filtered.(*image.RGBA).RGBAAt(0, 0); got == (color.RGBA{R: 20, G: 40, B: 60, A: 255}) {
		t.Fatal("expected filtered image to differ from source")
	}
}
