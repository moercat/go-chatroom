package main

import (
	"encoding/json"
	"fmt"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/enum"
	"net"
)

func main() {
	// 连接地址
	host := "localhost"
	// 连接端口
	port := "8000"
	// 拨号创建连接
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 连接后通过 defer 以防忘记关闭连接
	defer conn.Close()
	fmt.Println("已连接到", conn.RemoteAddr())

	//	定义基础信息，输入用户昵称
	var baseMsg model.Message
	fmt.Println("请输入用户昵称：")
	_, _ = fmt.Scanln(&baseMsg.Name)
	fmt.Println("用户昵称为：", baseMsg.Name)

	go Receive(conn, baseMsg)

	Login(conn, baseMsg)

	// 向服务端发送信息
	for {
		var msg = baseMsg
		fmt.Println("请输入想要的操作Op：")
		_, _ = fmt.Scanln(&msg.Op)

		switch msg.Op {
		case enum.Chat:
			Say(conn, msg)
		case enum.Logout:
			Quit(conn, msg)
		case enum.Login:
			fmt.Println("您已登录，输入无效,请重新输入")
		case enum.UpdateUser:
			UpdUser(conn, msg)
		default:
			fmt.Println("输入无效op,请重新输入")
		}

	}
}

func Say(conn net.Conn, m model.Message) {
	fmt.Println("请输入想要的发送的内容：")
	_, _ = fmt.Scanln(&m.Msg)

	msg, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("发送失败")
		return
	}
	fmt.Println("发送成功")
}

func Quit(conn net.Conn, m model.Message) {

	msg, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("离线失败")
		return
	}
	fmt.Println("离线成功")
}

func Receive(conn net.Conn, m model.Message) {
	for {
		data := make([]byte, 255)
		ml, err := conn.Read(data)
		if ml == 0 || err != nil {
			// 收到的参数错误忽略、
			continue
		}
		fmt.Println(string(data[:ml]))
	}
}

func Login(conn net.Conn, m model.Message) {
	// Login 即为我们本地维护的Op表
	msg, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("通知服务端登录信息发送失败")
		return
	}
}

func UpdUser(conn net.Conn, m model.Message) {
	var user model.User
	fmt.Println("请输入想要的更新的用户年龄：")
	_, _ = fmt.Scanln(&user.Age)
	fmt.Println("请输入想要的更新的用户性别：")
	_, _ = fmt.Scanln(&user.Sex)

	marshal, err := json.Marshal(user)
	if err != nil {
		return
	}

	m.Msg = string(marshal)

	msg, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("修改失败")
		return
	}
	fmt.Println("修改成功")
}
