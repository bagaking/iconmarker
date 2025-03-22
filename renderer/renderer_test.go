package renderer

import (
	"image"
	"image/color"
	"testing"

	"github.com/bagaking/iconmarker/assets"
	"github.com/bagaking/iconmarker/cache"
)

const testSVGTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 10 10"><rect x="0" y="0" width="10" height="10" fill="%s"/></svg>`

type testSVGOption struct {
	data          []byte
	width, height int
	err           error
}

type invalidRenderOption struct{}

func (invalidRenderOption) ValidateOption() error { return nil }

func (o testSVGOption) ValidateOption() error { return o.err }
func (o testSVGOption) GetSVGData() []byte    { return o.data }
func (o testSVGOption) GetDimensions() (int, int) {
	return o.width, o.height
}

type testTextOption struct {
	text                string
	maxWidth, maxHeight int
	fontSize            float64
	fontData            []byte
	fontColor           color.Color
	xOffset, yOffset    int
	err                 error
}

func (o testTextOption) ValidateOption() error { return o.err }
func (o testTextOption) GetText() string       { return o.text }
func (o testTextOption) GetMaxSize() (int, int) {
	return o.maxWidth, o.maxHeight
}
func (o testTextOption) GetFontSize() float64    { return o.fontSize }
func (o testTextOption) GetColor() interface{}   { return o.fontColor }
func (o testTextOption) GetPosition() (int, int) { return o.xOffset, o.yOffset }
func (o testTextOption) GetFontData() []byte     { return o.fontData }

func TestSVGRendererRendersDimensionsAndTransparentBackground(t *testing.T) {
	r := NewSVGRenderer(cache.NewResourceManager(4, 1, 1))
	data := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 10 10"><rect x="0" y="0" width="5" height="5" fill="#ff0000"/></svg>`)

	img, err := r.Render(testSVGOption{data: data, width: 20, height: 12})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got := img.Bounds(); got != image.Rect(0, 0, 20, 12) {
		t.Fatalf("bounds = %v, want %v", got, image.Rect(0, 0, 20, 12))
	}
	if got := color.NRGBAModel.Convert(img.At(2, 2)).(color.NRGBA); got.R != 255 || got.G != 0 || got.B != 0 || got.A != 255 {
		t.Fatalf("center pixel = %#v, want opaque red", got)
	}
	if got := color.NRGBAModel.Convert(img.At(19, 11)).(color.NRGBA); got.A != 0 {
		t.Fatalf("background alpha = %d, want transparent", got.A)
	}

	if got := r.resourceManager.GetSVGCache().Size(); got != 1 {
		t.Fatalf("SVG cache size = %d, want one entry", got)
	}
	// Dimensions are part of rendering, not the raw SVG resource key.  A
	// second size should reuse the same parsed-data cache entry safely.
	if _, err := r.Render(testSVGOption{data: data, width: 7, height: 9}); err != nil {
		t.Fatalf("second Render() error = %v", err)
	}
	if got := r.resourceManager.GetSVGCache().Size(); got != 1 {
		t.Fatalf("SVG cache size after resize = %d, want one entry", got)
	}
}

func TestSVGRendererValidatesOptionsAndPropagatesParseErrors(t *testing.T) {
	r := NewSVGRenderer(cache.NewResourceManager(2, 1, 1))

	if _, err := r.Render(testSVGOption{data: []byte("<svg/>")}); err == nil {
		t.Fatal("expected invalid dimensions error")
	}
	if _, err := r.Render(testSVGOption{data: []byte("not svg"), width: 4, height: 4}); err == nil {
		t.Fatal("expected SVG parse error")
	}
	if _, err := r.Render(testSVGOption{data: []byte("<svg>"), width: 4, height: 4}); err == nil {
		t.Fatal("expected malformed SVG error")
	}
	if _, err := r.Render(invalidRenderOption{}); err == nil {
		t.Fatal("expected option type error")
	}
}

