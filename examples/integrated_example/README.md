# 整合示例

这个示例展示了如何使用重组后的 Icon Marker 功能，包括：

1. 使用兼容的旧式 API 创建带有文本的图像
2. 使用新的 API 和滤镜功能创建带有文本和特效的图像
3. 先创建图像，然后单独应用滤镜
4. 对现有图像应用不透明度滤镜

## 运行方法

1. 首先，请确保放置了测试资源文件：
   - 将一个 TTF 字体文件放在 `examples/assets/font.ttf`
   - 将一个 JPG 背景图片放在 `examples/assets/background.jpg`

2. 执行以下命令运行示例：
   ```bash
   cd examples/integrated_example
   go run main.go
   ```

3. 输出图像将保存在 `examples/integrated_example/output` 目录中。

## 不同示例的说明

1. **old_api.png**: 使用兼容的旧式 API 创建的图像
2. **new_api_with_filters.png**: 使用新 API 和滤镜（灰度 + 蓝色色调）创建的图像
3. **inverted.png**: 使用反转滤镜处理后的图像
4. **transparent.png**: 应用了透明度滤镜的背景图像

## API 设计说明

新的 API 设计采用了更加模块化的方法：

- **Marker 实例**: 可以创建自定义的 `IconMarker` 实例，或者使用全局默认实例
- **滤镜系统**: 提供了统一的滤镜接口和多种内置滤镜
- **缓存系统**: 内部使用 LRU 缓存以提高性能
- **渲染器**: 使用专用渲染器处理不同类型的内容

这种设计保持向下兼容的同时，提供了更强大和灵活的功能。 