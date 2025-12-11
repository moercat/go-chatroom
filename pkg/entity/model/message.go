package model

import "go-chatroom/pkg/enum"

type Message struct {
	Name      string         `json:"name"`      // 用户名
	Op        enum.Operation `json:"op"`        // 操作服务
	Msg       string         `json:"msg"`       // 信息内容
	Target    string         `json:"target"`    // 目标用户(私聊时使用)
	Group     string         `json:"group"`     // 群组名称(群聊时使用)
	Timestamp int64          `json:"timestamp"` // 时间戳
	Area      enum.Area      `json:"area"`      // 聊天区域类型
}
