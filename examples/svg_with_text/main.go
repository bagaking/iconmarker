package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/assets"
	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/filter"
	"github.com/bagaking/iconmarker/renderer"
)

// 自定义SVG渲染选项，实现renderer.SVGRenderOption接口
type SVGRenderOpts struct {
	SVGData  []byte
	Width    int
	Height   int
	Position image.Point
}

// GetSVGData 返回SVG数据
func (o *SVGRenderOpts) GetSVGData() []byte {
	return o.SVGData
}

// GetDimensions 返回SVG尺寸
func (o *SVGRenderOpts) GetDimensions() (width, height int) {
	return o.Width, o.Height
}

// ValidateOption 验证选项
func (o *SVGRenderOpts) ValidateOption() error {
	if len(o.SVGData) == 0 {
		return fmt.Errorf("SVG data is empty")
	}
	if o.Width <= 0 || o.Height <= 0 {
		return fmt.Errorf("invalid dimensions: width=%d, height=%d", o.Width, o.Height)
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

	// 读取外部SVG图标
	svgFile := filepath.Join("..", "assets", "icon.svg")
	svgData, err := os.ReadFile(svgFile)
	if err != nil {
		fmt.Printf("无法读取外部SVG图标: %v\n", err)
		return
	}

	// 获取内嵌的SVG图标
	diamondSvgData, err := assets.IconDiamondMarker.Load()
	if err != nil {
		fmt.Printf("无法获取内嵌的菱形标记图标: %v\n", err)
		return
	}

	locationPinData, err := assets.IconLocationPin.Load()
	if err != nil {
		fmt.Printf("无法获取内嵌的位置标记图标: %v\n", err)
		return
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 创建IconMarker实例
	marker := iconmarker.NewIconMarker()

	// 示例1：SVG图标在左，文本在右（使用外部SVG）
	example1_IconLeft(marker, bgImg, svgData, outputDir)

	// 示例2：SVG图标在上，文本在下（使用内嵌菱形图标）
	example2_IconTop(marker, bgImg, diamondSvgData, outputDir)

	// 示例3：文本环绕SVG图标（使用内嵌位置标记图标）
	example3_TextAround(marker, bgImg, locationPinData, outputDir)

	// 示例4：带滤镜的组合效果（使用外部SVG）
	example4_WithFilters(marker, bgImg, svgData, outputDir)

	// 示例5：内嵌SVG图标示例
	example5_EmbeddedIcons(marker, bgImg, outputDir)

	fmt.Println("所有示例已完成，输出图像保存在", outputDir, "目录")
}

// 在图像上渲染SVG的辅助函数
func renderSVGOnImage(svgRenderer *renderer.SVGRenderer, img draw.Image, opts *SVGRenderOpts) error {
	// 首先使用Render方法获取SVG图像
	svgImg, err := svgRenderer.Render(opts)
	if err != nil {
		return err
	}

	// 然后将SVG图像绘制到目标图像上的指定位置
	draw.Draw(img, image.Rect(
		opts.Position.X,
		opts.Position.Y,
		opts.Position.X+opts.Width,
		opts.Position.Y+opts.Height),
		svgImg, image.Point{}, draw.Over)

	return nil
}

// SVG图标在左，文本在右的布局
func example1_IconLeft(marker *core.IconMarker, bgImg image.Image, svgData []byte, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 计算图标和文本的位置
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	centerY := imgHeight / 2

	// 创建SVG渲染选项 - 放在图像左侧
	svgSize := 150 // SVG图标大小
	svgOpts := &SVGRenderOpts{
		SVGData:  svgData,
		Width:    svgSize,
		Height:   svgSize,
		Position: image.Point{X: imgWidth/4 - svgSize/2, Y: centerY - svgSize/2}, // 放在左侧四分之一处
	}

	// 渲染SVG到图像
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)
	if err := renderSVGOnImage(svgRenderer, img, svgOpts); err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 创建文本选项 - 放在图像右侧
	textOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "左侧SVG图标 + 右侧文本",
		XOffset:   imgWidth / 4, // 向右侧偏移，放在右侧四分之三处
		YOffset:   0,            // 垂直居中
	}.SetStaticSize(48).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	// 渲染文本到图像
	if err := core.DrawCenteredFont(nil, img, textOpt); err != nil {
		fmt.Printf("渲染文本失败: %v\n", err)
		return
	}

	// 保存结果
	saveAsJPEG(img, filepath.Join(outputDir, "svg_left_text_right.jpg"))
}

