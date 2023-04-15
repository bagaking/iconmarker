package iconmarker

import "image/color"

type (
	FontEffect struct {
		Type    string      `json:"type"`
		Color   color.Color `json:"color"`
		XOffset int         `json:"x_offset"`
		YOffset int         `json:"y_offset"`
	}

	DrawTextOption struct {
		FontColor color.Color  `json:"font_color"`
		FontSize  float64      `json:"font_size"`
		MaxWidth  int          `json:"max_width"`
		MaxHeight int          `json:"max_height"`
		Text      string       `json:"text"`
		YOffset   int          `json:"y_offset"`
		XOffset   int          `json:"x_offset"`
		Effect    []FontEffect `json:"effect"`
	}
)

const (
	EShadow  = "shadow"
	EOutline = "outline"
)

func (o DrawTextOption) SetStaticSize(fontSize float64) DrawTextOption {
	o.FontSize = fontSize
	o.MaxWidth = 0
	o.MaxHeight = 0
	return o
}

func (o DrawTextOption) SetAdaptedSize(maxWidth, maxHeight int) DrawTextOption {
	o.FontSize = 0
	o.MaxWidth = maxWidth
	o.MaxHeight = maxHeight
	return o
}

func (o DrawTextOption) AddShadow(c color.Color, uniOffset int) DrawTextOption {
	o.Effect = append(o.Effect, FontEffect{
		Type:    EShadow,
		Color:   c,
		XOffset: uniOffset,
		YOffset: uniOffset,
	})
	return o
}

func (o DrawTextOption) AddOutline(c color.Color, uniOffset int) DrawTextOption {
	o.Effect = append(o.Effect, FontEffect{
		Type:    EOutline,
		Color:   c,
		XOffset: uniOffset,
		YOffset: uniOffset,
	})
	return o
}

func (o DrawTextOption) MoveOffset(x, y int) DrawTextOption {
	o.XOffset += x
	o.YOffset += y
	return o
}

func (o DrawTextOption) ToEffectGroup() []DrawTextOption {
	ret := make([]DrawTextOption, 0)
	for _, effect := range o.Effect {
		switch effect.Type {
		case EShadow:
			o2 := o.MoveOffset(effect.XOffset, effect.YOffset)
			o2.FontColor = effect.Color
			ret = append(ret, o2)

		case EOutline:
			// make a round by xoffset and yoffset
			for i := -effect.XOffset; i <= effect.XOffset; i++ {
				// Manufacturing rounded corners
				for j := -effect.YOffset; j <= effect.YOffset; j++ {
					if i*i+j*j <= effect.XOffset*effect.XOffset {

						o2 := o.MoveOffset(i, j)
						o2.FontColor = effect.Color
						ret = append(ret, o2)
					}
				}
			}
		}
	}
	return ret
}
