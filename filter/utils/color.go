package utils

import (
	"fmt"
	"image/color"
	"strings"
)

// ParseHexColor 将十六进制颜色字符串解析为 color.RGBA
// 支持格式: #RGB, #RGBA, RGB, RGBA
func ParseHexColor(colorStr string) (color.RGBA, error) {
	// 移除#号
	colorStr = strings.TrimPrefix(colorStr, "#")

	// 确保字符串长度为6或8
	if len(colorStr) != 6 && len(colorStr) != 8 {
		return color.RGBA{}, fmt.Errorf("invalid color format: %s", colorStr)
	}

	// 解析RGB值
	var r, g, b, a uint8
	if len(colorStr) == 6 {
		_, err := fmt.Sscanf(colorStr, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("failed to parse RGB values: %v", err)
		}
		a = 255 // 默认不透明
	} else {
		_, err := fmt.Sscanf(colorStr, "%02x%02x%02x%02x", &r, &g, &b, &a)
		if err != nil {
			return color.RGBA{}, fmt.Errorf("failed to parse RGBA values: %v", err)
		}
	}

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}

// LerpColor 在两个颜色之间进行线性插值
// t 是插值因子，范围 [0,1]
func LerpColor(c1, c2 color.RGBA, t float64) color.RGBA {
	if t <= 0 {
		return c1
	}
	if t >= 1 {
		return c2
	}

	return color.RGBA{
		R: uint8(float64(c1.R) + t*float64(c2.R-c1.R)),
		G: uint8(float64(c1.G) + t*float64(c2.G-c1.G)),
		B: uint8(float64(c1.B) + t*float64(c2.B-c1.B)),
		A: uint8(float64(c1.A) + t*float64(c2.A-c1.A)),
	}
}

// DarkenColor 将颜色变暗
// factor 是暗化因子，范围 [0,1]
func DarkenColor(c color.RGBA, factor float64) color.RGBA {
	if factor <= 0 {
		return c
	}
	if factor >= 1 {
		return color.RGBA{A: c.A}
	}

	return color.RGBA{
		R: uint8(float64(c.R) * (1 - factor)),
		G: uint8(float64(c.G) * (1 - factor)),
		B: uint8(float64(c.B) * (1 - factor)),
		A: c.A,
	}
}

// LightenColor 将颜色变亮
// factor 是亮化因子，范围 [0,1]
func LightenColor(c color.RGBA, factor float64) color.RGBA {
	if factor <= 0 {
		return c
	}
	if factor >= 1 {
		return color.RGBA{R: 255, G: 255, B: 255, A: c.A}
	}

	return color.RGBA{
		R: uint8(float64(c.R) + factor*(255-float64(c.R))),
		G: uint8(float64(c.G) + factor*(255-float64(c.G))),
		B: uint8(float64(c.B) + factor*(255-float64(c.B))),
		A: c.A,
	}
}

// AdjustOpacity 调整颜色的透明度
// opacity 是新的透明度值，范围 [0,1]
func AdjustOpacity(c color.RGBA, opacity float64) color.RGBA {
	if opacity <= 0 {
		return color.RGBA{A: 0}
	}
	if opacity >= 1 {
		return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255}
	}

	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: uint8(opacity * 255),
	}
}

// ToHexString 将 color.RGBA 转换为十六进制颜色字符串
func ToHexString(c color.RGBA) string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}
