package model

import "net"

type Client struct {
	Conn   net.Conn // 连接信息
	Name   string   // 别名
	IsQuit bool     // 是否退出
	User
}
