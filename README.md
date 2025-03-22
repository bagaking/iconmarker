# Iconmarker

Iconmarker is a Go library for drawing centered text onto JPEG images and
applying simple image filters. It exposes a compatibility package at
`github.com/bagaking/iconmarker` and lower-level packages under `core` and
`filter`.

## Preview

These checked-in preview assets are documentation snapshots copied from
runnable examples after running `bash examples/run_examples.sh`.
`docs/assets/simple_icons_with_text.png` comes from the
`examples/icon_dashboard` example's generated `simple_icons_with_text.png`
file, and `docs/assets/badge_resolved.png` comes from the `examples/badge`
example's generated `badge_resolved.png` file.

| Icon dashboard | Resolved badge |
| --- | --- |
| <img src="docs/assets/simple_icons_with_text.png" alt="Simple icons with text example output" width="360"> | <img src="docs/assets/badge_resolved.png" alt="Resolved badge example output" width="260"> |

## Installation

```sh
go get github.com/bagaking/iconmarker
```

Import the compatibility package for the common text-rendering API:

```go
import "github.com/bagaking/iconmarker"
```

Import the filter package when constructing filter options directly:

```go
import "github.com/bagaking/iconmarker/filter"
```

## Minimal Example

```go
package main

import (
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/bagaking/iconmarker"
)

func main() {
	backgroundBytes, err := os.ReadFile("background.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, err := iconmarker.CreateImg(
		nil,
		backgroundBytes,
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:      "Hello",
		}.SetStaticSize(48).AddOutline(color.RGBA{A: 255}, 2),
	)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		log.Fatal(err)
	}
}
```

## Examples Gallery

The repository includes runnable examples with lightweight source assets. They
cover text layout, text effects, SVG rendering, composed filters, badges, and
dashboard-style icon layouts.

| Example | Shows | Reproduce locally |
| --- | --- | --- |
| `examples/basic_text` | Legacy and current centered text APIs | `cd examples/basic_text && go run main.go` |
| `examples/text_effects` | Shadows, outlines, and combined text effects | `cd examples/text_effects && go run main.go` |
| `examples/svg_rendering` | SVG loading, resizing, and filter application | `cd examples/svg_rendering && go run main.go` |
| `examples/combined_filters` | Composite and sequential filter pipelines | `cd examples/combined_filters && go run main.go` |
| `examples/svg_with_text` | SVG and text layout combinations | `cd examples/svg_with_text && go run main.go` |
| `examples/badge` | Status badge generation from SVG, text, and color | `cd examples/badge && go run main.go` |
| `examples/icon_dashboard` | Embedded icons, filters, and labels in one image | `cd examples/icon_dashboard && go run main.go` |

Run all examples from the repository root:

```sh
bash examples/run_examples.sh
```

The script writes generated files under each example's `output/` directory and
returns a non-zero exit code if any example fails, so it is suitable for local
smoke testing and CI.

Documentation preview images live under `docs/assets/`. Per-example output
directories are generated artifacts and can be recreated with the commands
above.

## API Notes

- `CreateImg(fontBytes, backgroundBytes []byte, opts ...DrawTextOption)`
  decodes `backgroundBytes` as JPEG and returns an `*image.RGBA`.
- Passing `nil` or an empty slice for `fontBytes` uses the embedded default
  font. Passing bytes parses them as a TrueType font.
- `DrawTextOption.SetStaticSize` uses a fixed font size.
- `DrawTextOption.SetAdaptedSize` scales text to fit a maximum width and
  height.
- `AddShadow`, `AddOutline`, `XOffset`, and `YOffset` adjust centered text
  rendering.
- `SaveImage2File` writes an `image.Image` with the encoder supplied by the
  caller, such as `png.Encode` or a JPEG encoder wrapper.

## Filters

The default filter manager registers these filter names:

- `grayscale`
- `tint`
- `opacity`
- `invert`

Example:

```go
package main

import (
	"image"
	"log"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/filter"
)

func applyTint(src image.Image) image.Image {
	img, err := iconmarker.ApplyFilter(src, "tint", filter.TintOption{
		Color:     [3]uint8{255, 0, 0},
		Intensity: 0.4,
	})
	if err != nil {
		log.Fatal(err)
	}
	return img
}
```

Use `iconmarker.GetFilterManager()` or `filter.NewFilterManager()` when you need
to register a custom implementation of the `filter.Filter` interface.

## Input And Output Boundaries

- Text rendering currently expects JPEG background bytes.
- Text rendering returns an RGBA image in memory; output format is chosen by the
  caller's encoder.
- Filter APIs accept `image.Image` input and return a filtered image.
- Tint intensity and opacity values must be between `0` and `1`.
- Unknown filter names return `filter.ErrFilterNotFound`.

## Validation

Run the test suite and example smoke test with:

```sh
go test ./...
bash examples/run_examples.sh
```

`go test ./...` covers package-level behavior. `bash examples/run_examples.sh`
proves that the public examples still compile and can write their generated
outputs from the checked-in example assets.

## License

MIT. See [LICENSE](LICENSE).
