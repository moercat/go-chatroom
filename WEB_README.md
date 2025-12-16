# Go 聊天室 - Web界面

这个项目在原有的命令行聊天室基础上，添加了一个现代化的Web界面，使用户可以通过浏览器访问聊天室。

## 功能特性

- **现代化的用户界面**：响应式设计，支持桌面和移动设备
- **多聊天模式**：
  - 公屏聊天
  - 私聊
  - 群聊
- **用户管理**：
  - 用户名设置
  - 个人信息修改
  - 在线用户列表
  - 群组列表
- **实时通信**：使用WebSocket实现实时消息传输

## 目录结构

```
go-chatroom/
├── client/           # 原始命令行客户端
├── server/           # 原始TCP聊天服务器
├── web/              # Web界面文件
│   ├── index.html    # 主页面
│   ├── style.css     # 样式文件
│   └── script.js     # 前端JavaScript逻辑
├── api/              # Web API服务器
│   ├── main.go       # API服务器代码
│   └── go.mod        # 依赖管理
├── start_web_chatroom.sh  # 一键启动脚本
└── WEB_README.md     # 详细使用说明
```

## 启动方法

### 方法一：使用启动脚本（推荐）

```bash
./start_web_chatroom.sh
```

脚本会自动检查并启动TCP聊天服务器（如果未运行），然后启动Web API服务器，并自动打开浏览器。

### 方法二：手动启动

1. 确保TCP聊天服务器正在运行（监听8000端口）：
```bash
cd server
go run server.go
```

2. 启动Web API服务器：
```bash
cd api
go run main.go
```

3. 打开浏览器访问 `http://localhost:8080`

## 技术架构

- **前端**：HTML5、CSS3、JavaScript
- **后端**：Go语言
- **通信协议**：WebSocket（前端与API服务器），TCP（API服务器与原始聊天服务器）
- **协议转换**：Web API服务器将WebSocket消息转换为TCP JSON消息，并将TCP格式化响应转换为WebSocket消息

## 使用说明

1. 打开浏览器访问 `http://localhost:8080`
2. 输入用户名并点击"进入聊天室"
3. 使用顶部标签页切换不同聊天模式
4. 在侧边栏查看在线用户和群组列表
5. 在底部输入框输入消息并点击"发送"

## 注意事项

- 确保TCP聊天服务器（端口8000）在运行
- Web API服务器默认运行在8080端口
- 现有的命令行客户端仍然可以继续使用
- Web界面与命令行客户端可以同时使用
- Web API服务器会自动将用户消息转换为原始服务器可以理解的格式
- 服务器响应会被解析并转换为Web界面可以显示的格式