// SVG图标在上，文本在下的布局
func example2_IconTop(marker *core.IconMarker, bgImg image.Image, svgData []byte, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 计算图标和文本的位置
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	centerX := imgWidth / 2

	// 创建SVG渲染选项 - 放在图像上部
	svgSize := 180 // SVG图标大小
	svgOpts := &SVGRenderOpts{
		SVGData:  svgData,
		Width:    svgSize,
		Height:   svgSize,
		Position: image.Point{X: centerX - svgSize/2, Y: imgHeight/4 - svgSize/2}, // 放在上方四分之一处
	}

	// 渲染SVG到图像
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)
	if err := renderSVGOnImage(svgRenderer, img, svgOpts); err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 创建文本选项 - 放在图像下部
	textOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "上部SVG图标 + 下部文本",
		XOffset:   0,             // 水平居中
		YOffset:   imgHeight / 4, // 向下偏移，放在下方四分之三处
	}.SetStaticSize(48).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	// 渲染文本到图像
	if err := core.DrawCenteredFont(nil, img, textOpt); err != nil {
		fmt.Printf("渲染文本失败: %v\n", err)
		return
	}

	// 保存结果
	saveAsJPEG(img, filepath.Join(outputDir, "svg_top_text_bottom.jpg"))
}

// 文本环绕SVG图标的布局
func example3_TextAround(marker *core.IconMarker, bgImg image.Image, svgData []byte, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 计算中心位置
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	centerX := imgWidth / 2
	centerY := imgHeight / 2

	// 创建SVG渲染选项 - 放在图像中心位置
	svgSize := 200 // SVG图标大小
	svgOpts := &SVGRenderOpts{
		SVGData:  svgData,
		Width:    svgSize,
		Height:   svgSize,
		Position: image.Point{X: centerX - svgSize/2, Y: centerY - svgSize/2}, // 居中放置
	}

	// 渲染SVG到图像
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)
	if err := renderSVGOnImage(svgRenderer, img, svgOpts); err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 添加上方文本
	topTextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "上方文本",
		XOffset:   0,
		YOffset:   -svgSize/2 - 40, // 放在SVG上方
	}.SetStaticSize(42).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	if err := core.DrawCenteredFont(nil, img, topTextOpt); err != nil {
		fmt.Printf("渲染上方文本失败: %v\n", err)
		return
	}

	// 添加下方文本
	bottomTextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "下方文本",
		XOffset:   0,
		YOffset:   svgSize/2 + 40, // 放在SVG下方
	}.SetStaticSize(42).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	if err := core.DrawCenteredFont(nil, img, bottomTextOpt); err != nil {
		fmt.Printf("渲染下方文本失败: %v\n", err)
		return
	}

	// 添加左侧文本
	leftTextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "左侧",
		XOffset:   -svgSize/2 - 60, // 放在SVG左侧
		YOffset:   0,
	}.SetStaticSize(36).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	if err := core.DrawCenteredFont(nil, img, leftTextOpt); err != nil {
		fmt.Printf("渲染左侧文本失败: %v\n", err)
		return
	}

	// 添加右侧文本
	rightTextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "右侧",
		XOffset:   svgSize/2 + 60, // 放在SVG右侧
		YOffset:   0,
	}.SetStaticSize(36).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	if err := core.DrawCenteredFont(nil, img, rightTextOpt); err != nil {
		fmt.Printf("渲染右侧文本失败: %v\n", err)
		return
	}

	// 保存结果
	saveAsJPEG(img, filepath.Join(outputDir, "text_around_svg.jpg"))
}

