# go-chatroom Makefile
# 快速命令：
#   make            显示帮助
#   make web        一键启动 TCP + HTTP API + Web（推荐）
#   make tcp        只启动 TCP 服务
#   make api        只启动 HTTP API
#   make test       跑测试
#   make build      编译
#   make clean      清理构建产物

GO       ?= go
SERVER   := ./server/server.go
API      := ./api/main.go

.PHONY: help
help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: web
web: ## 一键启动：TCP server + HTTP API + 打开 Web
	bash start_web_chatroom.sh

.PHONY: tcp
tcp: ## 启动 TCP 服务（:8888）
	cd server && $(GO) run server.go

.PHONY: api
api: ## 启动 HTTP API（:8080）
	cd api && $(GO) run main.go

.PHONY: build
build: ## 编译所有
	cd server && $(GO) build -o chatroom-server .
	cd api && $(GO) build -o chatroom-api .

.PHONY: test
test: ## 跑测试
	$(GO) test ./...

.PHONY: vet
vet: ## 静态检查
	$(GO) vet ./...

.PHONY: clean
clean: ## 清理
	rm -f server/chatroom-server api/chatroom-api server/chat.log