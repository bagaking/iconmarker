package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/bagaking/iconmarker/assets"
	"github.com/bagaking/iconmarker/cache"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// FontResource represents a cacheable font resource.
type FontResource struct {
	font *truetype.Font
}

// Size implements cache.CacheItem.
func (r *FontResource) Size() int {
	// Approximation of font size in memory.
	return 100 * 1024
}

// Clone implements cache.Resource.  truetype.Font is immutable after parsing,
// so sharing the parsed value is safe and avoids an expensive deep copy.
func (r *FontResource) Clone() cache.Resource { return r }

// TextRenderer implements the Renderer interface for text rendering.
type TextRenderer struct {
	resourceManager *cache.ResourceManager
}

// NewTextRenderer creates a new text renderer.  A nil manager is replaced by a
// private default manager so a renderer is safe to construct in isolation.
func NewTextRenderer(resourceManager *cache.ResourceManager) *TextRenderer {
	if resourceManager == nil {
		resourceManager = cache.NewResourceManager(1, 16, 1)
	}
	return &TextRenderer{resourceManager: resourceManager}
}

// Render renders text on a new RGBA image.
func (r *TextRenderer) Render(options RenderOption) (image.Image, error) {
	textOptions, ok := options.(TextRenderOption)
	if !ok {
		return nil, fmt.Errorf("options is not TextRenderOption")
	}
	if err := textOptions.ValidateOption(); err != nil {
		return nil, err
	}

	width, height := textOptions.GetMaxSize()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid text canvas dimensions: width=%d, height=%d", width, height)
	}

	fontLoaded, err := r.getFontForOption(textOptions)
	if err != nil {
		return nil, err
	}
	fontSize := r.fontSizeForOption(fontLoaded, textOptions, width, height)
	if fontSize < 1 {
		return nil, fmt.Errorf("invalid font size: %f", fontSize)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	if err := drawText(img, img.Bounds(), textOptions, fontLoaded, fontSize); err != nil {
		return nil, err
	}
	return img, nil
}

// RenderOnImage renders text on an existing image.  The image's full bounds,
// including a non-zero origin, are respected when centering the text.
func (r *TextRenderer) RenderOnImage(img draw.Image, options RenderOption) error {
	if img == nil {
		return fmt.Errorf("destination image is nil")
	}
	textOptions, ok := options.(TextRenderOption)
	if !ok {
		return fmt.Errorf("options is not TextRenderOption")
	}
	if err := textOptions.ValidateOption(); err != nil {
		return err
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return fmt.Errorf("destination image has empty bounds: %v", bounds)
	}
	fontLoaded, err := r.getFontForOption(textOptions)
	if err != nil {
		return err
	}

	maxWidth, maxHeight := textOptions.GetMaxSize()
	if maxWidth <= 0 {
		maxWidth = bounds.Dx()
	}
	if maxHeight <= 0 {
		maxHeight = bounds.Dy()
	}
	fontSize := r.fontSizeForOptionWithConstraints(fontLoaded, textOptions, maxWidth, maxHeight)
	if fontSize < 1 {
		return fmt.Errorf("invalid font size: %f", fontSize)
	}
	return drawText(img, bounds, textOptions, fontLoaded, fontSize)
}

