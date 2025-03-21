# ICON MARKAR 

iconmarker supports attaching text to existing images and applying various filters

## Basic Usage

```go
package main

import "github.com/bagaking/iconmarker"

func main() {
    // ...
	return iconmarker.CreateImg(
		fontBytes,
		imgBytes,
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 200, G: 255, B: 255, A: 255},
			Text:      "Hello World",
		}.SetAdaptedSize(600, 300).AddOutline(color.RGBA{R: 16, G: 16, B: 16, A: 255}, 4),
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 210, G: 64, B: 32, A: 255},
			Text:      "iconmarker example",
			YOffset:   256,
		}.SetAdaptedSize(680, 80).AddShadow(color.RGBA{R: 128, G: 128, B: 128, A: 128}, ico.TitleShadowWidth),
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 64, G: 64, B: 45, A: 255},
			Text:      "from bagaking",
			YOffset:   320, 
		}.SetStaticSize(32),
	)
}
```

## Image Filters

The library includes a powerful filter system that allows you to apply various effects to your images:

### Available Filters

1. **Grayscale Filter** - Converts images to grayscale with optional alpha preservation
2. **Tint Filter** - Applies a color tint with adjustable intensity
3. **Opacity Filter** - Adjusts the transparency of images
4. **Invert Filter** - Inverts image colors with optional alpha inversion
5. **Composite Filter** - Combines multiple filters in sequence

### Using Filters

```go
package main

import (
    "image"
    "image/png"
    "os"
    
    "github.com/bagaking/iconmarker/filter"
)

func main() {
    // Load your image
    // ...
    
    // Create a filter manager
    filterManager := filter.NewFilterManager()
    
    // Apply a single filter
    grayImage, err := filterManager.QuickGrayscale(originalImage)
    
    // Apply a tint with custom color and intensity
    tintedImage, err := filterManager.QuickTint(originalImage, [3]uint8{255, 0, 0}, 0.7) // Red tint at 70% intensity
    
    // Apply multiple filters in sequence
    multiFilteredImage, err := filterManager.ApplyFilters(
        originalImage,
        []string{"grayscale", "tint"},
        []filter.FilterOption{
            filter.GrayscaleOption{PreserveAlpha: true},
            filter.TintOption{
                Color:     [3]uint8{0, 0, 255}, // Blue tint
                Intensity: 0.5,
            },
        },
    )
    
    // Save the result
    // ...
}
```

### Creating Custom Filters

You can create custom filters by implementing the `filter.Filter` interface:

```go
type MyCustomFilter struct{}

func NewMyCustomFilter() *MyCustomFilter {
    return &MyCustomFilter{}
}

func (f *MyCustomFilter) Apply(img draw.Image, options filter.FilterOption) error {
    // Implement your custom filter logic here
    return nil
}

// Register your custom filter
filterManager.Register("my-custom-filter", NewMyCustomFilter())
```

See the examples directory for more detailed usage examples.
