# IconMarker 示例集

这个目录包含了 IconMarker 库的各种使用示例，展示了库的主要功能和用法。

## 快速开始

直接运行设置脚本来准备环境和资源文件：

```bash
./setup_and_run.sh
```

这个脚本会：
1. 创建必要的测试资源文件（背景图片、SVG图标）和下载字体文件
2. 为每个示例创建输出目录
3. 设置正确的执行权限
4. 询问是否要立即运行所有示例

## 示例列表

### 1. 基础文本示例 (basic_text)

演示如何在图片上添加基本文本，包括：
- 使用传统API添加文本
- 使用新API添加文本

```bash
cd basic_text
go run main.go
```

### 2. 文本效果示例 (text_effects)

展示各种文本效果，包括：
- 添加带阴影的文本
- 添加带轮廓的文本
- 组合使用阴影和轮廓效果

```bash
cd text_effects
go run main.go
```

### 3. SVG渲染示例 (svg_rendering)

展示如何渲染SVG图标到图片上，包括：
- 渲染原始大小的SVG
- 渲染调整大小的SVG
- 应用滤镜后渲染SVG

```bash
cd svg_rendering
go run main.go
```

### 4. 组合滤镜示例 (combined_filters)

展示如何使用和组合多种滤镜，包括：
- 使用内置组合滤镜
- 顺序应用多个滤镜
- 使用自定义组合应用滤镜

```bash
cd combined_filters
go run main.go
```

### 5. 整合示例 (integrated_example)

展示如何综合使用IconMarker的各种功能，包括：
- 传统API与新API的对比
- 使用多种滤镜处理图像
- 使用高级渲染功能

```bash
cd integrated_example
go run main.go
```

### 6. SVG与文本组合示例 (svg_with_text)

演示如何在同一个图像上组合SVG图标与文本，包括：
- SVG图标在左，文本在右的布局
- SVG图标在上，文本在下的布局（使用内嵌菱形图标）
- 文本环绕SVG图标的布局（使用内嵌位置标记图标）
- 应用滤镜后组合SVG与文本
- 展示所有内嵌SVG图标

```bash
cd svg_with_text
go run main.go
```

## 内嵌SVG图标

IconMarker 提供 22 个内嵌 SVG 图标，无需读取外部文件。以
`assets.AllIcons()` 或 `assets.ListAvailableIcons()` 获取完整清单。

使用示例：
```go
import "github.com/bagaking/iconmarker/assets"

// 使用类型安全的枚举加载图标
svgData, err := assets.IconDiamondMarker.Load()

// 也可按名称读取
locationData, err := assets.GetSVGIcon("location-pin")

// 获取所有可用图标的名称列表
iconNames, err := assets.ListAvailableIcons()

// 通过名称获取任意图标
svgData, err := assets.GetSVGIcon("diamond-marker")
```

## 运行所有示例

使用提供的脚本一次性运行所有示例：

```bash
./run_examples.sh
```

所有示例的输出将保存在各自目录下的 `output` 文件夹中。

## 资源文件

示例使用的资源文件位于 `assets` 目录中：
- `background.jpg`: 用作背景的图片
- `icon.svg`: 用于SVG渲染示例的图标
- `font.ttf`: 用于文本渲染的字体文件

**关于字体文件**：示例可读取 `font.ttf`，库本身也内嵌了阿里巴巴普惠体作为默认字体；调用文本 API 时传入空字体字节即可使用默认字体。

## 自定义示例

如果要创建自己的示例，建议遵循以下模式：
1. 创建一个新的目录或单文件
2. 加载背景图像或创建空白图像
3. 创建IconMarker实例
4. 使用适当的API添加内容或应用滤镜
5. 保存结果到output目录

## 注意事项

- 所有示例都需要 Go 1.23 或更高版本
- 确保在运行示例前已通过`go mod tidy`安装所有依赖
- 如果您在运行示例时遇到问题，请查看IconMarker的文档或提交issue
