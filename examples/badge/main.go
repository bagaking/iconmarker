package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagaking/iconmarker"
	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/filter/utils"
	"github.com/bagaking/iconmarker/renderer"
)

// BadgeStyle 定义徽章样式
type BadgeStyle struct {
	Name        string
	CircleColor string // 外圈渐变色
	TextColor   string // 文本颜色
	StatusText  string // 状态文本
	IconColor   string // 图标颜色
	LargeText   string // 大字内容
}

// 预定义徽章样式
var (
	BadgeStyleUrgent = BadgeStyle{
		Name:        "urgent",
		CircleColor: "#DB4437", // 红色
		TextColor:   "#DB4437",
		StatusText:  "Biz No.001",
		IconColor:   "#FF5C51",
		LargeText:   "URGENT",
	}
	BadgeStyleResolved = BadgeStyle{
		Name:        "resolved",
		StatusText:  "Biz No.001",
		LargeText:   "SOLVED",
		CircleColor: "#34A853", // 绿色
		TextColor:   "#34A853",
		IconColor:   "#4CAF50",
	}
	BadgeStyleInProgress = BadgeStyle{
		Name:        "in_progress",
		StatusText:  "Biz No.001",
		LargeText:   "INPROC",
		CircleColor: "#F48542", // 橙色
		TextColor:   "#F48542",
		IconColor:   "#FF9B5B",
	}
	BadgeStyleProblem = BadgeStyle{
		Name:        "problem",
		StatusText:  "Biz No.001",
		LargeText:   "跟进中",
		CircleColor: "#4285F4", // 蓝色
		TextColor:   "#4285F4",
		IconColor:   "#5B9BFF",
	}
)

// 徽章SVG模板
const badgeTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="512" height="512" viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
  <!-- 定义渐变 -->
  <defs>
    <linearGradient id="circleGradient" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="{{.CircleColor}}" />
      <stop offset="50%" stop-color="{{.IconColor}}" />
      <stop offset="100%" stop-color="{{.CircleColor}}" />
    </linearGradient>
  </defs>
  
  <!-- 背景透明 -->
  <!-- 外圈 - 圆形带渐变 -->
  <circle cx="256" cy="256" r="240" stroke="url(#circleGradient)" stroke-width="24" fill="none"/>
  
  <!-- 标题背景 - 浅色矩形 -->
  <rect x="25" y="140" width="462" height="92" rx="12" fill="#E8F0FE"/>
</svg>`

// 图标SVG模板
const iconTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="100" height="100" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
  <!-- 外圈虚线 -->
  <circle cx="50" cy="50" r="48" stroke="{{.IconColor}}" stroke-width="2" stroke-dasharray="8 4" fill="none"/>
  
  <!-- 主体圆圈 -->
  <circle cx="50" cy="50" r="40" fill="{{.IconColor}}" stroke="{{.IconColor}}" stroke-width="2"/>
  
  <!-- 横穿线 -->
  <line x1="10" y1="50" x2="90" y2="50" stroke="white" stroke-width="6"/>
  
  <!-- 火车头图案 -->
  <!-- 车身主体 -->
  <path d="M25 55 L65 55 L65 37 L57 37 L57 30 L45 30 L45 37 L25 37 Z" fill="white" stroke="white" stroke-width="2"/>
  
  <!-- 前部车头 -->
  <path d="M65 40 L75 45 L75 50 L65 50 Z" fill="white" stroke="white" stroke-width="2"/>
  
  <!-- 烟囱 -->
  <rect x="50" y="25" width="6" height="12" rx="2" fill="white" stroke="white"/>
  
  <!-- 车轮 -->
  <circle cx="35" cy="62" r="5" fill="white" stroke="{{.IconColor}}" stroke-width="1"/>
  <circle cx="55" cy="62" r="5" fill="white" stroke="{{.IconColor}}" stroke-width="1"/>
</svg>`

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

