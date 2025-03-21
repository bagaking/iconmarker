package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/filter"
	"github.com/bagaking/iconmarker/renderer"
)

// 简化版的SVG渲染选项
type SVGRenderOpts struct {
	SVGData []byte
	Width   int
	Height  int
}

// 实现renderer.SVGRenderOption接口
func (o SVGRenderOpts) GetSVGData() []byte {
	return o.SVGData
}

func (o SVGRenderOpts) GetDimensions() (int, int) {
	return o.Width, o.Height
}

func (o SVGRenderOpts) ValidateOption() error {
	if len(o.SVGData) == 0 {
		return fmt.Errorf("SVG数据为空")
	}
	if o.Width <= 0 || o.Height <= 0 {
		return fmt.Errorf("无效的尺寸: %dx%d", o.Width, o.Height)
	}
	return nil
}

func main() {
	// 打开背景图像
	bgFile := filepath.Join("..", "assets", "background.jpg")
	bgImg, err := openImage(bgFile)
	if err != nil {
		fmt.Printf("无法打开背景图像: %v\n", err)
		return
	}

	// 加载SVG图标
	svgFile := filepath.Join("..", "assets", "icon.svg")
	svgData, err := os.ReadFile(svgFile)
	if err != nil {
		fmt.Printf("无法读取SVG文件: %v\n", err)
		return
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 创建IconMarker实例（用于获取渲染器和资源管理器）
	marker := iconmarker.NewIconMarker()

	// 示例1：渲染原始尺寸的SVG
	renderOriginalSVG(marker, svgData, bgImg, outputDir)

	// 示例2：渲染调整大小的SVG
	renderResizedSVG(marker, svgData, bgImg, outputDir)

	// 示例3：渲染带滤镜的SVG
	renderFilteredSVG(marker, svgData, bgImg, outputDir)
}

// 渲染原始SVG示例
func renderOriginalSVG(marker *core.IconMarker, svgData []byte, bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建SVG渲染选项
	svgOpts := SVGRenderOpts{
		SVGData: svgData,
		Width:   100,
		Height:  100,
	}

	// 创建SVG渲染器
	svgRenderer := renderer.NewSVGRenderer(marker.GetResourceManager())

	// 渲染SVG
	svgImg, err := svgRenderer.Render(svgOpts)
	if err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 将SVG绘制到背景图像上（居中位置）
	centerX := (bounds.Dx() - svgOpts.Width) / 2
	centerY := (bounds.Dy() - svgOpts.Height) / 2
	draw.Draw(img, image.Rect(centerX, centerY, centerX+svgOpts.Width, centerY+svgOpts.Height),
		svgImg, image.Point{}, draw.Over)

	// 保存结果
	outFile := filepath.Join(outputDir, "svg_original.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 渲染调整大小的SVG示例
func renderResizedSVG(marker *core.IconMarker, svgData []byte, bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建SVG渲染选项 - 更大尺寸
	svgOpts := SVGRenderOpts{
		SVGData: svgData,
		Width:   200,
		Height:  200,
	}

	// 创建SVG渲染器
	svgRenderer := renderer.NewSVGRenderer(marker.GetResourceManager())

	// 渲染SVG
	svgImg, err := svgRenderer.Render(svgOpts)
	if err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 将SVG绘制到背景图像上（居中位置）
	centerX := (bounds.Dx() - svgOpts.Width) / 2
	centerY := (bounds.Dy() - svgOpts.Height) / 2
	draw.Draw(img, image.Rect(centerX, centerY, centerX+svgOpts.Width, centerY+svgOpts.Height),
		svgImg, image.Point{}, draw.Over)

	// 保存结果
	outFile := filepath.Join(outputDir, "svg_resized.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 渲染带滤镜的SVG示例
func renderFilteredSVG(marker *core.IconMarker, svgData []byte, bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建SVG渲染选项
	svgOpts := SVGRenderOpts{
		SVGData: svgData,
		Width:   150,
		Height:  150,
	}

	// 创建SVG渲染器
	svgRenderer := renderer.NewSVGRenderer(marker.GetResourceManager())

	// 渲染SVG
	svgImg, err := svgRenderer.Render(svgOpts)
	if err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 应用滤镜 - 添加蓝色色调
	filteredImg, err := marker.ApplyFilter(svgImg, "tint", filter.TintOption{
		Color:     [3]uint8{0, 0, 255}, // 蓝色
		Intensity: 0.7,
	})
	if err != nil {
		fmt.Printf("应用滤镜失败: %v\n", err)
		return
	}

	// 将滤镜处理后的SVG合成到背景图上
	centerX := (bounds.Dx() - svgOpts.Width) / 2
	centerY := (bounds.Dy() - svgOpts.Height) / 2
	draw.Draw(img, image.Rect(centerX, centerY, centerX+svgOpts.Width, centerY+svgOpts.Height),
		filteredImg, image.Point{}, draw.Over)

	// 保存结果
	outFile := filepath.Join(outputDir, "svg_filtered.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 打开图像文件
func openImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// 保存图像为JPEG
func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
}
