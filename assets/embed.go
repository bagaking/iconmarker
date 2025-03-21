// Package assets 提供项目所需的内嵌资源文件
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// DefaultFontFS 包含默认字体文件的文件系统
//
//go:embed AlibabaPuHuiTi-3-105-Heavy.ttf
var DefaultFontFS embed.FS

// DefaultFontPath 是默认字体文件的路径
const DefaultFontPath = "AlibabaPuHuiTi-3-105-Heavy.ttf"

// IconsFS 包含所有默认SVG图标的文件系统
//
//go:embed icons/*.svg
var IconsFS embed.FS

// IconsDir 是图标目录的路径
const IconsDir = "icons"

// IconType 表示预定义的图标类型
type IconType string

// 预定义图标枚举常量
const (
	IconAlarmClock    IconType = "alarm-clock"
	IconAlert         IconType = "alert"
	IconCircle        IconType = "circle"
	IconCthulhu       IconType = "cthulhu"
	IconCyberpunk     IconType = "cyberpunk"
	IconDiamondMarker IconType = "diamond-marker"
	IconEvaluation    IconType = "evaluation"
	IconFlower        IconType = "flower"
	IconHeart         IconType = "heart"
	IconLightning     IconType = "lightning"
	IconLocationPin   IconType = "location-pin"
	IconNorthStar     IconType = "north-star"
	IconPaperPlane    IconType = "paper-plane"
	IconPerson        IconType = "person"
	IconProhibited    IconType = "prohibited"
	IconRobot         IconType = "robot"
	IconRoundedRect   IconType = "rounded-rect"
	IconSpaceship     IconType = "spaceship"
	IconTarget        IconType = "target"
	IconTeam          IconType = "team"
	IconTodo          IconType = "todo"
	IconWarning       IconType = "warning"
)

// AllIcons 返回所有预定义图标的切片
func AllIcons() []IconType {
	return []IconType{
		IconAlarmClock,
		IconAlert,
		IconCircle,
		IconCthulhu,
		IconCyberpunk,
		IconDiamondMarker,
		IconEvaluation,
		IconFlower,
		IconHeart,
		IconLightning,
		IconLocationPin,
		IconNorthStar,
		IconPaperPlane,
		IconPerson,
		IconProhibited,
		IconRobot,
		IconRoundedRect,
		IconSpaceship,
		IconTarget,
		IconTeam,
		IconTodo,
		IconWarning,
	}
}

// String 实现Stringer接口，返回图标类型的字符串表示
func (i IconType) String() string {
	return string(i)
}

// Load 加载并返回图标的SVG数据
func (i IconType) Load() ([]byte, error) {
	return GetSVGIcon(string(i))
}

// ParseIconType 将字符串解析为IconType
// 如果字符串不匹配任何预定义图标，返回错误
func ParseIconType(name string) (IconType, error) {
	// 移除可能的.svg后缀
	name = strings.TrimSuffix(name, ".svg")

	// 检查是否为有效的图标名称
	for _, icon := range AllIcons() {
		if string(icon) == name {
			return icon, nil
		}
	}

	return "", fmt.Errorf("无效的图标名称: %s", name)
}

// GetDefaultFont 返回默认字体的字节数据
func GetDefaultFont() ([]byte, error) {
	return DefaultFontFS.ReadFile(DefaultFontPath)
}

// GetSVGIcon 返回指定名称的SVG图标的字节数据
// 参数 name 不需要包含 .svg 扩展名和 icons/ 路径前缀
func GetSVGIcon(name string) ([]byte, error) {
	if !strings.HasSuffix(name, ".svg") {
		name = name + ".svg"
	}

	path := filepath.Join(IconsDir, name)
	return IconsFS.ReadFile(path)
}

// ListAvailableIcons 返回所有可用图标的名称列表（不含.svg扩展名）
func ListAvailableIcons() ([]string, error) {
	var icons []string

	err := fs.WalkDir(IconsFS, IconsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".svg") {
			name := filepath.Base(path)
			name = strings.TrimSuffix(name, ".svg")
			icons = append(icons, name)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("列出图标失败: %w", err)
	}

	return icons, nil
}
