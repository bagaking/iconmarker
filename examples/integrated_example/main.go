package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/filter"
	"github.com/golang/freetype/truetype"
)

func main() {
	// 加载背景图像
	bgFile := filepath.Join("..", "assets", "background.jpg")
	bgImg, err := openImage(bgFile)
	if err != nil {
		fmt.Printf("无法打开背景图像: %v\n", err)
		return
	}

	// 将背景图像转换为字节数组
	bgBuf := new(bytes.Buffer)
	err = jpeg.Encode(bgBuf, bgImg, &jpeg.Options{Quality: 90})
	if err != nil {
		fmt.Printf("无法编码背景图像: %v\n", err)
		return
	}
	bgBytes := bgBuf.Bytes()

	// 尝试加载字体文件（可选，失败时将使用默认字体）
	var fontBytes []byte
	fontFile := filepath.Join("..", "assets", "font.ttf")
	var loadedFontSuccess bool = false

	fontBytes, err = ioutil.ReadFile(fontFile)
	if err == nil {
		// 验证字体文件是否有效
		_, err = truetype.Parse(fontBytes)
		if err == nil {
			loadedFontSuccess = true
		} else {
			fmt.Printf("字体文件解析失败: %v\n", err)
			fmt.Println("将使用默认嵌入字体")
			// 创建一个空字节数组，让API使用默认字体
			fontBytes = []byte{}
		}
	} else {
		fmt.Printf("无法读取字体文件: %v\n", err)
		fmt.Println("将使用默认嵌入字体")
		// 创建一个空字节数组，让API使用默认字体
		fontBytes = []byte{}
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	if loadedFontSuccess {
		fmt.Println("使用成功加载的外部字体文件")
	} else {
		fmt.Println("使用默认嵌入字体")
	}

	// 方法1：使用旧的兼容 API
	img1, err := iconmarker.CreateImg(
		fontBytes,
		bgBytes,
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:      "Hello World",
		}.SetAdaptedSize(600, 300).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 4),
	)
	if err != nil {
		fmt.Printf("Error creating image with old API: %v\n", err)
		return
	}

	// 保存使用旧 API 创建的图像
	if err := iconmarker.SaveImage2File(img1, filepath.Join(outputDir, "old_api.png"), png.Encode); err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	// 方法2：使用新的 API 与滤镜
	// 创建一个自定义 IconMarker 实例
	marker := iconmarker.NewIconMarker()

	// 使用滤镜创建图像
	img2, err := marker.CreateImgWithFilters(
		fontBytes,
		bgBytes,
		[]string{"grayscale", "tint"},
		[]filter.FilterOption{
			filter.GrayscaleOption{PreserveAlpha: true},
			filter.TintOption{
				Color:     [3]uint8{0, 0, 255}, // 蓝色色调
				Intensity: 0.5,
			},
		},
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:      "Text with Filters",
		}.SetAdaptedSize(600, 300).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 4),
	)
	if err != nil {
		fmt.Printf("Error creating image with filters: %v\n", err)
		return
	}

	// 保存使用新 API 创建的图像
	if err := marker.SaveImage2File(img2, filepath.Join(outputDir, "new_api_with_filters.png"), png.Encode); err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	// 方法3：使用全局 API 创建图像，然后应用滤镜
	img3, err := iconmarker.CreateImg(
		fontBytes,
		bgBytes,
		iconmarker.DrawTextOption{
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:      "Apply Filters Later",
		}.SetAdaptedSize(600, 300).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 4),
	)
	if err != nil {
		fmt.Printf("Error creating image: %v\n", err)
		return
	}

	// 应用反转滤镜
	invertedImg, err := iconmarker.ApplyFilter(img3, "invert", filter.InvertOption{InvertAlpha: false})
	if err != nil {
		fmt.Printf("Error applying invert filter: %v\n", err)
		return
	}

	// 保存处理后的图像
	outputFile3, err := os.Create(filepath.Join(outputDir, "inverted.png"))
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile3.Close()
	if err := png.Encode(outputFile3, invertedImg); err != nil {
		fmt.Printf("Error encoding PNG: %v\n", err)
		return
	}

	// 方法4：加载 JPEG 图像，应用不透明度滤镜，然后保存
	// 应用不透明度滤镜
	transparentImg, err := iconmarker.ApplyFilter(bgImg, "opacity", filter.OpacityOption{Opacity: 0.5})
	if err != nil {
		fmt.Printf("Error applying opacity filter: %v\n", err)
		return
	}

	// 保存处理后的图像
	outputFile4, err := os.Create(filepath.Join(outputDir, "transparent.png"))
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile4.Close()
	if err := png.Encode(outputFile4, transparentImg); err != nil {
		fmt.Printf("Error encoding PNG: %v\n", err)
		return
	}

	fmt.Println("All examples completed successfully!")
	fmt.Println("Output images are saved in the", outputDir, "directory.")
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
