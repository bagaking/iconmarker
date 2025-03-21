package filter

import (
	"image"
	"image/draw"
)

// ApplyFilters applies multiple filters to an image and returns a new image
func (fm *FilterManager) ApplyFilters(src image.Image, filterNames []string, optionsList []FilterOption) (image.Image, error) {
	// Create a new RGBA image to work with
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)

	// Apply each filter in sequence
	for i, name := range filterNames {
		filter, ok := fm.Get(name)
		if !ok {
			return nil, ErrFilterNotFound
		}

		var option FilterOption
		if i < len(optionsList) {
			option = optionsList[i]
		}

		if err := filter.Apply(dst, option); err != nil {
			return nil, err
		}
	}

	return dst, nil
}

// ApplyNamedFilters applies multiple filters to an image by name with default options
func (fm *FilterManager) ApplyNamedFilters(src image.Image, filterNames []string) (image.Image, error) {
	return fm.ApplyFilters(src, filterNames, nil)
}

// QuickGrayscale applies a grayscale filter with default options
func (fm *FilterManager) QuickGrayscale(src image.Image) (image.Image, error) {
	return fm.ApplyNamedFilters(src, []string{"grayscale"})
}

// QuickTint applies a tint filter with the specified color and intensity
func (fm *FilterManager) QuickTint(src image.Image, color [3]uint8, intensity float64) (image.Image, error) {
	option := TintOption{
		Color:     color,
		Intensity: intensity,
	}
	return fm.ApplyFilters(src, []string{"tint"}, []FilterOption{option})
}

// QuickOpacity applies an opacity filter with the specified opacity
func (fm *FilterManager) QuickOpacity(src image.Image, opacity float64) (image.Image, error) {
	option := OpacityOption{
		Opacity: opacity,
	}
	return fm.ApplyFilters(src, []string{"opacity"}, []FilterOption{option})
}

// QuickInvert applies an invert filter
func (fm *FilterManager) QuickInvert(src image.Image, invertAlpha bool) (image.Image, error) {
	option := InvertOption{
		InvertAlpha: invertAlpha,
	}
	return fm.ApplyFilters(src, []string{"invert"}, []FilterOption{option})
}
