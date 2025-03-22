# IconMarker

IconMarker is a small Go library for generating labeled visual assets. It can
draw centered text on JPEG backgrounds, render embedded or custom SVG icons,
apply composable image filters, and compose reproducible badges or dashboards.

The compatibility API remains available alongside lower-level `core`,
`renderer`, `filter`, and `assets` packages.

Requires Go 1.23 or newer.

## Preview

The checked-in images below are documentation snapshots copied from verified
example output; they are not test baselines.

| Icon dashboard | Resolved badge |
| --- | --- |
| <img src="docs/assets/simple_icons_with_text.png" alt="Simple icons with text example output" width="360"> | <img src="docs/assets/badge_resolved.png" alt="Resolved badge example output" width="260"> |

## Install

```sh
go get github.com/bagaking/iconmarker
```

For the compatibility API:

```go
import "github.com/bagaking/iconmarker"
```

For lower-level composition:

```go
import (
    "github.com/bagaking/iconmarker/assets"
    "github.com/bagaking/iconmarker/core"
    "github.com/bagaking/iconmarker/filter"
    "github.com/bagaking/iconmarker/renderer"
)
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
    }.SetAdaptedSize(600, 120).AddOutline(color.RGBA{A: 255}, 2),
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

`DrawTextOption.SetStaticSize` uses a fixed font size. `SetAdaptedSize` scales
text to fit the requested maximum width and height. `AddShadow`, `AddOutline`,
and `MoveOffset` adjust text effects and placement.

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

Use `assets.AllIcons()` for the typed icon list,
`assets.ListAvailableIcons()` for names, or `assets.GetSVGIcon(name)` to load a
specific embedded icon.

## Render text directly

The renderer API keeps font bytes and text color as separate fields. Empty
`FontData` uses the embedded default font.

```go
textRenderer := renderer.NewTextRenderer(marker.GetResourceManager())
textImage, err := textRenderer.Render(renderer.TextRenderOptions{
    Text:      "Hello",
    Width:     320,
    Height:    100,
    FontColor: color.White,
    FontSize:  42,
})
```

## Filters

The default filter manager registers `grayscale`, `tint`, `opacity`, `invert`,
and `composite`. Filter application returns a new RGBA image and leaves the
source unchanged.

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

Tint intensity and opacity values must be between `0` and `1`; unknown names
return `filter.ErrFilterNotFound`.

## Reproducible examples

The repository examples are the command-line entry points for generating local
artifacts. Generated files are written below each example's `output/` directory
and are ignored by git.

```sh
git clone https://github.com/bagaking/iconmarker.git
cd iconmarker
go test ./...
bash examples/run_examples.sh
```

| Example | Generates | Reproduce one example |
| --- | --- | --- |
| `examples/basic_text` | centered text overlays using compatibility and core APIs | `cd examples/basic_text && go run main.go` |
| `examples/text_effects` | text with shadows, outlines, and combined effects | `cd examples/text_effects && go run main.go` |
| `examples/svg_rendering` | SVG rendering, resizing, and filtered SVG output | `cd examples/svg_rendering && go run main.go` |
| `examples/combined_filters` | composite, sequential, and custom filter outputs | `cd examples/combined_filters && go run main.go` |
| `examples/integrated_example` | older API, newer API, invert, and opacity outputs | `cd examples/integrated_example && go run main.go` |
| `examples/svg_with_text` | SVG/text layouts and embedded icon sheets | `cd examples/svg_with_text && go run main.go` |
| `examples/badge` | four status badge PNGs from SVG templates and text | `cd examples/badge && go run main.go` |
| `examples/icon_dashboard` | one icon dashboard PNG from embedded icons and labels | `cd examples/icon_dashboard && go run main.go` |

## Asset provenance

- Runtime embedded assets live under `assets/`: the default font and SVG icon
  set loaded by the `assets` package.
- Example-only inputs live under `examples/assets/`: `background.jpg`,
  `font.ttf`, and `icon.svg`.
- The `badge` example builds its SVG templates directly in Go.
- README preview snapshots live under `docs/assets/` and are copied from
  generated outputs after running the corresponding examples.

## Input and output boundaries

- Text rendering expects JPEG background bytes and returns an RGBA image.
- Filter APIs accept `image.Image` input and return a filtered image.
- `SaveImage2File` lets the caller choose an encoder such as PNG or JPEG.
- Generated example outputs are reproducible local artifacts, not golden files.

## Development and validation

Run the package tests, race detector, static analysis, and example smoke test:

```sh
go test ./... -count=1
go test -race ./...
go vet ./...
bash examples/run_examples.sh
git diff --check
```

After running examples, inspect `git status --short --ignored` to confirm that
only intended source or documentation files are staged.

## License

MIT. See [LICENSE](LICENSE).
