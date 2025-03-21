package filter

import (
	"image/color"
	"image/draw"
	"math"
)

// TintFilter applies a color tint to an image
type TintFilter struct{}

// NewTintFilter creates a new tint filter
func NewTintFilter() *TintFilter {
	return &TintFilter{}
}

// Apply applies the tint filter
func (f *TintFilter) Apply(img draw.Image, options FilterOption) error {
	// Cast options to TintOption
	opt, ok := options.(TintOption)
	if !ok {
		return ErrInvalidColor
	}

	// Validate options
	if err := opt.ValidateOption(); err != nil {
		return err
	}

	// Extract tint color components
	tintR, tintG, tintB := opt.Color[0], opt.Color[1], opt.Color[2]
	intensity := opt.Intensity

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

			// 如果像素完全透明，直接保留原样不处理
			if a8 == 0 {
				continue
			}

			// 根据不透明度调整tint影响程度
			// 半透明像素应该比完全不透明的像素受到更少的tint影响
			actualIntensity := intensity * float64(a8) / 255.0

			// Calculate luminance from original RGB values
			lum := float64(r8)*0.299 + float64(g8)*0.587 + float64(b8)*0.114

			// Apply tint based on luminance and intensity, taking alpha into account
			// 对于RGB值，我们混合原始亮度和目标颜色
			newR := uint8(math.Round(lum*(1-actualIntensity) + float64(tintR)*actualIntensity))
			newG := uint8(math.Round(lum*(1-actualIntensity) + float64(tintG)*actualIntensity))
			newB := uint8(math.Round(lum*(1-actualIntensity) + float64(tintB)*actualIntensity))

			// 设置新颜色，保留原始alpha值
			img.Set(x, y, color.RGBA{newR, newG, newB, a8})
		}
	}

	return nil
}
