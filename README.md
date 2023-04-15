# ICON MARKAR 

iconmarker supports attaching text to existing images

## Usage

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
			MaxWidth:  600,
			MaxHeight: 300,
		}.AddOutline(color.RGBA{R: 16, G: 16, B: 16, A: 255}, 4),
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 210, G: 64, B: 32, A: 255},
			Text:      "iconmarker example",
			YOffset:   256,
			MaxWidth:  670,
			MaxHeight: 60,
		}.AddShadow(color.RGBA{R: 128, G: 128, B: 128, A: 128}, ico.TitleShadowWidth),
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 64, G: 64, B: 45, A: 255},
			FontSize:  32,
			Text:      "from bagaking",
			YOffset:   320,
			MaxWidth:  500,
			MaxHeight: 50,
		},
	)
}
```
