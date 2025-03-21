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
	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/renderer"
)

// 自定义SVG渲染选项，实现renderer.SVGRenderOption接口
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

	// 创建安全徽章SVG - 直接在代码中定义SVG内容
	securityBadgeSVG := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="512" height="512" viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
  <!-- 定义渐变 -->
  <defs>
    <linearGradient id="circleGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4285F4" />
      <stop offset="50%" stop-color="#5B9BFF" />
      <stop offset="100%" stop-color="#4285F4" />
    </linearGradient>
    <linearGradient id="textGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#DB4437" />
      <stop offset="50%" stop-color="#FF5C51" />
      <stop offset="100%" stop-color="#DB4437" />
    </linearGradient>
  </defs>
  
  <!-- 背景透明 -->
  <!-- 外圈 - 蓝色圆形带渐变 -->
  <circle cx="256" cy="256" r="240" stroke="url(#circleGradient)" stroke-width="24" fill="none"/>
  
  <!-- 标题背景 - 浅蓝色矩形（调整位置和尺寸） -->
  <rect x="25" y="140" width="462" height="92" rx="12" fill="#E8F0FE"/>
</svg>`

	// 创建画布
	width, height := 512, 512
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白色背景
	draw.Draw(background, background.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// 渲染SVG底图
	badgeOpts := &SVGRenderOpts{
		SVGData:  []byte(securityBadgeSVG),
		Width:    width,
		Height:   height,
		Position: image.Point{X: 0, Y: 0},
	}

	badgeImg, err := svgRenderer.Render(badgeOpts)
	if err != nil {
		fmt.Printf("渲染徽章底图失败: %v\n", err)
		return
	}

	// 将SVG底图绘制到背景上
	draw.Draw(background, background.Bounds(), badgeImg, image.Point{}, draw.Over)

	// 添加"Security..."文本
	securityTextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 66, G: 133, B: 244, A: 255}, // 谷歌蓝
		Text:      "badge ...",
		XOffset:   0,
		YOffset:   -60, // 向上调整，因为底图向上移动了
	}.SetStaticSize(60)

	if err := core.DrawCenteredFont(nil, background, securityTextOpt); err != nil {
		fmt.Printf("渲染 badge 文本失败: %v\n", err)
		return
	}

	// 添加"E0"文本 - 红色大字
	e0TextOpt := core.DrawTextOption{
		FontColor: color.RGBA{R: 55, G: 68, B: 219, A: 255},
		Text:      "OK",
		XOffset:   0,
		YOffset:   120,
	}.SetStaticSize(168)

	if err := core.DrawCenteredFont(nil, background, e0TextOpt); err != nil {
		fmt.Printf("渲染 OK 文本失败: %v\n", err)
		return
	}

	// 在OK下方添加一个蓝色圆形图标，内含火车头
	// 直接定义带颜色的SVG图标，而不是后期应用滤镜
	trainIconSVG := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="100" height="100" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
  <!-- 外圈虚线 -->
  <circle cx="50" cy="50" r="48" stroke="#1A3A47" stroke-width="2" stroke-dasharray="8 4" fill="none"/>
  
  <!-- 主体圆圈 -->
  <circle cx="50" cy="50" r="40" fill="#2A98B7" stroke="#1A3A47" stroke-width="2"/>
  
  <!-- 横穿线 -->
  <line x1="10" y1="50" x2="90" y2="50" stroke="#7CCDE0" stroke-width="6"/>
  
  <!-- 火车头图案 -->
  <!-- 车身主体 -->
  <path d="M25 55 L65 55 L65 37 L57 37 L57 30 L45 30 L45 37 L25 37 Z" fill="white" stroke="white" stroke-width="2"/>
  
  <!-- 前部车头 -->
  <path d="M65 40 L75 45 L75 50 L65 50 Z" fill="white" stroke="white" stroke-width="2"/>
  
  <!-- 烟囱 -->
  <rect x="50" y="25" width="6" height="12" rx="2" fill="white" stroke="white"/>
  
  <!-- 车轮 -->
  <circle cx="35" cy="62" r="5" fill="white" stroke="#1A3A47" stroke-width="1"/>
  <circle cx="55" cy="62" r="5" fill="white" stroke="#1A3A47" stroke-width="1"/>
</svg>`

	// 图标大小和位置
	iconSize := 100
	iconPosition := image.Point{X: (width - iconSize) / 2, Y: 400} // OK文本下方居中

	// 渲染SVG图标
	trainIconOpts := &SVGRenderOpts{
		SVGData:  []byte(trainIconSVG),
		Width:    iconSize,
		Height:   iconSize,
		Position: image.Point{X: 0, Y: 0}, // 临时位置
	}

	trainIconImg, err := svgRenderer.Render(trainIconOpts)
	if err != nil {
		fmt.Printf("渲染图标失败: %v\n", err)
		return
	}

	// 将图标绘制到背景上正确的位置
	draw.Draw(background, image.Rect(
		iconPosition.X,
		iconPosition.Y,
		iconPosition.X+iconSize,
		iconPosition.Y+iconSize),
		trainIconImg, image.Point{}, draw.Over)

	// 保存结果
	saveAsPNG(background, filepath.Join(outputDir, "badge.png"))
	fmt.Println("徽章图像已保存到", filepath.Join(outputDir, "badge.png"))
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
