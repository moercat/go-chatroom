# Go Chatroom

> 基于 Go TCP + Web 的多人在线聊天室

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Stars](https://img.shields.io/badge/Stars-8-brightgreen)
![Status](https://img.shields.io/badge/Status-学习参考-blue)
[![CI](https://github.com/moercat/go-chatroom/actions/workflows/go.yml/badge.svg)](https://github.com/moercat/go-chatroom/actions/workflows/go.yml)

## 这是什么

一个完整的 TCP + Web 双协议聊天室，**纯 Go 标准库实现**（仅 `net`、`encoding/json`、`sync`）。涵盖单聊、群聊、公屏、上下线通知、在线列表、聊天日志、个人资料修改等功能。适合：

- Go 网络编程学习
- TCP 长连接服务器设计参考
- Web 前后端实时通信示例

## 架构

```
┌─────────────┐         ┌─────────────┐
│  TCP Client │◄───────►│             │
│  client/    │         │  TCP Server │◄──── chat.log
└─────────────┘         │  server/    │       (聊天日志)
                        │  :8888      │
┌─────────────┐         │             │
│  Web Client │◄───────►│             │
│  web/       │  WSS    │  HTTP API   │
└─────────────┘  HTTP   │  api/       │
                        │  :8080      │
                        └─────────────┘
                             │
                             ▼
                        ┌─────────────┐
                        │ ConnMap     │  在线连接
                        │ GroupMap    │  群组成员
                        └─────────────┘
```

## 启动

### 方式一：一键启动（推荐）

```bash
./start_web_chatroom.sh
```

会自动启动 TCP server (8888) + HTTP API (8080) + 打开 Web 界面。

### 方式二：手动启动

```bash
# 终端 1：启动 TCP 服务
cd server && go run server.go

# 终端 2：启动 HTTP API
cd api && go run main.go

# 终端 3：打开 Web 界面
open web/index.html
# 或浏览器访问 file:// 路径
```

### 方式三：TCP 命令行客户端

```bash
# 多开几个终端测试
cd client && go run client.go
```

## 功能矩阵

| 功能 | TCP | Web | 状态 |
| --- | --- | --- | --- |
| 加入 / 退出 | ✅ | ✅ | |
| 公屏聊天 | ✅ | ✅ | |
| 私聊 | ✅ | ✅ | |
| 群聊（创建/发送） | ✅ | ✅ | |
| 上下线通知 | ✅ | ✅ | |
| 在线用户列表 | ✅ | ✅ | |
| 群组列表 | ✅ | ✅ | |
| 个人资料修改 | ✅ | ✅ | |
| 聊天日志（服务端文件） | ✅ | — | server/chat.log |
| 消息持久化 | ❌ | ❌ | 重启清空 |
| 离线消息 | ❌ | ❌ | |
| 端到端加密 | ❌ | ❌ | |

## 协议

详见 [PROTOCOL.md](./PROTOCOL.md)。

简要：JSON over TCP，每条消息一行 `\n` 分隔。

## 目录

```
go-chatroom/
├── server/         ← TCP 聊天服务器 (核心)
├── client/         ← 命令行聊天客户端
├── api/            ← HTTP API (WebSocket 桥接)
├── web/            ← Web 界面 (HTML/JS)
├── pkg/
│   ├── chat/       ← 消息路由
│   ├── enum/       ← 消息类型定义
│   ├── entity/     ← 数据模型
│   └── log/        ← 聊天日志
├── start_web_chatroom.sh  ← 一键启动脚本
├── PROTOCOL.md     ← 协议说明
└── README.md
```

## 技术栈

| 组件 | 选型 |
| --- | --- |
| 语言 | Go 1.21 |
| 网络 | net 标准库 (TCP) |
| 序列化 | encoding/json |
| 并发 | sync (Mutex + RWMutex) |
| 前端 | 原生 HTML/CSS/JS (无框架) |

## 学习要点

1. **TCP 长连接管理** — `net.Listener` 接受连接，每个连接一个 goroutine
2. **消息路由** — `Read` 函数按 `Op` 字段分发到 `ntyLogin` / `Quit` / `SendPrivateMessage` / ...
3. **广播机制** — `for client := range ConnMap { client.Send(msg) }`
4. **并发安全** — `sync.RWMutex` 保护 `ConnMap` 和 `GroupMap`
5. **JSON over TCP** — 用换行符做消息分隔（替代粘包处理）

## 已知限制

- 无消息持久化（重启清空）
- 群组成员在内存（重启丢失）
- 无心跳保活（实际需要客户端每 30s 发一次消息）
- 私聊无加密
- 无身份验证（任何人可起任意用户名）

## 后续方向

- [ ] WebSocket 替换 HTTP polling
- [ ] 消息持久化（SQLite/BoltDB）
- [ ] 群组成员持久化
- [ ] TLS 加密传输
- [ ] 单元测试（协议层）

## 许可证

MIT