// 生成徽章
func generateBadge(style BadgeStyle) error {
	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 创建IconMarker实例
	marker := iconmarker.NewIconMarker()
	resourceMgr := marker.GetResourceManager()
	svgRenderer := renderer.NewSVGRenderer(resourceMgr)

	// 替换模板中的颜色变量
	badgeSvg := strings.ReplaceAll(badgeTemplate, "{{.CircleColor}}", style.CircleColor)
	badgeSvg = strings.ReplaceAll(badgeSvg, "{{.IconColor}}", style.IconColor)

	// 创建画布
	width, height := 512, 512
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	bgColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白色背景
	draw.Draw(background, background.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// 渲染SVG底图
	badgeOpts := &SVGRenderOpts{
		SVGData:  []byte(badgeSvg),
		Width:    width,
		Height:   height,
		Position: image.Point{X: 0, Y: 0},
	}

	badgeImg, err := svgRenderer.Render(badgeOpts)
	if err != nil {
		return fmt.Errorf("渲染徽章底图失败: %v", err)
	}

	// 将SVG底图绘制到背景上
	draw.Draw(background, background.Bounds(), badgeImg, image.Point{}, draw.Over)

	// 解析基础颜色
	baseColor, err := utils.ParseHexColor(style.TextColor)
	if err != nil {
		return fmt.Errorf("解析文本颜色失败: %v", err)
	}

	// 为状态文本创建浅色版本 (与白色混合)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	statusTextColor := utils.LerpColor(baseColor, white, 0.3) // 30% 偏向白色，让状态文本更柔和

	// 为大字文本创建深色版本 (与黑色混合)
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	largeTextColor := utils.LerpColor(baseColor, black, 0.1) // 10% 偏向黑色，让大字文本更醒目

	// 添加状态文本 (使用浅色)
	statusTextOpt := core.DrawTextOption{
		FontColor: statusTextColor,
		Text:      style.StatusText,
		XOffset:   0,
		YOffset:   -60,
	}.SetStaticSize(60)

	if err := core.DrawCenteredFont(nil, background, statusTextOpt); err != nil {
		return fmt.Errorf("渲染状态文本失败: %v", err)
	}

	// 添加大字文本 (使用深色)
	largeTextOpt := core.DrawTextOption{
		FontColor: largeTextColor,
		Text:      style.LargeText,
		XOffset:   0,
		YOffset:   90,
	}.SetStaticSize(120)

	if err := core.DrawCenteredFont(nil, background, largeTextOpt); err != nil {
		return fmt.Errorf("渲染大字文本失败: %v", err)
	}

	// 替换图标模板中的颜色变量
	iconSvg := strings.ReplaceAll(iconTemplate, "{{.IconColor}}", style.IconColor)

	// 图标大小和位置
	iconSize := 100
	iconPosition := image.Point{X: (width - iconSize) / 2, Y: 400}

	// 渲染SVG图标
	iconOpts := &SVGRenderOpts{
		SVGData:  []byte(iconSvg),
		Width:    iconSize,
		Height:   iconSize,
		Position: image.Point{X: 0, Y: 0},
	}

	iconImg, err := svgRenderer.Render(iconOpts)
	if err != nil {
		return fmt.Errorf("渲染图标失败: %v", err)
	}

	// 将图标绘制到背景上
	draw.Draw(background, image.Rect(
		iconPosition.X,
		iconPosition.Y,
		iconPosition.X+iconSize,
		iconPosition.Y+iconSize),
		iconImg, image.Point{}, draw.Over)

	// 保存结果
	outputPath := filepath.Join(outputDir, fmt.Sprintf("badge_%s.png", style.Name))
	if err := saveAsPNG(background, outputPath); err != nil {
		return fmt.Errorf("保存徽章失败: %v", err)
	}

	fmt.Printf("徽章已保存到 %s\n", outputPath)
	return nil
}

// 保存为PNG图像
func saveAsPNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("编码PNG失败: %v", err)
	}
	return nil
}

func main() {
	// 生成所有样式的徽章
	styles := []BadgeStyle{
		BadgeStyleUrgent,
		BadgeStyleResolved,
		BadgeStyleInProgress,
		BadgeStyleProblem,
	}

	for _, style := range styles {
		if err := generateBadge(style); err != nil {
			fmt.Printf("生成%s样式徽章失败: %v\n", style.Name, err)
		}
	}
}
