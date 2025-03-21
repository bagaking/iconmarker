package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/filter"
)

func main() {
	// 打开背景图像
	bgFile := filepath.Join("..", "assets", "background.jpg")
	bgImg, err := openImage(bgFile)
	if err != nil {
		fmt.Printf("无法打开背景图像: %v\n", err)
		return
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 示例1：使用内置组合滤镜
	compositeSingleFilter(bgImg, outputDir)

	// 示例2：顺序应用多个滤镜
	sequentialFilters(bgImg, outputDir)

	// 示例3：自定义组合滤镜
	customCompositeFilter(bgImg, outputDir)
}

// 使用内置组合滤镜示例
func compositeSingleFilter(bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建IconMarker实例
	marker := iconmarker.NewIconMarker()

	// 创建组合滤镜选项
	compositeOpt := filter.CompositeOption{
		Filters: []filter.Filter{
			filter.NewGrayscaleFilter(),
			filter.NewInvertFilter(),
		},
		Options: []filter.FilterOption{
			filter.GrayscaleOption{PreserveAlpha: true},
			filter.InvertOption{InvertAlpha: false},
		},
		StopOnError: true,
	}

	// 应用组合滤镜
	filteredImg, err := marker.ApplyFilter(img, "composite", compositeOpt)
	if err != nil {
		fmt.Printf("应用组合滤镜失败: %v\n", err)
		return
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "composite_filter.jpg")
	if err := saveImage(filteredImg, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 顺序应用多个滤镜示例
func sequentialFilters(bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建IconMarker实例
	marker := iconmarker.NewIconMarker()

	// 应用滤镜1: 灰度
	grayImg, err := marker.ApplyFilter(img, "grayscale", filter.GrayscaleOption{
		PreserveAlpha: true,
	})
	if err != nil {
		fmt.Printf("应用灰度滤镜失败: %v\n", err)
		return
	}

	// 应用滤镜2: 蓝色色调
	tintedImg, err := marker.ApplyFilter(grayImg, "tint", filter.TintOption{
		Color:     [3]uint8{0, 0, 255}, // 蓝色
		Intensity: 0.5,                 // 中等强度
	})
	if err != nil {
		fmt.Printf("应用色调滤镜失败: %v\n", err)
		return
	}

	// 应用滤镜3: 降低不透明度
	finalImg, err := marker.ApplyFilter(tintedImg, "opacity", filter.OpacityOption{
		Opacity: 0.8, // 80%不透明度
	})
	if err != nil {
		fmt.Printf("应用不透明度滤镜失败: %v\n", err)
		return
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "sequential_filters.jpg")
	if err := saveImage(finalImg, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 自定义组合滤镜示例 - 使用ApplyFilters
func customCompositeFilter(bgImg image.Image, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建滤镜和选项数组
	filterNames := []string{
		"grayscale",
		"tint",
		"invert",
	}

	filterOptions := []filter.FilterOption{
		filter.GrayscaleOption{PreserveAlpha: true},
		filter.TintOption{
			Color:     [3]uint8{255, 0, 0}, // 红色
			Intensity: 0.3,                 // 低强度
		},
		filter.InvertOption{InvertAlpha: false},
	}

	// 使用全局API应用多个滤镜
	filteredImg, err := iconmarker.ApplyFilters(img, filterNames, filterOptions)
	if err != nil {
		fmt.Printf("应用多个滤镜失败: %v\n", err)
		return
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "custom_composite.jpg")
	if err := saveImage(filteredImg, outFile); err != nil {
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
