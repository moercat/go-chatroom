# 协议说明

go-chatroom 使用 JSON over TCP（HTTP API 用 WebSocket 桥接）。

## 消息格式

每条消息都是 JSON 字符串，以换行符 (`\n`) 分隔。

```json
{
  "name":      "alice",
  "op":        1,
  "msg":       "hello world",
  "target":    "bob",
  "group":     "",
  "timestamp": 1700000000,
  "area":      "public_screen"
}
```

### 字段说明

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | string | 是 | 发送者用户名（登录后由服务端填充） |
| `op` | int | 是 | 操作类型，见下表 |
| `msg` | string | 否 | 消息内容 |
| `target` | string | 私聊时必填 | 私聊目标用户名 |
| `group` | string | 群聊时必填 | 群组名 |
| `timestamp` | int64 | 否 | Unix 时间戳（秒） |
| `area` | string | 否 | 区域：`public_screen` / `group_chat` / `private_chat` |

### 操作类型（Op）

| 值 | 常量 | 说明 |
| --- | --- | --- |
| 1 | `Chat` | 公屏聊天 |
| 2 | `Logout` | 登出 |
| 3 | `Login` | 登录（首次连接） |
| 4 | `UpdateUser` | 修改个人资料 |
| 5 | `PrivateChat` | 私聊 |
| 6 | `GroupChat` | 群聊 |
| 7 | `CreateGroup` | 创建群组 |
| 8 | `ListGroups` | 列出所有群组 |
| 9 | `ListUsers` | 列出在线用户 |

### 区域（Area）

| 值 | 说明 |
| --- | --- |
| `public_screen` | 公屏（默认） |
| `group_chat` | 群聊 |
| `private_chat` | 私聊 |

## 交互示例

### 1. 登录（首次连接后立即发送）

客户端 → 服务端：
```json
{"op": 3, "name": "alice"}
```

服务端 → 所有客户端（广播上线通知）：
```json
{"op": 3, "name": "alice", "msg": "alice 上线了"}
```

### 2. 公屏聊天

客户端 → 服务端：
```json
{"op": 1, "name": "alice", "msg": "大家好", "area": "public_screen"}
```

服务端 → 所有客户端（广播）：
```json
{"op": 1, "name": "alice", "msg": "大家好", "area": "public_screen", "timestamp": 1700000000}
```

### 3. 私聊

发送方 → 服务端：
```json
{"op": 5, "name": "alice", "target": "bob", "msg": "在吗？", "area": "private_chat"}
```

服务端 → 接收方 (bob)：
```json
{"op": 5, "name": "alice", "target": "bob", "msg": "在吗？", "area": "private_chat"}
```

### 4. 创建群组

客户端 → 服务端：
```json
{"op": 7, "name": "alice", "group": "tech-team"}
```

服务端 → 创建者：
```json
{"op": 7, "name": "alice", "group": "tech-team", "msg": "群组创建成功"}
```

### 5. 群聊

客户端 → 服务端：
```json
{"op": 6, "name": "alice", "group": "tech-team", "msg": "开会", "area": "group_chat"}
```

服务端 → 群组所有成员（广播）：
```json
{"op": 6, "name": "alice", "group": "tech-team", "msg": "开会", "area": "group_chat"}
```

### 6. 列出在线用户

客户端 → 服务端：
```json
{"op": 9, "name": "alice"}
```

服务端 → 请求者：
```json
{"op": 9, "msg": "alice, bob, carol"}
```

### 7. 列出群组

客户端 → 服务端：
```json
{"op": 8, "name": "alice"}
```

服务端 → 请求者：
```json
{"op": 8, "msg": "tech-team, friends, family"}
```

### 8. 修改个人资料

客户端 → 服务端：
```json
{"op": 4, "name": "alice", "msg": "新昵称"}
```

服务端 → 广播：
```json
{"op": 4, "name": "alice", "msg": "alice 已更名为 新昵称"}
```

### 9. 登出

客户端 → 服务端：
```json
{"op": 2, "name": "alice"}
```

服务端 → 所有客户端（广播下线通知）：
```json
{"op": 2, "name": "alice", "msg": "alice 下线了"}
```

## 客户端实现示例

参考 `client/client.go`：

```go
conn, _ := net.Dial("tcp", "localhost:8888")
encoder := json.NewEncoder(conn)

// 登录
encoder.Encode(model.Message{Op: enum.Login, Name: "alice"})

// 公屏聊天
encoder.Encode(model.Message{Op: enum.Chat, Name: "alice", Msg: "hello", Area: enum.PublicScreen})
```

## 服务端端口

| 端口 | 服务 |
| --- | --- |
| 8888 | TCP 聊天服务（`server/`） |
| 8080 | HTTP API（`api/`，WebSocket 桥接） |

## 错误处理

当前实现：
- 无 token / session 验证（依赖用户名即可连接）
- 消息体超过 1KB 会被截断
- 非法 JSON 会被忽略
- 心跳：客户端需每 30s 发送一次消息保活（服务端 60s 无消息断开）

## 已知限制

- 不持久化消息历史（重启清空）
- 群组成员关系在内存，重启丢失
- 不支持离线消息
- 私聊无加密
