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

// FontResource represents a cacheable font resource
type FontResource struct {
	font *truetype.Font
}

// Size implements cache.CacheItem
func (r *FontResource) Size() int {
	// Approximation of font size in memory
	return 100 * 1024 // 100KB is a reasonable approximation
}

// Clone implements cache.Resource
func (r *FontResource) Clone() cache.Resource {
	// Fonts are immutable, so we can return the same instance
	return r
}

// TextRenderer implements the Renderer interface for text rendering
type TextRenderer struct {
	resourceManager *cache.ResourceManager
}

// NewTextRenderer creates a new text renderer
func NewTextRenderer(resourceManager *cache.ResourceManager) *TextRenderer {
	return &TextRenderer{
		resourceManager: resourceManager,
	}
}

// Render renders text on an image
func (r *TextRenderer) Render(options RenderOption) (image.Image, error) {
	// Cast options to TextRenderOption
	textOptions, ok := options.(TextRenderOption)
	if !ok {
		return nil, fmt.Errorf("options is not TextRenderOption")
	}

	// Validate options
	if err := textOptions.ValidateOption(); err != nil {
		return nil, err
	}

	// Create a new RGBA image
	width, height := textOptions.GetMaxSize()
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Get the font
	var fontData []byte
	var err error

	fontColor, ok := textOptions.GetColor().([]byte)
	if ok {
		fontData = fontColor
	} else {
		// 使用默认字体数据作为兜底
		fontData, err = assets.GetDefaultFont()
		if err != nil {
			return nil, fmt.Errorf("failed to get font data and default font: %w", err)
		}
	}

	fontLoaded, err := r.getFont(fontData)
	if err != nil {
		return nil, err
	}

	// Calculate font size and position
	fontSize := textOptions.GetFontSize()
	if fontSize <= 0 {
		fontSize = r.adaptFontSize(fontLoaded, textOptions.GetText(), width, height, 0)
	}

	// Create drawer
	textColor, ok := textOptions.GetColor().(color.Color)
	if !ok {
		textColor = color.RGBA{255, 255, 255, 255} // Default to white
	}

	face := truetype.NewFace(fontLoaded, &truetype.Options{
		Size: fontSize,
		DPI:  72,
	})

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	// Calculate text position
	xOffset, yOffset := textOptions.GetPosition()
	txtWidth := drawer.MeasureString(textOptions.GetText()).Round()
	txtHeight := int(fontSize)

	// Center text horizontally and vertically
	x := (width-txtWidth)/2 + xOffset
	y := (height+txtHeight)/2 - drawer.Face.Metrics().Descent.Round() + yOffset

	// Draw text
	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(textOptions.GetText())

	return img, nil
}

// RenderOnImage renders text on an existing image
func (r *TextRenderer) RenderOnImage(img draw.Image, options RenderOption) error {
	// Cast options to TextRenderOption
	textOptions, ok := options.(TextRenderOption)
	if !ok {
		return fmt.Errorf("options is not TextRenderOption")
	}

	// Validate options
	if err := textOptions.ValidateOption(); err != nil {
		return err
	}

	// Get the font
	var fontData []byte
	var err error

	fontColor, ok := textOptions.GetColor().([]byte)
	if ok {
		fontData = fontColor
	} else {
		// 使用默认字体数据作为兜底
		fontData, err = assets.GetDefaultFont()
		if err != nil {
			return fmt.Errorf("failed to get font data and default font: %w", err)
		}
	}

	fontLoaded, err := r.getFont(fontData)
	if err != nil {
		return err
	}

	// Get image bounds
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calculate font size and position
	fontSize := textOptions.GetFontSize()
	maxWidth, maxHeight := textOptions.GetMaxSize()
	if maxWidth <= 0 {
		maxWidth = width
	}
	if maxHeight <= 0 {
		maxHeight = height
	}

	if fontSize <= 0 {
		fontSize = r.adaptFontSize(fontLoaded, textOptions.GetText(), maxWidth, maxHeight, 0)
	}

	// Create drawer
	textColor, ok := textOptions.GetColor().(color.Color)
	if !ok {
		textColor = color.RGBA{255, 255, 255, 255} // Default to white
	}

	face := truetype.NewFace(fontLoaded, &truetype.Options{
		Size: fontSize,
		DPI:  72,
	})

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	// Calculate text position
	xOffset, yOffset := textOptions.GetPosition()
	txtWidth := drawer.MeasureString(textOptions.GetText()).Round()
	txtHeight := int(fontSize)

	// Center text horizontally and vertically
	x := (width-txtWidth)/2 + xOffset
	y := (height+txtHeight)/2 - drawer.Face.Metrics().Descent.Round() + yOffset

	// Draw text
	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(textOptions.GetText())

	return nil
}

// getFont loads a font from cache or parses it
func (r *TextRenderer) getFont(fontData []byte) (*truetype.Font, error) {
	// Generate key for font cache
	key := r.resourceManager.GenerateKeyFromData(fontData)

	// Try to get from cache
	item, found := r.resourceManager.GetResource("font", key, r.resourceManager.GetFontCache())
	if found {
		if fontResource, ok := item.(*FontResource); ok {
			return fontResource.font, nil
		}
	}

	// Parse font
	font, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %w", err)
	}

	// Cache font
	r.resourceManager.PutResource("font", key, r.resourceManager.GetFontCache(), &FontResource{font: font})

	return font, nil
}

// adaptFontSize calculates the appropriate font size for the given text and dimensions
func (r *TextRenderer) adaptFontSize(f *truetype.Font, text string, maxW, maxH int, fontSize float64) float64 {
	if maxH > 0 {
		fontSize = float64(maxH)
	} else {
		fontSize = float64(maxW)
	}
	opt := truetype.Options{
		Size: fontSize,
		DPI:  72,
	}

	face := truetype.NewFace(f, &opt)
	drawer := &font.Drawer{
		Face: face,
	}

	isInside := func() bool {
		width := drawer.MeasureString(text).Round()
		return width <= maxW && face.Metrics().Height.Ceil() <= maxH
	}

	changeFontSize := func(size float64) {
		opt.Size = size
		newFace := truetype.NewFace(f, &opt)
		drawer.Face = newFace
		face = newFace
	}

	// Binary search to reduce font size until it fits
	for {
		if isInside() {
			break
		}

		fontSize /= 2
		changeFontSize(fontSize)

		if fontSize < 3 {
			fontSize = 3
			break
		}
	}

	// Enlarge font size to find the optimal size
	for {
		assumedSize := (fontSize + 1) * 1.1
		changeFontSize(assumedSize)

		if !isInside() {
			break
		}
		fontSize = assumedSize
	}

	return fontSize
}
