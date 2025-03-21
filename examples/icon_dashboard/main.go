package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/assets"
	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/filter"
	"github.com/bagaking/iconmarker/renderer"
)

// SVGRenderOption 实现接口
type SVGRenderOpts struct {
	SVGData  []byte
	Width    int
	Height   int
	Position image.Point
}

func (o *SVGRenderOpts) GetSVGData() []byte {
	return o.SVGData
}

func (o *SVGRenderOpts) GetDimensions() (width, height int) {
	return o.Width, o.Height
}

func (o *SVGRenderOpts) ValidateOption() error {
	if len(o.SVGData) == 0 {
		return fmt.Errorf("SVG数据为空")
	}
	if o.Width <= 0 || o.Height <= 0 {
		return fmt.Errorf("无效的尺寸: width=%d, height=%d", o.Width, o.Height)
	}
	return nil
}

func main() {
	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 创建IconMarker实例
	marker := iconmarker.NewIconMarker()
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)

	// 创建背景
	width, height := 800, 400
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{R: 240, G: 245, B: 250, A: 255} // 浅蓝色背景
	draw.Draw(background, background.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// 获取两个不同的SVG图标 - 使用IconType枚举代替字符串
	icon1Type := assets.IconAlert
	icon2Type := assets.IconHeart

	// 使用IconType的Load方法加载图标数据
	icon1Data, err := icon1Type.Load()
	if err != nil {
		fmt.Printf("获取图标1失败: %v\n", err)
		return
	}

	icon2Data, err := icon2Type.Load()
	if err != nil {
		fmt.Printf("获取图标2失败: %v\n", err)
		return
	}

	// 渲染第一个SVG图标 - 保持原始颜色
	icon1Opts := &SVGRenderOpts{
		SVGData:  icon1Data,
		Width:    150,
		Height:   150,
		Position: image.Point{X: 100, Y: 125},
	}

	icon1Img, err := svgRenderer.Render(icon1Opts)
	if err != nil {
		fmt.Printf("渲染图标1失败: %v\n", err)
		return
	}

	// 将SVG图像绘制到目标图像上
	draw.Draw(background, image.Rect(
		icon1Opts.Position.X,
		icon1Opts.Position.Y,
		icon1Opts.Position.X+icon1Opts.Width,
		icon1Opts.Position.Y+icon1Opts.Height),
		icon1Img, image.Point{}, draw.Over)

	// 渲染第二个SVG图标 - 应用颜色滤镜
	icon2Opts := &SVGRenderOpts{
		SVGData:  icon2Data,
		Width:    150,
		Height:   150,
		Position: image.Point{X: 550, Y: 125},
	}

	icon2Img, err := svgRenderer.Render(icon2Opts)
	if err != nil {
		fmt.Printf("渲染图标2失败: %v\n", err)
		return
	}

	// 使用修复后的filter滤镜为图标应用红色
	filterManager := filter.NewFilterManager()
	redColor := [3]uint8{220, 50, 50}
	filteredIcon2, err := filterManager.QuickTint(icon2Img, redColor, 0.7)
	if err != nil {
		fmt.Printf("应用滤镜失败: %v\n", err)
		return
	}

	// 将滤镜处理后的SVG图像绘制到目标图像上
	draw.Draw(background, image.Rect(
		icon2Opts.Position.X,
		icon2Opts.Position.Y,
		icon2Opts.Position.X+icon2Opts.Width,
		icon2Opts.Position.Y+icon2Opts.Height),
		filteredIcon2, image.Point{}, draw.Over)

	// 添加中央文本
	textOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 30, G: 41, B: 59, A: 255},
		Text:      "图标和文本组合示例",
		XOffset:   0,
		YOffset:   -150,
	}.SetStaticSize(32).AddOutline(color.RGBA{R: 255, G: 255, B: 255, A: 150}, 1)

	if err := core.DrawCenteredFont(nil, background, textOpt); err != nil {
		fmt.Printf("渲染标题失败: %v\n", err)
		return
	}

	// 添加说明文本 - 显示图标枚举名称
	descOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 71, G: 85, B: 105, A: 255},
		Text:      fmt.Sprintf("左侧为 %s 图标，右侧为应用了红色滤镜的 %s 图标", icon1Type, icon2Type),
		XOffset:   0,
		YOffset:   -100,
	}.SetStaticSize(20)

	if err := core.DrawCenteredFont(nil, background, descOpt); err != nil {
		fmt.Printf("渲染描述失败: %v\n", err)
		return
	}

	// 添加底部说明
	footerOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 100, G: 116, B: 139, A: 255},
		Text:      "通过iconmarker库可以轻松将SVG图标和文字组合在一起",
		XOffset:   0,
		YOffset:   150,
	}.SetStaticSize(18)

	if err := core.DrawCenteredFont(nil, background, footerOpt); err != nil {
		fmt.Printf("渲染底部说明失败: %v\n", err)
		return
	}

	// 保存结果
	saveAsPNG(background, filepath.Join(outputDir, "simple_icons_with_text.png"))
	fmt.Println("简单图标与文本组合已保存到", filepath.Join(outputDir, "simple_icons_with_text.png"))
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
}
