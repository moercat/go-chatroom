package model

import "go-chatroom/pkg/enum"

type Message struct {
	Name string         // 用户名
	Op   enum.Operation // 操作服务
	Msg  string         // 信息内容
}
