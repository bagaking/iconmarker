// Package renderer provides rendering capabilities for Icon Marker
package renderer

import (
	"fmt"
	"image"
	"image/color"
)

// RenderOption defines options for rendering operations
type RenderOption interface {
	// ValidateOption validates the rendering options
	ValidateOption() error
}

// Renderer defines the interface for all renderers
type Renderer interface {
	// Render performs the rendering operation with the given options
	Render(options RenderOption) (image.Image, error)
}

// TextRenderOption defines options for text rendering
type TextRenderOption interface {
	RenderOption
	// GetText returns the text to render
	GetText() string
	// GetMaxSize returns the maximum width and height for the text
	GetMaxSize() (width, height int)
	// GetFontSize returns the font size
	GetFontSize() float64
	// GetColor returns the text color
	GetColor() interface{}
	// GetPosition returns the x and y offsets for positioning
	GetPosition() (x, y int)
}

// FontDataOption is an optional extension implemented by text options that
// provide their own TrueType/OpenType bytes.  It intentionally remains
// separate from TextRenderOption so existing user-defined options continue to
// satisfy the original interface.
type FontDataOption interface {
	GetFontData() []byte
}

// TextRenderOptions is the ready-to-use implementation of TextRenderOption.
// Width and Height describe the canvas returned by Render.  MaxWidth and
// MaxHeight constrain the text size; when omitted, the canvas dimensions are
// used for fitting.  FontData is optional and falls back to the embedded font.
type TextRenderOptions struct {
	Text      string
	FontData  []byte
	FontColor color.Color
	FontSize  float64
	MaxWidth  int
	MaxHeight int
	Width     int
	Height    int
	XOffset   int
	YOffset   int
}

// TextOption is a concise compatibility alias for TextRenderOptions.
type TextOption = TextRenderOptions

func (o TextRenderOptions) ValidateOption() error {
	if o.Width < 0 || o.Height < 0 || o.MaxWidth < 0 || o.MaxHeight < 0 {
		return fmt.Errorf("text dimensions must not be negative")
	}
	if o.FontSize < 0 {
		return fmt.Errorf("font size must not be negative: %f", o.FontSize)
	}
	if o.Width == 0 && o.MaxWidth == 0 {
		return fmt.Errorf("text width must be positive")
	}
	if o.Height == 0 && o.MaxHeight == 0 {
		return fmt.Errorf("text height must be positive")
	}
	return nil
}

func (o TextRenderOptions) GetText() string { return o.Text }

func (o TextRenderOptions) GetMaxSize() (width, height int) {
	width, height = o.Width, o.Height
	if width == 0 {
		width = o.MaxWidth
	}
	if height == 0 {
		height = o.MaxHeight
	}
	return width, height
}

func (o TextRenderOptions) GetFontSize() float64 { return o.FontSize }

func (o TextRenderOptions) GetColor() interface{} { return o.FontColor }

func (o TextRenderOptions) GetPosition() (x, y int) { return o.XOffset, o.YOffset }

func (o TextRenderOptions) GetFontData() []byte {
	if len(o.FontData) == 0 {
		return nil
	}
	return append([]byte(nil), o.FontData...)
}

// SVGRenderOption defines options for SVG rendering
type SVGRenderOption interface {
	RenderOption
	// GetSVGData returns the SVG data
	GetSVGData() []byte
	// GetDimensions returns the desired width and height
	GetDimensions() (width, height int)
}

// SVGRenderOptions is the ready-to-use implementation of SVGRenderOption.
type SVGRenderOptions struct {
	SVGData []byte
	Width   int
	Height  int
}

// SVGOption is a concise compatibility alias for SVGRenderOptions.
type SVGOption = SVGRenderOptions

func (o SVGRenderOptions) ValidateOption() error {
	if len(o.SVGData) == 0 {
		return fmt.Errorf("SVG data is empty")
	}
	if o.Width <= 0 || o.Height <= 0 {
		return fmt.Errorf("invalid dimensions: width=%d, height=%d", o.Width, o.Height)
	}
	return nil
}

func (o SVGRenderOptions) GetSVGData() []byte {
	return append([]byte(nil), o.SVGData...)
}

func (o SVGRenderOptions) GetDimensions() (width, height int) {
	return o.Width, o.Height
}
