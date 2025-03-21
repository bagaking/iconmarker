// Package iconmarker - supports attaching text to existing images
package iconmarker

import (
	"image"
	"io"

	"github.com/bagaking/iconmarker/core"
	"github.com/bagaking/iconmarker/filter"
	"github.com/golang/freetype/truetype"
)

// 向下兼容 - 原有类型别名
type (
	DrawTextOption = core.DrawTextOption
	FontEffect     = core.FontEffect
)

// 原有常量别名
const (
	EShadow  = core.EShadow
	EOutline = core.EOutline
)

// 全局默认icon marker
var defaultIconMarker = core.NewIconMarker()

// CreateImg 创建带有文本的图像（兼容旧API）
func CreateImg(fontBytes, backgroundBytes []byte, drawFontOpt ...DrawTextOption) (*image.RGBA, error) {
	return core.CreateImg(fontBytes, backgroundBytes, drawFontOpt...)
}

// SaveImage2File save image to file (兼容旧API)
func SaveImage2File(img image.Image, path string, encoder func(io.Writer, image.Image) error) error {
	return core.SaveImage2File(img, path, encoder)
}

// DrawCenteredFont draws text on image with center alignment (兼容旧API)
func DrawCenteredFont(f *truetype.Font, outI *image.RGBA, opt DrawTextOption) error {
	return core.DrawCenteredFont(f, outI, opt)
}

// 公开一些核心工具函数
var (
	Bytes2Base64  = core.Bytes2Base64
	Base642Bytes  = core.Base642Bytes
	SaveValToFile = core.SaveValToFile
	PersistFile   = core.PersistFile
)

// PersistStr 字符串模板
var PersistStr = core.PersistStr

// NewIconMarker 创建新的图标标记器
func NewIconMarker() *core.IconMarker {
	return core.NewIconMarker()
}

// ApplyFilter 应用单个滤镜到图像（使用默认IconMarker）
func ApplyFilter(img image.Image, filterName string, option filter.FilterOption) (image.Image, error) {
	return defaultIconMarker.ApplyFilter(img, filterName, option)
}

// ApplyFilters 应用多个滤镜到图像（使用默认IconMarker）
func ApplyFilters(img image.Image, filterNames []string, options []filter.FilterOption) (image.Image, error) {
	return defaultIconMarker.ApplyFilters(img, filterNames, options)
}

// CreateImgWithFilters 创建带有文本和滤镜的图像（使用默认IconMarker）
func CreateImgWithFilters(fontBytes, backgroundBytes []byte,
	filters []string, filterOptions []filter.FilterOption,
	drawFontOpt ...DrawTextOption) (*image.RGBA, error) {
	return defaultIconMarker.CreateImgWithFilters(fontBytes, backgroundBytes, filters, filterOptions, drawFontOpt...)
}

// GetFilterManager 获取默认IconMarker的滤镜管理器
func GetFilterManager() *filter.FilterManager {
	return defaultIconMarker.GetFilterManager()
}
