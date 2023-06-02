package enum

import "strconv"

type Operation int

const (
	Chat Operation = iota + 1
	Logout
	Login
	UpdateUser
)

func MsgToOperation(msg string) (op Operation) {

	opInt, _ := strconv.Atoi(msg)

	return Operation(opInt)
}
