# IconMarker

IconMarker is a small Go library for drawing centered text on JPEG backgrounds,
rendering embedded or custom SVG icons, and applying composable image filters.
It includes an embedded default font and 22 embedded SVG icons.

Requires Go 1.23 or newer.

## Install

```bash
go get github.com/bagaking/iconmarker
```

## Draw text

`CreateImg` keeps the original package-level API. Pass an empty font slice to
use the embedded default font; the background is JPEG data.

```go
img, err := iconmarker.CreateImg(
    nil,
    jpegBytes,
    iconmarker.DrawTextOption{
        Text:      "Hello, IconMarker",
        FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
    }.
        SetAdaptedSize(600, 120).
        AddOutline(color.RGBA{A: 255}, 2),
)
if err != nil {
    return err
}
if err := iconmarker.SaveImage2File(img, "output.png", png.Encode); err != nil {
    return err
}
```

For an instance with its own filters and caches:

```go
marker := iconmarker.NewIconMarker()
img, err := marker.CreateImgWithFilters(
    nil,
    jpegBytes,
    []string{"grayscale", "tint"},
    []filter.FilterOption{
        filter.GrayscaleOption{PreserveAlpha: true},
        filter.TintOption{Color: [3]uint8{40, 100, 255}, Intensity: 0.35},
    },
    iconmarker.DrawTextOption{Text: "TDD", FontColor: color.White}.SetStaticSize(48),
)
```

## Render SVG

```go
marker := iconmarker.NewIconMarker()
svgData, err := assets.IconDiamondMarker.Load()
if err != nil {
    return err
}

svgRenderer := renderer.NewSVGRenderer(marker.GetResourceManager())
svgImage, err := svgRenderer.Render(renderer.SVGRenderOptions{
    SVGData: svgData,
    Width:   128,
    Height:  128,
})
```

Use `assets.AllIcons()` for the typed icon list or
`assets.ListAvailableIcons()` for names.

## Render text directly

The renderer API keeps font bytes and text color as separate fields:

```go
textRenderer := renderer.NewTextRenderer(marker.GetResourceManager())
textImage, err := textRenderer.Render(renderer.TextRenderOptions{
    Text:      "Hello",
    Width:     320,
    Height:    100,
    FontColor: color.White,
    FontSize:  42,
    // FontData is optional; empty uses the embedded default font.
})
```

## Filters

Built-in names are `grayscale`, `tint`, `opacity`, `invert`, and `composite`.
Filter application returns a new RGBA image and leaves the source unchanged.

```go
manager := filter.NewFilterManager()
result, err := manager.ApplyFilters(
    source,
    []string{"grayscale", "opacity"},
    []filter.FilterOption{
        filter.GrayscaleOption{PreserveAlpha: true},
        filter.OpacityOption{Opacity: 0.8},
    },
)
```

Custom filters implement `filter.Filter` and are registered by name:

```go
type MyFilter struct{}

func (MyFilter) Apply(img draw.Image, option filter.FilterOption) error {
    // Mutate img or return an explanatory error.
    return nil
}

manager.Register("mine", MyFilter{})
```

## Development

```bash
go test ./... -count=1
go test -race ./...
go vet ./...
```

Runnable examples live in [`examples`](examples/README.md). They use the
fixtures in `examples/assets`; `examples/run_examples.sh` exits non-zero if any
example fails.
