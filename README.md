# Go 聊天室

> 基于 Go TCP 的多人在线聊天室

![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue)
![Status](https://img.shields.io/badge/Status-学习参考-blue)
![Stars](https://img.shields.io/badge/Stars-8-brightgreen)
[![CI](https://github.com/moercat/go-chatroom/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/go-chatroom/actions/workflows/go.yml)

## 定位

**解决了什么问题**：从零实现一个基于 Go TCP 的多人在线聊天室，包含全部常用功能（单聊、群聊、公屏、上下线通知）。

**目标用户**：Go 网络编程学习者。

**学习价值**：纯 Go 标准库实现（net 包），展示 TCP server/client 架构设计。

## 快速开始

```bash
git clone https://github.com/moercat/go-chatroom.git
cd go-chatroom

# 启动 TCP 服务端
cd server && go run server.go

# 启动 Web API（新标签页）
cd api && go run main.go

# 客户端（开多个窗口）
cd client && go run client.go
```

环境要求：Go 1.18+。

## 架构

```
go-chatroom/
├── server/     ← TCP 聊天服务器（核心）
├── client/     ← 命令行聊天客户端
├── api/        ← Web API 服务（WebSocket 桥接）
├── web/        ← 前端界面（HTML/CSS/JS）
└── start_web_chatroom.sh  ← 一键启动脚本
```

## 功能

- [x] 加入/退出聊天室
- [x] 公屏聊天（广播）
- [x] 单聊（私信）
- [x] 群聊（创建/加入群组）
- [x] 上下线通知
- [x] 在线用户列表
- [x] 聊天日志
- [x] 个人资料修改
- [x] Web 界面（新功能）

## 技术栈

| 组件 | 选型 |
|------|------|
| 语言 | Go 1.18 |
| 网络 | net 标准库（TCP） |
| 前端 | 原生 HTML/CSS/JS |

## 状态

| 维度 | 说明 |
|------|------|
| 当前阶段 | 功能完整，可作为学习参考 |
| 最后更新 | 2025-12-16 |
| 活跃度 | 💤 休眠 |