func drawText(dst draw.Image, bounds image.Rectangle, options TextRenderOption, f *truetype.Font, fontSize float64) error {
	textColor := colorForOption(options)
	face := truetype.NewFace(f, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	defer face.Close()

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	text := options.GetText()
	textWidth := drawer.MeasureString(text).Round()
	textHeight := face.Metrics().Height.Ceil()
	xOffset, yOffset := options.GetPosition()

	// Include bounds.Min so sub-images (for example image.Rect(10, 20, ...))
	// are centered in their actual coordinate space rather than shifted toward
	// the origin.
	x := bounds.Min.X + ((bounds.Dx() - textWidth) >> 1) + xOffset
	y := bounds.Min.Y + ((bounds.Dy() + textHeight) >> 1) - face.Metrics().Descent.Round() + yOffset
	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(text)
	return nil
}

func colorForOption(options TextRenderOption) color.Color {
	if textColor, ok := options.GetColor().(color.Color); ok && textColor != nil {
		return textColor
	}
	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

func (r *TextRenderer) getFontForOption(options TextRenderOption) (*truetype.Font, error) {
	fontData, err := fontDataForOption(options)
	if err != nil {
		return nil, err
	}
	return r.getFont(fontData)
}

func fontDataForOption(options TextRenderOption) ([]byte, error) {
	// New options expose font bytes explicitly.  This keeps GetColor dedicated
	// to color and makes invalid font data observable instead of silently
	// falling back to the embedded font.
	if provider, ok := options.(FontDataOption); ok {
		if data := provider.GetFontData(); len(data) > 0 {
			return append([]byte(nil), data...), nil
		}
	}

	// Preserve the pre-existing, undocumented []byte convention for callers
	// that implemented the old interface.  It is only a fallback; concrete
	// options and new implementations should use FontDataOption.
	if legacyData, ok := options.GetColor().([]byte); ok && len(legacyData) > 0 {
		return append([]byte(nil), legacyData...), nil
	}

	fontData, err := assets.GetDefaultFont()
	if err != nil {
		return nil, fmt.Errorf("failed to load default font: %w", err)
	}
	return fontData, nil
}

// getFont loads a font from cache or parses it.
func (r *TextRenderer) getFont(fontData []byte) (*truetype.Font, error) {
	key := r.resourceManager.GenerateKeyFromData(fontData)
	if item, found := r.resourceManager.GetResource("font", key, r.resourceManager.GetFontCache()); found {
		if fontResource, ok := item.(*FontResource); ok && fontResource.font != nil {
			return fontResource.font, nil
		}
	}

	f, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}
	r.resourceManager.PutResource("font", key, r.resourceManager.GetFontCache(), &FontResource{font: f})
	return f, nil
}

func (r *TextRenderer) fontSizeForOption(f *truetype.Font, options TextRenderOption, width, height int) float64 {
	return r.fontSizeForOptionWithConstraints(f, options, width, height)
}

func (r *TextRenderer) fontSizeForOptionWithConstraints(f *truetype.Font, options TextRenderOption, maxWidth, maxHeight int) float64 {
	requested := options.GetFontSize()
	if requested > 0 {
		// Explicit sizes are honored unless the option's bounds require scaling
		// down to keep the text inside the canvas.
		return r.shrinkToFit(f, options.GetText(), maxWidth, maxHeight, requested)
	}
	return r.adaptFontSize(f, options.GetText(), maxWidth, maxHeight, 0)
}

func (r *TextRenderer) shrinkToFit(f *truetype.Font, text string, maxW, maxH int, requested float64) float64 {
	if requested < 1 {
		requested = 1
	}
	fits := makeTextFits(f, text, maxW, maxH)
	if fits(requested) {
		return requested
	}
	low, high := 1.0, requested
	if !fits(low) {
		return low
	}
	for i := 0; i < 32; i++ {
		mid := (low + high) / 2
		if fits(mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

// adaptFontSize calculates a fitting size and then chooses the largest size
// that satisfies the supplied limits.  A one-pixel minimum guarantees
// termination even when the text cannot fit the requested bounds.
func (r *TextRenderer) adaptFontSize(f *truetype.Font, text string, maxW, maxH int, fontSize float64) float64 {
	if maxW <= 0 && maxH <= 0 {
		if fontSize > 0 {
			return fontSize
		}
		return 1
	}

	requested := fontSize
	if requested <= 0 {
		if maxH > 0 {
			requested = float64(maxH)
		} else {
			requested = float64(maxW)
		}
	}
	if requested < 1 {
		requested = 1
	}

	fits := makeTextFits(f, text, maxW, maxH)
	if !fits(requested) {
		if !fits(1) {
			return 1
		}
		low, high := 1.0, requested
		for i := 0; i < 32; i++ {
			mid := (low + high) / 2
			if fits(mid) {
				low = mid
			} else {
				high = mid
			}
		}
		return low
	}

	// Grow an automatically selected size until the first non-fitting bound,
	// then binary-search the final interval.
	low, high := requested, requested*2
	for i := 0; i < 24 && fits(high); i++ {
		low, high = high, high*2
	}
	for i := 0; i < 32; i++ {
		mid := (low + high) / 2
		if fits(mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func makeTextFits(f *truetype.Font, text string, maxW, maxH int) func(float64) bool {
	return func(size float64) bool {
		face := truetype.NewFace(f, &truetype.Options{Size: size, DPI: 72, Hinting: font.HintingNone})
		defer face.Close()
		width := (&font.Drawer{Face: face}).MeasureString(text).Round()
		if maxW > 0 && width > maxW {
			return false
		}
		return maxH <= 0 || face.Metrics().Height.Ceil() <= maxH
	}
}
