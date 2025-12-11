package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/enum"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// 连接地址
	host := "localhost"
	// 连接端口
	port := "8000"
	// 拨号创建连接
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	// 连接后通过 defer 以防忘记关闭连接
	defer conn.Close()
	fmt.Println("已连接到", conn.RemoteAddr())

	//	定义基础信息，输入用户昵称
	var baseMsg model.Message
	fmt.Println("请输入用户昵称：")
	_, err = fmt.Scanln(&baseMsg.Name)
	if err != nil {
		fmt.Println("读取用户昵称失败:", err)
		return
	}
	fmt.Println("用户昵称为：", baseMsg.Name)

	go Receive(conn)

	// Send initial login message
	loginMsg := model.Message{
		Name: baseMsg.Name,
		Op:   enum.Login,
		Msg:  "",
		Area: enum.PublicScreen,
	}
	Login(conn, loginMsg)

	// 向服务端发送信息
	scanner := bufio.NewScanner(os.Stdin)

	// 显示菜单选项
	showMenu()

	for {
		fmt.Print("请选择操作: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		var msg = baseMsg

		switch input {
		case "1":
			// 发送公屏消息
			msg.Op = enum.Chat
			msg.Area = enum.PublicScreen
			Say(conn, msg)
		case "2":
			// 发送私聊消息
			msg.Op = enum.PrivateChat
			msg.Area = enum.PrivateArea
			SendPrivateMessage(conn, msg, scanner)
		case "3":
			// 发送群聊消息
			msg.Op = enum.GroupChat
			msg.Area = enum.GroupArea
			SendGroupMessage(conn, msg, scanner)
		case "4":
			// 创建群组
			msg.Op = enum.CreateGroup
			CreateGroup(conn, msg, scanner)
		case "5":
			// 列出群组
			msg.Op = enum.ListGroups
			ListGroups(conn, msg)
		case "6":
			// 列出在线用户
			msg.Op = enum.ListUsers
			ListUsers(conn, msg)
		case "7":
			// 更新用户信息
			msg.Op = enum.UpdateUser
			UpdUser(conn, msg, scanner)
		case "8":
			// 退出
			Quit(conn, msg)
			return
		default:
			fmt.Println("输入无效，请选择正确的选项")
			showMenu()
		}

	}
}

func showMenu() {
	fmt.Println("\n=== 聊天室功能菜单 ===")
	fmt.Println("1 - 发送公屏消息")
	fmt.Println("2 - 发送私聊消息")
	fmt.Println("3 - 发送群聊消息")
	fmt.Println("4 - 创建群组")
	fmt.Println("5 - 查看群组列表")
	fmt.Println("6 - 查看在线用户")
	fmt.Println("7 - 修改个人信息")
	fmt.Println("8 - 退出聊天室")
	fmt.Println("=====================")
}

func Say(conn net.Conn, m model.Message) {
	fmt.Print("请输入想要发送的内容: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("读取消息失败")
		return
	}
	m.Msg = scanner.Text()

	sendMessage(conn, m)
}

func SendPrivateMessage(conn net.Conn, m model.Message, scanner *bufio.Scanner) {
	fmt.Print("请输入目标用户名: ")
	if !scanner.Scan() {
		fmt.Println("读取目标用户名失败")
		return
	}
	m.Target = strings.TrimSpace(scanner.Text())

	fmt.Print("请输入私聊内容: ")
	if !scanner.Scan() {
		fmt.Println("读取消息失败")
		return
	}
	m.Msg = scanner.Text()

	sendMessage(conn, m)
}

func SendGroupMessage(conn net.Conn, m model.Message, scanner *bufio.Scanner) {
	fmt.Print("请输入群组名称: ")
	if !scanner.Scan() {
		fmt.Println("读取群组名称失败")
		return
	}
	m.Group = strings.TrimSpace(scanner.Text())

	fmt.Print("请输入群聊内容: ")
	if !scanner.Scan() {
		fmt.Println("读取消息失败")
		return
	}
	m.Msg = scanner.Text()

	sendMessage(conn, m)
}

func CreateGroup(conn net.Conn, m model.Message, scanner *bufio.Scanner) {
	fmt.Print("请输入要创建的群组名称: ")
	if !scanner.Scan() {
		fmt.Println("读取群组名称失败")
		return
	}
	m.Msg = strings.TrimSpace(scanner.Text())

	sendMessage(conn, m)
}

func ListGroups(conn net.Conn, m model.Message) {
	sendMessage(conn, m)
}

func ListUsers(conn net.Conn, m model.Message) {
	sendMessage(conn, m)
}

func sendMessage(conn net.Conn, m model.Message) {
	// 设置时间戳
	m.Timestamp = time.Now().Unix()

	msgData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msgData)
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}
}

func Quit(conn net.Conn, m model.Message) {
	m.Op = enum.Logout

	msgData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msgData)
	if err != nil {
		fmt.Println("离线失败:", err)
		return
	}
	fmt.Println("离线成功")
}

func Receive(conn net.Conn) {
	for {
		data := make([]byte, 2048) // Increased buffer size
		ml, err := conn.Read(data)
		if ml == 0 || err != nil {
			fmt.Printf("接收数据失败: %v\n", err)
			return // Exit when connection is closed or error occurs
		}
		fmt.Print(string(data[:ml]))
	}
}

func Login(conn net.Conn, m model.Message) {
	// Login 即为我们本地维护的Op表
	msgData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}
	_, err = conn.Write(msgData)
	if err != nil {
		fmt.Println("通知服务端登录信息发送失败:", err)
		return
	}
}

func UpdUser(conn net.Conn, m model.Message, scanner *bufio.Scanner) {
	var user model.User
	fmt.Print("请输入想要的更新的用户年龄: ")
	if !scanner.Scan() {
		fmt.Println("读取年龄失败")
		return
	}
	user.Age = scanner.Text()

	fmt.Print("请输入想要的更新的用户性别: ")
	if !scanner.Scan() {
		fmt.Println("读取性别失败")
		return
	}
	user.Sex = scanner.Text()

	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal user error:", err)
		return
	}

	m.Msg = string(userData)
	m.Op = enum.UpdateUser

	msgData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	_, err = conn.Write(msgData)
	if err != nil {
		fmt.Println("修改失败:", err)
		return
	}
	fmt.Println("修改成功")
}
