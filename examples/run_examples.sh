#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 确保资源文件存在
if [ ! -d "assets" ]; then
  mkdir -p assets
  echo -e "${RED}警告: 资源目录不存在，已创建。${NC}"
  echo -e "${RED}请确保在assets目录中放置以下文件:${NC}"
  echo -e "${RED}- background.jpg (背景图片)${NC}"
  echo -e "${RED}- font.ttf (字体文件)${NC}"
  echo -e "${RED}- icon.svg (SVG图标)${NC}"
  exit 1
fi

# 提示
echo -e "${GREEN}运行所有 IconMarker 示例${NC}"
echo "=================================================="

# 运行单文件示例
run_single_example() {
  local example_file=$1
  local example_name=$(basename "$example_file" .go)
  
  echo -e "${GREEN}运行示例: $example_name${NC}"
  
  # 创建输出目录
  mkdir -p "output/$example_name"
  
  # 运行示例
  go run "$example_file"
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 成功运行: $example_name${NC}"
  else
    echo -e "${RED}✗ 运行失败: $example_name${NC}"
  fi
  echo "--------------------------------------------------"
}

# 运行目录示例
run_dir_example() {
  local example_dir=$1
  local example_name=$(basename "$example_dir")
  
  echo -e "${GREEN}运行示例: $example_name${NC}"
  
  # 进入示例目录
  cd "$example_dir"
  
  # 创建输出目录
  mkdir -p "output"
  
  # 运行示例
  go run main.go
  
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 成功运行: $example_name${NC}"
  else
    echo -e "${RED}✗ 运行失败: $example_name${NC}"
  fi
  
  # 返回原目录
  cd "$SCRIPT_DIR"
  echo "--------------------------------------------------"
}

# 运行根目录下的单文件示例
echo -e "${GREEN}运行单文件示例:${NC}"
for example in *.go; do
  if [ -f "$example" ] && [[ "$example" != *"_test.go" ]]; then
    run_single_example "$example"
  fi
done

# 运行子目录示例
echo -e "${GREEN}运行目录示例:${NC}"
for dir in */; do
  if [ -d "$dir" ] && [ "$dir" != "assets/" ] && [ "$dir" != "output/" ]; then
    if [ -f "${dir}main.go" ]; then
      run_dir_example "$dir"
    fi
  fi
done

echo -e "${GREEN}所有示例运行完成${NC}"
echo "可以在各示例目录的output目录中查看输出文件" 