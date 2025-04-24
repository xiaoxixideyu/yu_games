#!/bin/bash

# 检查 Docker 是否已安装
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "错误: 请先安装 Docker 和 Docker Compose"
    exit 1
fi

# 构建和启动容器
echo "开始构建和启动容器..."
docker-compose up -d --build

# 检查容器状态
echo "检查容器状态..."
docker-compose ps

echo "应用已启动!"
echo "访问地址: http://localhost"
echo "使用以下命令查看日志:"
echo "  前端和反向代理: docker-compose logs -f nginx"
echo "  后端API: docker-compose logs -f backend"
echo "  数据库: docker-compose logs -f db"
echo ""
echo "使用以下命令停止应用:"
echo "  docker-compose down" 