func TestSVGRendererRenderMultiplePreservesOrderAndReportsIndex(t *testing.T) {
	r := NewSVGRenderer(cache.NewResourceManager(4, 1, 1))
	red := testSVGOption{data: []byte(sprintfSVG("#ff0000")), width: 3, height: 3}
	blue := testSVGOption{data: []byte(sprintfSVG("#0000ff")), width: 3, height: 3}
	images, err := r.RenderMultiple([]SVGRenderOption{red, blue})
	if err != nil {
		t.Fatalf("RenderMultiple() error = %v", err)
	}
	if got := color.NRGBAModel.Convert(images[0].At(1, 1)).(color.NRGBA); got.R != 255 || got.B != 0 {
		t.Fatalf("first image = %#v, want red", got)
	}
	if got := color.NRGBAModel.Convert(images[1].At(1, 1)).(color.NRGBA); got.B != 255 || got.R != 0 {
		t.Fatalf("second image = %#v, want blue", got)
	}

	bad := testSVGOption{data: []byte("<svg>"), width: 3, height: 3}
	if _, err := r.RenderMultiple([]SVGRenderOption{red, bad}); err == nil {
		t.Fatal("expected batch error")
	} else if got := err.Error(); !contains(got, "SVG 1") {
		t.Fatalf("batch error = %q, want index", got)
	}
}

func TestTextRendererUsesSeparateFontDataAndColor(t *testing.T) {
	fontData, err := assets.GetDefaultFont()
	if err != nil {
		t.Fatalf("GetDefaultFont() error = %v", err)
	}
	r := NewTextRenderer(cache.NewResourceManager(1, 2, 1))
	opt := testTextOption{
		text:     "TDD",
		maxWidth: 120, maxHeight: 60,
		fontData:  fontData,
		fontColor: color.RGBA{R: 220, G: 20, B: 30, A: 255},
	}
	img, err := r.Render(opt)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got := img.Bounds(); got != image.Rect(0, 0, 120, 60) {
		t.Fatalf("text bounds = %v, want 120x60", got)
	}
	if !hasColorNear(img, color.RGBA{R: 220, G: 20, B: 30, A: 255}) {
		t.Fatal("rendered text did not contain the requested color")
	}

	// An explicit, invalid font must fail even though the color is valid.  This
	// catches the old GetColor-as-font-data coupling.
	opt.fontData = []byte("not a font")
	if _, err := r.Render(opt); err == nil {
		t.Fatal("expected invalid explicit font data error")
	}
}

func TestTextRendererRenderOnImageHonorsNonZeroOrigin(t *testing.T) {
	fontData, err := assets.GetDefaultFont()
	if err != nil {
		t.Fatalf("GetDefaultFont() error = %v", err)
	}
	r := NewTextRenderer(cache.NewResourceManager(1, 2, 1))
	img := image.NewRGBA(image.Rect(10, 20, 130, 80))
	opt := testTextOption{
		text: "A", maxWidth: 120, maxHeight: 60,
		fontData: fontData, fontColor: color.White,
	}
	if err := r.RenderOnImage(img, opt); err != nil {
		t.Fatalf("RenderOnImage() error = %v", err)
	}
	if !hasAnyNonTransparent(img) {
		t.Fatal("RenderOnImage() did not draw any pixels")
	}
}

func TestConcreteRenderOptionsProvideUsableDefaults(t *testing.T) {
	text := TextRenderOptions{Text: "hello", Width: 40, Height: 20}
	if err := text.ValidateOption(); err != nil {
		t.Fatalf("TextRenderOptions.ValidateOption() error = %v", err)
	}
	if width, height := text.GetMaxSize(); width != 40 || height != 20 {
		t.Fatalf("TextRenderOptions max size = %dx%d, want 40x20", width, height)
	}

	svg := SVGRenderOptions{SVGData: []byte("<svg/>"), Width: 2, Height: 3}
	if err := svg.ValidateOption(); err != nil {
		t.Fatalf("SVGRenderOptions.ValidateOption() error = %v", err)
	}
}

func sprintfSVG(fill string) string {
	// Keep formatting local to avoid pulling fmt into every test helper call.
	if fill == "#ff0000" {
		return `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 10 10"><rect x="0" y="0" width="10" height="10" fill="#ff0000"/></svg>`
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 10 10"><rect x="0" y="0" width="10" height="10" fill="#0000ff"/></svg>`
}

func hasColorNear(img image.Image, want color.RGBA) bool {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			got := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if got.A > 0 && absInt(int(got.R)-int(want.R)) < 8 && absInt(int(got.G)-int(want.G)) < 8 && absInt(int(got.B)-int(want.B)) < 8 {
				return true
			}
		}
	}
	return false
}

func hasAnyNonTransparent(img image.Image) bool {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0 {
				return true
			}
		}
	}
	return false
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
