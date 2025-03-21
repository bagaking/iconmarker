// Package assets 提供项目所需的内嵌资源文件
package assets

import (
	"embed"
)

// DefaultFontFS 包含默认字体文件的文件系统
//
//go:embed MPLUSRounded1c-ExtraBold.ttf
var DefaultFontFS embed.FS

// DefaultFontPath 是默认字体文件的路径
const DefaultFontPath = "MPLUSRounded1c-ExtraBold.ttf"

// GetDefaultFont 返回默认字体的字节数据
func GetDefaultFont() ([]byte, error) {
	return DefaultFontFS.ReadFile(DefaultFontPath)
}
