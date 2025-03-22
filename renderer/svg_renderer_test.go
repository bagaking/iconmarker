package renderer

import (
	"errors"
	"image"
	"strings"
	"testing"

	"github.com/bagaking/iconmarker/cache"
)

type stubSVGOption struct {
	svgData []byte
	width   int
	height  int
	err     error
}

func (o stubSVGOption) ValidateOption() error {
	return o.err
}

func (o stubSVGOption) GetSVGData() []byte {
	return o.svgData
}

func (o stubSVGOption) GetDimensions() (int, int) {
	return o.width, o.height
}

type stubRenderOption struct{}

func (o stubRenderOption) ValidateOption() error {
	return nil
}

func TestSVGRendererRenderRejectsInvalidOptions(t *testing.T) {
	renderer := NewSVGRenderer(cache.NewResourceManager(4, 4, 4))

	_, err := renderer.Render(stubRenderOption{})
	if err == nil || !strings.Contains(err.Error(), "options is not SVGRenderOption") {
		t.Fatalf("SVGRenderer.Render(stubRenderOption{}) error = %v, want SVGRenderOption type error", err)
	}

	validateErr := errors.New("invalid svg option")
	_, err = renderer.Render(stubSVGOption{err: validateErr})
	if !errors.Is(err, validateErr) {
		t.Fatalf("SVGRenderer.Render(option with validation error) error = %v, want %v", err, validateErr)
	}

	_, err = renderer.Render(stubSVGOption{
		svgData: []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`),
		width:   0,
		height:  10,
	})
	if err == nil || !strings.Contains(err.Error(), "invalid dimensions: width=0, height=10") {
		t.Fatalf("SVGRenderer.Render(width=0 height=10) error = %v, want invalid dimensions error", err)
	}

	_, err = renderer.Render(stubSVGOption{
		svgData: []byte(`<svg xmlns="http://www.w3.org/2000/svg"><path d="`),
		width:   10,
		height:  10,
	})
	if err == nil || !strings.Contains(err.Error(), "error parsing SVG") {
		t.Fatalf("SVGRenderer.Render(malformed svg) error = %v, want parse error", err)
	}
}

func TestSVGRendererRenderReturnsRequestedBounds(t *testing.T) {
	renderer := NewSVGRenderer(cache.NewResourceManager(4, 4, 4))
	opt := stubSVGOption{
		svgData: []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"><rect width="10" height="10" fill="red"/></svg>`),
		width:   12,
		height:  8,
	}

	got, err := renderer.Render(opt)
	if err != nil {
		t.Fatalf("SVGRenderer.Render(%dx%d valid svg) error = %v, want nil", opt.width, opt.height, err)
	}
	if got.Bounds() != image.Rect(0, 0, opt.width, opt.height) {
		t.Errorf("SVGRenderer.Render(%dx%d valid svg).Bounds() = %v, want %v", opt.width, opt.height, got.Bounds(), image.Rect(0, 0, opt.width, opt.height))
	}
}

func TestSVGRendererRenderCachesIndependentSVGData(t *testing.T) {
	renderer := NewSVGRenderer(cache.NewResourceManager(4, 4, 4))
	svgText := `<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"><rect width="10" height="10" fill="red"/></svg>`
	svgData := []byte(svgText)
	opt := stubSVGOption{
		svgData: svgData,
		width:   10,
		height:  10,
	}

	if _, err := renderer.Render(opt); err != nil {
		t.Fatalf("SVGRenderer.Render(%dx%d valid svg before caller mutation) error = %v, want nil", opt.width, opt.height, err)
	}

	for i := range svgData {
		svgData[i] = 0
	}

	opt.svgData = []byte(svgText)
	got, err := renderer.Render(opt)
	if err != nil {
		t.Fatalf("SVGRenderer.Render(%dx%d cached svg from fresh caller data after prior caller mutation) error = %v, want nil", opt.width, opt.height, err)
	}
	if got.Bounds() != image.Rect(0, 0, opt.width, opt.height) {
		t.Errorf("SVGRenderer.Render(%dx%d cached svg from fresh caller data after prior caller mutation).Bounds() = %v, want %v", opt.width, opt.height, got.Bounds(), image.Rect(0, 0, opt.width, opt.height))
	}
}

func TestSVGRendererRenderMultipleReportsIndexedError(t *testing.T) {
	renderer := NewSVGRenderer(cache.NewResourceManager(4, 4, 4))
	options := []SVGRenderOption{
		stubSVGOption{
			svgData: []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="4"></svg>`),
			width:   4,
			height:  4,
		},
		stubSVGOption{
			svgData: []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="4"></svg>`),
			width:   -1,
			height:  4,
		},
	}

	_, err := renderer.RenderMultiple(options)
	if err == nil || !strings.Contains(err.Error(), "error rendering SVG 1: invalid dimensions: width=-1, height=4") {
		t.Fatalf("SVGRenderer.RenderMultiple(options with invalid second item) error = %v, want indexed dimensions error", err)
	}
}
