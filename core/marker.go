// Package core contains the core functionality of the iconmarker
package core

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"time"

	"github.com/bagaking/iconmarker/assets"
	"github.com/bagaking/iconmarker/cache"
	"github.com/bagaking/iconmarker/filter"
	"github.com/bagaking/iconmarker/renderer"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// IconMarker 提供图标标记功能的主要结构
type IconMarker struct {
	resourceManager *cache.ResourceManager
	filterManager   *filter.FilterManager
	textRenderer    *renderer.TextRenderer
	svgRenderer     *renderer.SVGRenderer
}

// NewIconMarker 创建一个新的图标标记器
func NewIconMarker() *IconMarker {
	// 创建资源管理器，适当的缓存大小
	resourceManager := cache.NewResourceManager(100, 50, 200) // SVG, Font, Image caches
	resourceManager.SetTTL(30 * time.Minute)                  // 设置缓存生存时间

	// 创建滤镜管理器
	filterManager := filter.NewFilterManager()

	// 创建渲染器
	textRenderer := renderer.NewTextRenderer(resourceManager)
	svgRenderer := renderer.NewSVGRenderer(resourceManager)

	return &IconMarker{
		resourceManager: resourceManager,
		filterManager:   filterManager,
		textRenderer:    textRenderer,
		svgRenderer:     svgRenderer,
	}
}

// CreateImg 创建带有文本的图像（兼容旧API）
func (im *IconMarker) CreateImg(fontBytes, backgroundBytes []byte, drawFontOpt ...DrawTextOption) (*image.RGBA, error) {
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

	// 解析背景图片
	img, err := jpeg.Decode(bytes.NewReader(backgroundBytes))
	if err != nil {
		return nil, fmt.Errorf("%w, error decoding background", err)
	}

	outI := image.NewRGBA(img.Bounds())
	draw.Draw(outI, outI.Bounds(), img, img.Bounds().Min, draw.Src)

	// 在图像上绘制文本
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

// CreateImgWithFilters 创建带有文本和滤镜的图像
func (im *IconMarker) CreateImgWithFilters(fontBytes, backgroundBytes []byte,
	filters []string, filterOptions []filter.FilterOption,
	drawFontOpt ...DrawTextOption) (*image.RGBA, error) {

	// 先创建基本图像
	img, err := im.CreateImg(fontBytes, backgroundBytes, drawFontOpt...)
	if err != nil {
		return nil, err
	}

	// 应用滤镜
	if len(filters) > 0 {
		filteredImg, err := im.filterManager.ApplyFilters(img, filters, filterOptions)
		if err != nil {
			return nil, fmt.Errorf("error applying filters: %w", err)
		}

		// 将filteredImg转换为RGBA
		if rgba, ok := filteredImg.(*image.RGBA); ok {
			return rgba, nil
		}

		// 如果不是RGBA，转换它
		resultImg := image.NewRGBA(filteredImg.Bounds())
		draw.Draw(resultImg, resultImg.Bounds(), filteredImg, image.Point{}, draw.Src)
		return resultImg, nil
	}

	return img, nil
}

// SaveImage2File 将图像保存到文件
func (im *IconMarker) SaveImage2File(img image.Image, path string, encoder func(io.Writer, image.Image) error) error {
	return SaveImage2File(img, path, encoder)
}

// GetFilterManager 返回滤镜管理器，允许注册自定义滤镜
func (im *IconMarker) GetFilterManager() *filter.FilterManager {
	return im.filterManager
}

// GetResourceManager 返回资源管理器
func (im *IconMarker) GetResourceManager() *cache.ResourceManager {
	return im.resourceManager
}

// ApplyFilter 对图像应用单个滤镜
func (im *IconMarker) ApplyFilter(img image.Image, filterName string, option filter.FilterOption) (image.Image, error) {
	// 创建一个新的RGBA图像
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	// 应用滤镜
	if err := im.filterManager.Apply(dst, filterName, option); err != nil {
		return nil, err
	}

	return dst, nil
}

// ApplyFilters 对图像应用多个滤镜
func (im *IconMarker) ApplyFilters(img image.Image, filterNames []string, options []filter.FilterOption) (image.Image, error) {
	return im.filterManager.ApplyFilters(img, filterNames, options)
}