// 带滤镜的组合效果
func example4_WithFilters(marker *core.IconMarker, bgImg image.Image, svgData []byte, outputDir string) {
	// 复制背景图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 首先应用滤镜到背景
	filteredImg, err := marker.ApplyFilters(img,
		[]string{"grayscale", "tint"},
		[]filter.FilterOption{
			filter.GrayscaleOption{PreserveAlpha: true},
			filter.TintOption{
				Color:     [3]uint8{50, 50, 180}, // 添加蓝色色调
				Intensity: 0.3,                   // 色调强度
			},
		})
	if err != nil {
		fmt.Printf("应用滤镜失败: %v\n", err)
		return
	}

	// 将滤镜结果转为RGBA以便进一步处理
	resultImg := image.NewRGBA(bounds)
	draw.Draw(resultImg, bounds, filteredImg, image.Point{}, draw.Src)

	// 计算中心位置
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 创建SVG渲染选项 - 放在图像右上角
	svgSize := 150
	svgOpts := &SVGRenderOpts{
		SVGData:  svgData,
		Width:    svgSize,
		Height:   svgSize,
		Position: image.Point{X: imgWidth - svgSize - 50, Y: 50}, // 右上角
	}

	// 渲染SVG到图像
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)
	if err := renderSVGOnImage(svgRenderer, resultImg, svgOpts); err != nil {
		fmt.Printf("渲染SVG失败: %v\n", err)
		return
	}

	// 添加标题文本
	titleOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 215, B: 0, A: 255}, // 金色
		Text:      "滤镜效果示例",
		XOffset:   0,
		YOffset:   -imgHeight / 4, // 放在上方
	}.SetStaticSize(64).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 4).AddShadow(color.RGBA{R: 0, G: 0, B: 0, A: 180}, 5)

	if err := core.DrawCenteredFont(nil, resultImg, titleOpt); err != nil {
		fmt.Printf("渲染标题文本失败: %v\n", err)
		return
	}

	// 添加描述文本
	descOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 220, G: 220, B: 255, A: 255}, // 淡蓝色
		Text:      "背景使用灰度+蓝色色调滤镜",
		XOffset:   0,
		YOffset:   imgHeight / 4, // 放在下方
	}.SetStaticSize(36).AddOutline(color.RGBA{R: 30, G: 30, B: 80, A: 255}, 2)

	if err := core.DrawCenteredFont(nil, resultImg, descOpt); err != nil {
		fmt.Printf("渲染描述文本失败: %v\n", err)
		return
	}

	// 添加说明文本在SVG图标旁边
	noteOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "SVG图标",
		XOffset:   -(imgWidth/2 - (imgWidth - svgSize - 50) - svgSize/2),
		YOffset:   -imgHeight/2 + 50 + svgSize/2,
	}.SetStaticSize(28).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 2)

	if err := core.DrawCenteredFont(nil, resultImg, noteOpt); err != nil {
		fmt.Printf("渲染说明文本失败: %v\n", err)
		return
	}

	// 保存结果
	saveAsPNG(resultImg, filepath.Join(outputDir, "svg_text_with_filters.png"))
}

// 示例5：内嵌SVG图标示例
func example5_EmbeddedIcons(marker *core.IconMarker, bgImg image.Image, outputDir string) {
	// 获取所有可用的内嵌图标
	iconNames, err := assets.ListAvailableIcons()
	if err != nil {
		fmt.Printf("获取图标列表失败: %v\n", err)
		return
	}

	// 创建图像
	bounds := bgImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, bgImg, image.Point{}, draw.Src)

	// 添加标题
	titleOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:      "内嵌SVG图标示例",
		XOffset:   0,
		YOffset:   -bounds.Dy() / 3, // 放在上方
	}.SetStaticSize(48).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 3)

	if err := core.DrawCenteredFont(nil, img, titleOpt); err != nil {
		fmt.Printf("渲染标题失败: %v\n", err)
		return
	}

	// 实例化SVG渲染器
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)

	// 水平排列图标
	iconSize := 120
	spacing := 40
	totalWidth := len(iconNames) * (iconSize + spacing)
	startX := (bounds.Dx() - totalWidth) / 2

	// 渲染每个图标
	for i, name := range iconNames {
		// 获取图标数据
		iconData, err := assets.GetSVGIcon(name)
		if err != nil {
			fmt.Printf("获取图标 %s 失败: %v\n", name, err)
			continue
		}

		// 计算位置
		x := startX + i*(iconSize+spacing)
		y := bounds.Dy()/2 - iconSize/2

		// 设置渲染选项
		svgOpts := &SVGRenderOpts{
			SVGData:  iconData,
			Width:    iconSize,
			Height:   iconSize,
			Position: image.Point{X: x, Y: y},
		}

		// 渲染SVG到图像
		if err := renderSVGOnImage(svgRenderer, img, svgOpts); err != nil {
			fmt.Printf("渲染图标 %s 失败: %v\n", name, err)
			continue
		}

		// 添加图标名称
		labelOpt := core.DrawTextOption{
			FontColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			Text:      name,
			XOffset:   x + iconSize/2 - bounds.Dx()/2, // 相对于中心点的x偏移
			YOffset:   y + iconSize + 30,              // 图标下方30像素
		}.SetStaticSize(20).AddOutline(color.RGBA{R: 0, G: 0, B: 0, A: 255}, 2)

		if err := core.DrawCenteredFont(nil, img, labelOpt); err != nil {
			fmt.Printf("渲染图标名称失败: %v\n", err)
			continue
		}
	}

	// 保存结果
	saveAsPNG(img, filepath.Join(outputDir, "embedded_icons.png"))
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

// 保存为JPEG图像
func saveAsJPEG(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建输出文件失败: %v\n", err)
		return
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
		fmt.Printf("编码JPEG失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", filename)
}

// 保存为PNG图像
func saveAsPNG(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建输出文件失败: %v\n", err)
		return
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Printf("编码PNG失败: %v\n", err)
		return
	}

	fmt.Printf("已保存: %s\n", filename)
}
