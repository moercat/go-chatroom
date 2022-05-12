package main

import (
	"fmt"
	"net"
	"strconv"
)

const (
	Say = iota + 1
	Quit
	Login
)

type Message struct {
	Name string // 用户名
	Op   int    // 操作服务
	Msg  string // 信息内容
}

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
	var baseMsg Message
	fmt.Println("请输入用户昵称：")
	_, _ = fmt.Scanln(&baseMsg.Name)
	fmt.Println("用户昵称为：", baseMsg.Name)

	go baseMsg.Receive(conn)

	baseMsg.Login(conn)

	// 向服务端发送信息
	for {
		var msg = baseMsg
		fmt.Println("请输入想要的操作Op：")
		_, _ = fmt.Scanln(&msg.Op)

		switch msg.Op {
		case Say:
			msg.Say(conn)
		case Quit:
			msg.Quit(conn)
		case Login:
			fmt.Println("您已登录，输入无效,请重新输入")
		default:
			fmt.Println("输入无效op,请重新输入")
		}

	}
}

func (m Message) Say(conn net.Conn) {
	fmt.Println("请输入想要的发送的内容：")
	_, _ = fmt.Scanln(&m.Msg)

	msg := m.Name + "|" + strconv.Itoa(m.Op) + "|" + m.Msg

	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("发送失败")
		return
	}
	fmt.Println("发送成功")
}

func (m Message) Quit(conn net.Conn) {
	fmt.Println("quit 方法")
}

func (m Message) Receive(conn net.Conn) {
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

func (m Message) Login(conn net.Conn) {
	// Login 即为我们本地维护的Op表
	msg := m.Name + "|" + strconv.Itoa(Login) + "|"

	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("通知服务端登录信息发送失败")
		return
	}
}
