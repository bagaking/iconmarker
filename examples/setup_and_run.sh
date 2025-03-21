#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# 函数：显示分隔线
print_separator() {
  echo -e "${YELLOW}=====================================${NC}"
}

# 脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 显示欢迎消息
print_separator
echo -e "${GREEN}欢迎使用 IconMarker 示例设置脚本${NC}"
print_separator

# 步骤2: 确保每个示例目录都有输出目录
echo -e "${GREEN}步骤 1: 创建输出目录${NC}"
mkdir -p output

# 为每个子目录创建output目录
for dir in */; do
  if [ -d "$dir" ] && [ "$dir" != "assets/" ] && [ "$dir" != "output/" ]; then
    mkdir -p "${dir}output"
    echo "为 ${dir} 创建输出目录"
  fi
done

# 步骤3: 给运行脚本添加执行权限
echo -e "${GREEN}步骤 2: 添加执行权限${NC}"
chmod +x run_examples.sh
echo "已添加 run_examples.sh 的执行权限"

# 步骤4: 运行示例
print_separator
echo -e "${GREEN}准备工作已完成！${NC}"
echo -e "现在您可以运行示例了，有以下选项："
echo -e "1. 全部运行: ${YELLOW}./run_examples.sh${NC}"
echo -e "2. 单独运行某个示例:"
echo -e "   ${YELLOW}cd basic_text && go run main.go${NC}"
echo -e "   ${YELLOW}cd text_effects && go run main.go${NC}"
echo -e "   ${YELLOW}cd svg_rendering && go run main.go${NC}"
echo -e "   ${YELLOW}cd combined_filters && go run main.go${NC}"
echo -e "   ${YELLOW}cd integrated_example && go run main.go${NC}"
echo -e "   ${YELLOW}cd svg_with_text && go run main.go${NC}"
echo -e "   ${YELLOW}go run filter_example.go${NC} (如果有)"

# 询问用户是否要立即运行所有示例
echo
read -p "是否立即运行所有示例? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  print_separator
  echo -e "${GREEN}运行所有示例...${NC}"
  ./run_examples.sh
else
  echo -e "${GREEN}设置完成。您可以稍后手动运行示例。${NC}"
fi

print_separator 