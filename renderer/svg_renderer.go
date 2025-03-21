package renderer

import (
	"bytes"
	"fmt"
	"image"
	"sync"

	"github.com/bagaking/iconmarker/cache"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// SVGResource represents a cacheable SVG resource
// 存储原始 SVG 数据而不是已解析的图标，这样我们可以避免并发问题
type SVGResource struct {
	svgData []byte // 存储原始SVG数据，而不是解析后的图标
}

// Size implements cache.CacheItem
// 返回 SVG 数据的实际大小
func (r *SVGResource) Size() int {
	// 返回实际SVG数据的大小
	return len(r.svgData)
}

// Clone implements cache.Resource
// 深度复制 SVG 数据以避免并发修改问题
func (r *SVGResource) Clone() cache.Resource {
	if r.svgData == nil || len(r.svgData) == 0 {
		return &SVGResource{}
	}

	// 深度复制数据
	dataCopy := make([]byte, len(r.svgData))
	copy(dataCopy, r.svgData)

	return &SVGResource{
		svgData: dataCopy,
	}
}

// SVGRenderer implements the Renderer interface for SVG rendering
type SVGRenderer struct {
	resourceManager *cache.ResourceManager
}

// NewSVGRenderer creates a new SVG renderer
func NewSVGRenderer(resourceManager *cache.ResourceManager) *SVGRenderer {
	return &SVGRenderer{
		resourceManager: resourceManager,
	}
}

// Render renders an SVG to an image
func (r *SVGRenderer) Render(options RenderOption) (image.Image, error) {
	// Cast options to SVGRenderOption
	svgOptions, ok := options.(SVGRenderOption)
	if !ok {
		return nil, fmt.Errorf("options is not SVGRenderOption")
	}

	// Validate options
	if err := svgOptions.ValidateOption(); err != nil {
		return nil, err
	}

	// Get dimensions
	width, height := svgOptions.GetDimensions()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width=%d, height=%d", width, height)
	}

	// Parse and render SVG
	svgData := svgOptions.GetSVGData()
	img, err := r.renderSVG(svgData, width, height)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// RenderMultiple renders multiple SVGs in parallel
func (r *SVGRenderer) RenderMultiple(options []SVGRenderOption) ([]image.Image, error) {
	results := make([]image.Image, len(options))
	errors := make([]error, len(options))

	var wg sync.WaitGroup
	for i, opt := range options {
		wg.Add(1)
		go func(idx int, opt SVGRenderOption) {
			defer wg.Done()
			img, err := r.Render(opt)
			results[idx] = img
			errors[idx] = err
		}(i, opt)
	}

	wg.Wait()

	// Check for errors
	for i, err := range errors {
		if err != nil {
			return nil, fmt.Errorf("error rendering SVG %d: %w", i, err)
		}
	}

	return results, nil
}

// renderSVG renders an SVG to an RGBA image
// 每次渲染时都重新解析 SVG 数据以避免并发问题
func (r *SVGRenderer) renderSVG(svgData []byte, width, height int) (*image.RGBA, error) {
	// Generate key for cache
	key := r.resourceManager.GenerateKeyFromData(svgData)

	// Try to get from cache
	var svgResource *SVGResource
	item, found := r.resourceManager.GetResource("svg", key, r.resourceManager.GetSVGCache())
	if found {
		if res, ok := item.(*SVGResource); ok {
			svgResource = res
		}
	}

	// 确保我们有SVG数据
	if svgResource == nil {
		// 缓存SVG数据
		svgResource = &SVGResource{
			svgData: svgData,
		}
		r.resourceManager.PutResource("svg", key, r.resourceManager.GetSVGCache(), svgResource)
	}

	// 每次都从数据创建新的SvgIcon，避免并发修改问题
	svgIcon, err := oksvg.ReadIconStream(bytes.NewReader(svgResource.svgData))
	if err != nil {
		return nil, fmt.Errorf("error parsing SVG: %w", err)
	}

	// Set dimensions
	svgIcon.SetTarget(0, 0, float64(width), float64(height))

	// Create output image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Use high-quality rendering
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)

	// Draw SVG
	svgIcon.Draw(raster, 1.0)

	return img, nil
}
