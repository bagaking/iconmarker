package filter

import (
	"image/color"
	"image/draw"
)

// GrayscaleFilter converts an image to grayscale
type GrayscaleFilter struct{}

// NewGrayscaleFilter creates a new grayscale filter
func NewGrayscaleFilter() *GrayscaleFilter {
	return &GrayscaleFilter{}
}

// Apply applies the grayscale filter
func (f *GrayscaleFilter) Apply(img draw.Image, options FilterOption) error {
	// Cast options to GrayscaleOption
	opt, ok := options.(GrayscaleOption)
	if !ok {
		// Use default options if not provided
		opt = GrayscaleOption{
			PreserveAlpha: true,
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

			// Calculate grayscale using ITU-R BT.709 coefficients
			gray := uint8((int(r8)*299 + int(g8)*587 + int(b8)*114) / 1000)

			if opt.PreserveAlpha {
				img.Set(x, y, color.RGBA{gray, gray, gray, a8})
			} else {
				img.Set(x, y, color.RGBA{gray, gray, gray, gray})
			}
		}
	}

	return nil
}
