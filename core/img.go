package core

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"

	"github.com/bagaking/iconmarker/assets"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// CreateImg 是旧API的兼容函数，未来应该使用 IconMarker.CreateImg
func CreateImg(fontBytes, backgroundBytes []byte, drawFontOpt ...DrawTextOption) (*image.RGBA, error) {
	var font *truetype.Font
	var err error

	// 如果字体字节数组为空，尝试加载默认字体
	if len(fontBytes) == 0 {
		defaultFontBytes, err := assets.GetDefaultFont()
		if err != nil {
			return nil, fmt.Errorf("failed to get default font: %w", err)
		}
		font, err = freetype.ParseFont(defaultFontBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default font: %w", err)
		}
	} else {
		// 解析用户提供的字体
		font, err = freetype.ParseFont(fontBytes)
		if err != nil {
			return nil, fmt.Errorf("%w, error parsing font file", err)
		}
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
func adaptSize(f *truetype.Font, text string, maxW, maxH int, fontSize float64) (realSize float64) {
	if maxH > 0 {
		realSize = float64(maxH)
	} else {
		realSize = float64(maxW)
	}
	opt := truetype.Options{
		Size:    realSize,
		DPI:     72,
		Hinting: font.HintingNone,
	}
	d := &font.Drawer{
		Face: truetype.NewFace(f, &opt),
	}

	isInside := func() bool {
		width := d.MeasureString(text).Round()
		return width <= maxW && d.Face.Metrics().Height.Ceil() <= maxH
	}

	changeFontSize := func(size float64) {
		opt.Size = size
		d.Face = truetype.NewFace(f, &opt)
	}

	// if font size is specified, use it directly when it is smaller than max width
	if fontSize > 0 {
		if isInside() {
			return fontSize
		}
		// if realSize > fontSize, init realSize to fontSize
		// cuz the fontsize is already larger than max, it will make the binary
		// search faster
		// but is makes different result from different font size, thus the adapting
		// result is stable only when font size is not specified (fontSize == 0)
		if realSize > fontSize {
			realSize = fontSize
		}
	}

	// binary search to reduce font size until it is smaller than max width
	for {
		if isInside() {
			break
		}

		realSize /= 2
		changeFontSize(realSize)

		if realSize < 3 {
			realSize = 3
			break
		}
	}

	// font size is already smaller than max width, now we need to enlarge font
	// size to make it just smaller than max width
	for {
		assumedSize := (realSize + 1) * 1.1
		changeFontSize(assumedSize)

		if !isInside() {
			break
		}
		realSize = assumedSize
	}

	return realSize
}

// DrawCenteredFont draws text on image with center alignment
// if opt.MaxWidth > 0 and opt.fontsize == 0, the font size will
// be adapted to fit the max width and height (height is ignored
// if it is 0), otherwise the font size will be opt.FontSize. for
// any specified font size, the adapting result will be stable
//
// if opt.FontSize > 0, the font size will be used directly and be
// scaled down to fit max width
//
// when adapt font size are not used, the smaller font size is 1,
// if the font size is smaller than 1, an error will be returned
//
// to easily draw text with different effects, use DrawTextOption's
// pipe operators, such as DrawTextOption.SetStaticSize or
// DrawTextOption.SetAdaptedSize
//
// its also possible to draw text with different effects,
// see DrawTextOption
func DrawCenteredFont(f *truetype.Font, outI *image.RGBA, opt DrawTextOption) error {
	// 如果传入的字体为nil，尝试加载默认字体
	if f == nil {
		fontData, err := assets.GetDefaultFont()
		if err != nil {
			return fmt.Errorf("failed to get default font: %w", err)
		}

		var parseErr error
		f, parseErr = freetype.ParseFont(fontData)
		if parseErr != nil {
			return fmt.Errorf("failed to parse default font: %w", parseErr)
		}
	}

	if opt.MaxWidth > 0 {
		opt.FontSize = adaptSize(f, opt.Text, opt.MaxWidth, opt.MaxHeight, opt.FontSize)
	} else if opt.FontSize < 1 {
		return fmt.Errorf("invalid font size: %f", opt.FontSize)
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
