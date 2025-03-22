package iconmarker_test

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"testing"

	iconmarker "github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/filter"
	"github.com/bagaking/iconmarker/renderer"
)

func TestEndToEndJPEGTextSVGFilterAndPNG(t *testing.T) {
	background := image.NewRGBA(image.Rect(0, 0, 96, 64))
	draw.Draw(background, background.Bounds(), image.NewUniform(color.RGBA{R: 35, G: 45, B: 55, A: 255}), image.Point{}, draw.Src)
	var jpegData bytes.Buffer
	if err := jpeg.Encode(&jpegData, background, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode JPEG fixture: %v", err)
	}

	marker := iconmarker.NewIconMarker()
	img, err := marker.CreateImgWithFilters(
		nil,
		jpegData.Bytes(),
		[]string{"tint", "opacity"},
		[]filter.FilterOption{
			filter.TintOption{Color: [3]uint8{220, 80, 40}, Intensity: 0.2},
			filter.OpacityOption{Opacity: 1},
		},
		iconmarker.DrawTextOption{
			Text:      "TDD",
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		}.SetStaticSize(18),
	)
	if err != nil {
		t.Fatalf("CreateImgWithFilters() error = %v", err)
	}
	if got := img.Bounds(); got != background.Bounds() {
		t.Fatalf("text image bounds = %v, want %v", got, background.Bounds())
	}

	svg := renderer.NewSVGRenderer(marker.GetResourceManager())
	icon, err := svg.Render(renderer.SVGRenderOptions{
		SVGData: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><circle cx="5" cy="5" r="5" fill="#00ff00"/></svg>`),
		Width:   20,
		Height:  20,
	})
	if err != nil {
		t.Fatalf("SVG Render() error = %v", err)
	}
	draw.Draw(img, image.Rect(38, 22, 58, 42), icon, image.Point{}, draw.Over)

	encoded, err := encodePNG(img)
	if err != nil {
		t.Fatalf("PNG encode error = %v", err)
	}
	decoded, format, err := image.Decode(bytes.NewReader(encoded))
	if err != nil {
		t.Fatalf("PNG decode error = %v", err)
	}
	if format != "png" || decoded.Bounds() != background.Bounds() {
		t.Fatalf("decoded output format/bounds = %q/%v", format, decoded.Bounds())
	}
	if got := color.NRGBAModel.Convert(decoded.At(48, 32)).(color.NRGBA); got.G < 200 || got.A == 0 {
		t.Fatalf("SVG pixel = %#v, want opaque green contribution", got)
	}
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
