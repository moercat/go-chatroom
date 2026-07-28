# Go Chatroom

> 🎯 **从零实现 TCP IM** — 纯 Go 标准库，协议文档和源码一一对应

> 一个 TCP + Web 双协议多人聊天室。**纯 `net` / `encoding/json` / `sync` 实现**，零第三方依赖。涵盖单聊、群聊、公屏、上下线通知、在线列表、聊天日志等功能。

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Stars](https://img.shields.io/badge/Stars-8-brightgreen)
![Status](https://img.shields.io/badge/Status-完整可跑-blue)
[![CI](https://github.com/moercat/go-chatroom/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/go-chatroom/actions/workflows/go.yml)
[![Protocol](https://img.shields.io/badge/协议-PROTOCOL.md-green)](./PROTOCOL.md)

## ✨ 为什么是 go-chatroom

- 🎓 **教学价值高** — TCP 长连接 + 消息路由 + 广播 + 并发控制一次学完
- 📖 **协议文档完整** — [PROTOCOL.md](./PROTOCOL.md) 与源码一一对应
- 🔌 **双协议** — TCP（学习用）+ Web（实用）
- 📦 **零依赖** — 仅 Go 标准库

## 🚀 一键启动

```bash
make web     # 等价于 ./start_web_chatroom.sh
```

或手动开 3 个终端：

```bash
make tcp     # TCP 8888
make api     # HTTP 8080
open web/index.html
```

## 📚 协议（PROTOCOL.md）

JSON over TCP，9 种操作、3 种区域：

| Op | 操作 | Op | 操作 |
| --- | --- | --- | --- |
| 1 | 公屏聊天 | 6 | 群聊 |
| 2 | 登出 | 7 | 创建群组 |
| 3 | 登录 | 8 | 列出群组 |
| 4 | 修改资料 | 9 | 列出在线用户 |
| 5 | 私聊 | | |

完整字段 / 示例 / 端口 / 错误处理 → [PROTOCOL.md](./PROTOCOL.md)

## ✅ 功能矩阵

| 功能 | TCP | Web |
| --- | --- | --- |
| 公屏 / 私聊 / 群聊 | ✅ | ✅ |
| 上下线通知 | ✅ | ✅ |
| 在线列表 / 群组列表 | ✅ | ✅ |
| 个人资料修改 | ✅ | ✅ |
| 聊天日志（文件） | ✅ | — |

## 🏗️ 架构

```
TCP Client ──┐
             ├─→ TCP Server (8888) ─→ chat.log
Web Client ──┘
             └─→ HTTP API (8080) ─→ ConnMap + GroupMap
```

## 🎯 学习要点

1. **TCP 长连接** — `net.Listener` + 一个连接一个 goroutine
2. **消息路由** — `Read` 函数按 `Op` 字段分发
3. **广播机制** — `for client := range ConnMap { client.Send(msg) }`
4. **并发安全** — `sync.RWMutex` 保护共享 map
5. **粘包处理** — 用 `\n` 做消息分隔

## 📁 目录

```
go-chatroom/
├── server/         ← TCP 聊天服务器（核心）
├── client/         ← 命令行客户端
├── api/            ← HTTP API（Web 桥接）
├── web/            ← Web 界面
├── pkg/
│   ├── chat/       ← 消息路由
│   ├── enum/       ← 消息类型
│   ├── entity/     ← 数据模型
│   └── log/        ← 聊天日志
├── PROTOCOL.md     ← 协议文档
├── start_web_chatroom.sh  ← 一键启动
└── README.md
