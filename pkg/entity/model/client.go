package model

import "net"

type Client struct {
	Conn   net.Conn      `json:"-"`       // 连接信息 - skip JSON serialization
	Name   string        `json:"name"`    // 别名
	IsQuit bool          `json:"is_quit"` // 是否退出
	User   `json:"user"` // Embedded user info
}
