package filter

import (
	"image/color"
	"image/draw"
)

// InvertFilter inverts the colors of an image
type InvertFilter struct{}

// NewInvertFilter creates a new invert filter
func NewInvertFilter() *InvertFilter {
	return &InvertFilter{}
}

// Apply applies the invert filter
func (f *InvertFilter) Apply(img draw.Image, options FilterOption) error {
	// Cast options to InvertOption
	opt, ok := options.(InvertOption)
	if !ok {
		// Use default options if not provided
		opt = InvertOption{
			InvertAlpha: false, // Default to not inverting alpha
		}
	}

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			// Convert to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Invert colors (255 - value)
			newR := 255 - r8
			newG := 255 - g8
			newB := 255 - b8

			// Optionally invert alpha
			newA := a8
			if opt.InvertAlpha {
				newA = 255 - a8
			}

			img.Set(x, y, color.RGBA{newR, newG, newB, newA})
		}
	}

	return nil
}
