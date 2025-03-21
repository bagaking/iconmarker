// Package renderer provides rendering capabilities for Icon Marker
package renderer

import (
	"image"
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

// SVGRenderOption defines options for SVG rendering
type SVGRenderOption interface {
	RenderOption
	// GetSVGData returns the SVG data
	GetSVGData() []byte
	// GetDimensions returns the desired width and height
	GetDimensions() (width, height int)
}
