# Go 聊天室

Advanced chat room

![example](https://img.shields.io/badge/Go-1.16-blue)

## 目的

本项目采用了简单的客户端服务端架构，实现一个功能丰富的多人在线聊天室

## 特色
从零开始实现一个基于Go的多人在线聊天室，功能包括：单聊、群聊、昵称、上下线通知、聊天日志等等

[X]加入聊天室
[X]广播通知
[X]公屏聊天
[X]聊天日志
[X]群聊
[X]单聊
[X]退出
[X]创建群组
[X]查看群组列表
[X]查看在线用户列表
[X]修改个人资料

## 运行

### 原始命令行客户端

服务端
```shell
cd server
go run server.go
```

客户端(要开启几个客户端就开几个窗口)
```shell
cd client
go run client.go
```

### Web界面（新功能）

现在项目还支持现代化的Web界面，可通过浏览器访问：

1. 启动TCP服务器（同上）
2. 启动Web API服务器：
```shell
cd api
go run main.go
```
3. 访问 http://localhost:8080

#### 使用一键启动脚本

或者使用一键启动脚本：
```shell
./start_web_chatroom.sh
```

该脚本会自动检查并启动TCP聊天服务器（如果未运行），然后启动Web API服务器，并自动打开浏览器。

详细信息请查看 WEB_README.md

## 项目结构

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
└── WEB_README.md     # Web界面详细使用说明
```