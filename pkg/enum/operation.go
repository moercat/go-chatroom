package enum

import "strconv"

type Operation int

const (
	Chat Operation = iota + 1
	Logout
	Login
	UpdateUser
	PrivateChat // 私聊
	GroupChat   // 群聊
	CreateGroup // 创建群组
	ListGroups  // 列出群组
	ListUsers   // 列出在线用户
)

func MsgToOperation(msg string) (op Operation) {

	opInt, _ := strconv.Atoi(msg)

	return Operation(opInt)
}
