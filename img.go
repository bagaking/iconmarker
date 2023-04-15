// Package iconmarker - supports attaching text to existing images
package iconmarker

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func CreateImg(fontBytes, backgroundBytes []byte, drawFontOpt ...DrawTextOption) (*image.RGBA, error) {
	// Parse font file
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("%w, error parsing font file", err)
	}

	// Parse bg image
	img, err := jpeg.Decode(bytes.NewReader(backgroundBytes))
	if err != nil {
		return nil, fmt.Errorf("%w, error decoding background", err)
	}

	outI := image.NewRGBA(img.Bounds())
	draw.Draw(outI, outI.Bounds(), img, img.Bounds().Min, draw.Src)

	// Draw text on image
	for _, opt := range drawFontOpt {
		for _, eop := range opt.ToEffectGroup() {
			err = DrawCenteredFont(font, outI, eop)
			if err != nil {
				return nil, fmt.Errorf("%w, error drawing text", err)
			}
		}

		if err = DrawCenteredFont(font, outI, opt); err != nil {
			return nil, fmt.Errorf("%w, error drawing text", err)
		}
	}

	return outI, nil
}

// SaveImage2File save image to file
//
//	imgEncoder: png.Encode or func(w io.Writer, m image.Image) error {
//		return jpeg.Encode(w, m, &jpeg.Options{Quality: 100})
//	}
//
// if you want to save image to byte stream, just call encoder outside
// e.g.
// var buf bytes.Buffer
// err = png.Encode(&buf, img)
func SaveImage2File(img image.Image, path string, encoder func(io.Writer, image.Image) error) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w, error creating output file", err)
	}
	defer outputFile.Close()

	// Encode image as JPEG
	err = encoder(outputFile, img)
	if err != nil {
		return fmt.Errorf("%w, error encoding image", err)
	}
	return nil
}

// adaptSize returns the real font size that fits the max width and height
func adaptSize(f *truetype.Font, text string, maxW, maxH int) (realSize float64) {
	if maxH > 0 {
		realSize = float64(maxH)
	} else {
		realSize = float64(maxW)
	}
	d := &font.Drawer{
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    realSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}

	// binary search to reduce font size until it is smaller than max width
	for {
		width := d.MeasureString(text).Round()
		if width <= maxW && d.Face.Metrics().Height.Ceil() <= maxH {
			break
		}
		realSize /= 2
		d.Face = truetype.NewFace(f, &truetype.Options{
			Size:    realSize,
			DPI:     72,
			Hinting: font.HintingNone,
		})

		if realSize < 3 {
			realSize = 3
			break
		}
	}

	// font size is already smaller than max width, now we need to enlarge font
	// size to make it just smaller than max width
	for {
		assumedSize := (realSize + 1) * 1.1
		d.Face = truetype.NewFace(f, &truetype.Options{
			Size:    assumedSize,
			DPI:     72,
			Hinting: font.HintingNone,
		})
		width := d.MeasureString(text).Round()
		if width > maxW || d.Face.Metrics().Height.Ceil() > maxH {
			break
		}
		realSize = assumedSize
	}

	return realSize
}

func DrawCenteredFont(f *truetype.Font, outI *image.RGBA, opt DrawTextOption) error {
	if opt.MaxWidth > 0 {
		opt.FontSize = adaptSize(f, opt.Text, opt.MaxWidth, opt.MaxHeight)
	}

	d := &font.Drawer{
		Dst: outI,
		Src: image.NewUniform(opt.FontColor),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    opt.FontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}

	// 计算
	size := d.MeasureString(opt.Text)
	width := size.Round()
	height := int(opt.FontSize)

	// 计算文本位置以居中绘制, 将文本的视觉中心点对齐到图像中心点
	// 计算 x 坐标，让文本中心对齐图像中心
	x := (outI.Bounds().Dx()-width)/2 + opt.XOffset
	// 计算 y 坐标，让文本 baseline 对齐图像中心
	// 这里的文本 baseline 是指文本底部距离 baseline 的距离，也就是 font.Metrics().Descent
	y := (outI.Bounds().Dy()+height)/2 - d.Face.Metrics().Descent.Round() + opt.YOffset

	// 设置文本位置
	d.Dot = fixed.P(x, y)

	// 绘制文本到图像
	d.DrawString(opt.Text)

	return nil
}
