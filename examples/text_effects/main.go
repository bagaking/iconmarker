package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/core"
	"github.com/golang/freetype/truetype"
)

func main() {
	// 打开背景图像
	bgFile := filepath.Join("..", "assets", "background.jpg")
	bgImg, err := openImage(bgFile)
	if err != nil {
		fmt.Printf("无法打开背景图像: %v\n", err)
		return
	}

	// 加载字体文件（可选，失败时会使用默认字体）
	var font *truetype.Font
	fontFile := filepath.Join("..", "assets", "font.ttf")
	fontData, err := ioutil.ReadFile(fontFile)
	if err == nil {
		font, err = truetype.Parse(fontData)
		if err != nil {
			fmt.Printf("警告: 无法解析字体文件: %v\n", err)
			// 继续执行，将使用默认字体
		}
	} else {
		fmt.Printf("警告: 无法读取字体文件: %v\n", err)
		fmt.Println("将使用默认嵌入字体")
		// 继续执行，将使用默认字体
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 示例1：带有阴影效果的文本
	shadowExample(bgImg, font, outputDir)

	// 示例2：带有轮廓效果的文本
	outlineExample(bgImg, font, outputDir)

	// 示例3：组合多种效果
	combinedEffectsExample(bgImg, font, outputDir)
}

// 阴影效果示例
func shadowExample(bgImg image.Image, font *truetype.Font, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建带有阴影效果的文本选项
	opt := iconmarker.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "阴影效果文本",
		XOffset:   0, // 中心对齐
		YOffset:   0, // 中心对齐
	}.SetStaticSize(48).AddShadow(color.RGBA{R: 0, G: 0, B: 0, A: 192}, 4)

	// 添加文本 - 如果font为nil，底层会使用默认字体
	err := iconmarker.DrawCenteredFont(font, img, opt)
	if err != nil {
		fmt.Printf("添加文本失败: %v\n", err)
		// 如果渲染真的失败，仍然使用模拟渲染作为最后的兜底
		textImg := simulateTextRender(img, "阴影效果文本 (模拟)", 48)
		draw.Draw(img, bounds, textImg, image.Point{}, draw.Over)
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "text_shadow.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 轮廓效果示例
func outlineExample(bgImg image.Image, font *truetype.Font, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建带有轮廓效果的文本选项
	opt := iconmarker.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 0, A: 255}, // 黄色文本
		Text:      "轮廓效果文本",
		XOffset:   0, // 中心对齐
		YOffset:   0, // 中心对齐
	}.SetStaticSize(48).AddOutline(color.RGBA{R: 0, G: 0, B: 128, A: 255}, 3) // 深蓝色轮廓

	// 添加文本 - 如果font为nil，底层会使用默认字体
	err := iconmarker.DrawCenteredFont(font, img, opt)
	if err != nil {
		fmt.Printf("添加文本失败: %v\n", err)
		// 如果渲染真的失败，仍然使用模拟渲染作为最后的兜底
		textImg := simulateTextRender(img, "轮廓效果文本 (模拟)", 48)
		draw.Draw(img, bounds, textImg, image.Point{}, draw.Over)
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "text_outline.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 组合效果示例
func combinedEffectsExample(bgImg image.Image, font *truetype.Font, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 创建带有阴影和轮廓效果的文本选项
	opt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // 白色文本
		Text:      "组合效果文本",
		XOffset:   0, // 中心对齐
		YOffset:   0, // 中心对齐
	}.SetStaticSize(48).
		AddOutline(color.RGBA{R: 0, G: 0, B: 128, A: 255}, 3). // 深蓝色轮廓
		AddShadow(color.RGBA{R: 0, G: 0, B: 0, A: 192}, 5)     // 黑色阴影

	// 添加文本 - 如果font为nil，底层会使用默认字体
	err := core.DrawCenteredFont(font, img, opt)
	if err != nil {
		fmt.Printf("添加文本失败: %v\n", err)
		// 如果渲染真的失败，仍然使用模拟渲染作为最后的兜底
		textImg := simulateTextRender(img, "组合效果文本 (模拟)", 48)
		draw.Draw(img, bounds, textImg, image.Point{}, draw.Over)
	}

	// 保存结果
	outFile := filepath.Join(outputDir, "text_combined_effects.jpg")
	if err := saveImage(img, outFile); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", outFile)
}

// 模拟文本渲染 - 当字体加载失败时的备选方案
func simulateTextRender(baseImg image.Image, text string, fontSize int) image.Image {
	bounds := baseImg.Bounds()
	textImg := image.NewRGBA(bounds)

	// 获取图像中心位置
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	// 创建一个简单的文本指示器
	rectWidth := len(text) * fontSize / 2
	rectHeight := fontSize * 2

	// 绘制文本背景
	for y := centerY - rectHeight/2; y < centerY+rectHeight/2; y++ {
		for x := centerX - rectWidth/2; x < centerX+rectWidth/2; x++ {
			// 确保坐标在图像范围内
			if x >= 0 && x < bounds.Dx() && y >= 0 && y < bounds.Dy() {
				// 半透明黑色背景
				textImg.Set(x, y, color.RGBA{0, 0, 0, 128})
			}
		}
	}

	// 在背景上绘制简单的"T"字母，表示这是文本位置
	tWidth := fontSize / 2
	tHeight := fontSize

	// 绘制"T"的横线
	for y := centerY - tHeight/2; y < centerY-tHeight/2+tHeight/5; y++ {
		for x := centerX - tWidth; x < centerX+tWidth; x++ {
			if x >= 0 && x < bounds.Dx() && y >= 0 && y < bounds.Dy() {
				textImg.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	// 绘制"T"的竖线
	for y := centerY - tHeight/2; y < centerY+tHeight/2; y++ {
		for x := centerX - tWidth/5; x < centerX+tWidth/5; x++ {
			if x >= 0 && x < bounds.Dx() && y >= 0 && y < bounds.Dy() {
				textImg.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	return textImg
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
