package filter

import (
	"image/color"
	"image/draw"
)

// OpacityFilter adjusts the opacity of an image
type OpacityFilter struct{}

// NewOpacityFilter creates a new opacity filter
func NewOpacityFilter() *OpacityFilter {
	return &OpacityFilter{}
}

// Apply applies the opacity filter
func (f *OpacityFilter) Apply(img draw.Image, options FilterOption) error {
	// Cast options to OpacityOption
	opt, ok := options.(OpacityOption)
	if !ok {
		// Use default options if not provided
		opt = OpacityOption{
			Opacity: 1.0, // Default to fully opaque
		}
	}

	// Validate options
	if err := opt.ValidateOption(); err != nil {
		return err
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

			// Apply opacity
			newA := uint8(float64(a8) * opt.Opacity)

			img.Set(x, y, color.RGBA{r8, g8, b8, newA})
		}
	}

	return nil
}
