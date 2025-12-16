#!/bin/bash

echo "Go聊天室Web界面启动脚本"
echo "=========================="

# 获取当前脚本所在目录
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# 检查TCP服务器是否在运行
if lsof -i :8000 | grep LISTEN > /dev/null; then
    echo "✓ TCP聊天服务器已在端口8000运行"
else
    echo "⚠ TCP聊天服务器未运行"
    read -p "是否启动TCP聊天服务器? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "启动TCP聊天服务器..."
        cd "$SCRIPT_DIR/server"
        go run server.go &
        SERVER_PID=$!
        cd "$SCRIPT_DIR"
        echo "TCP服务器启动，PID: $SERVER_PID"

        # 等待服务器启动
        sleep 2
    fi
fi

# 检查Web API服务器是否在运行
if lsof -i :8080 | grep LISTEN > /dev/null; then
    echo "⚠ Web API服务器已在端口8080运行"
    read -p "是否关闭现有服务并重启? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        lsof -ti:8080 | xargs kill -9 2>/dev/null
        echo "关闭现有Web API服务器"
    else
        echo "使用现有的Web API服务器"
        echo "打开浏览器访问 http://localhost:8080"
        open http://localhost:8080
        exit 0
    fi
fi

echo "启动Web API服务器..."
cd "$SCRIPT_DIR/api"
go run main.go &
API_PID=$!
echo "Web API服务器启动，PID: $API_PID"

# 等待API服务器启动
sleep 2

echo "打开浏览器访问 http://localhost:8080"
cd "$SCRIPT_DIR"
open http://localhost:8080

# 等待进程
wait $API_PID