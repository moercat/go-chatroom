package main

import (
	"encoding/json"
	"fmt"
	"go-chatroom/pkg/chat"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/enum"
	"net"
	"time"
)

var ConnMap = make(map[string]model.Client)

func main() {
	// 使用 net 包的 Listen 函数监听 127.0.0.1:8000 上的 tcp 连接
	listen, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("聊天室开启失败！error:%v", err)
		return
	}
	// 使用 defer 在运行结束后优雅的关闭
	defer listen.Close()

	fmt.Println("聊天室开启成功！正在监听8000端口")

	for {
		// 当接收到连接请求时
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("conn fail ...")
			continue
		}
		// conn.RemoteAddr() 连接的客户端地址
		fmt.Println(conn.RemoteAddr(), "connect successed")

		// handle 为每一个客户端开单独的协程进行业务操作
		go handle(conn)
	}

}

func handle(conn net.Conn) {

	for {
		// 通过 Read 获取数据到 data中
		// ml 即为数据长度
		data := make([]byte, 255)
		ml, err := conn.Read(data)
		if ml == 0 || err != nil {
			// 收到的参数错误忽略、
			continue
		}

		// 解析协议
		var cMsg model.Message
		err = json.Unmarshal(data[0:ml], &cMsg)
		if err != nil {
			fmt.Println("json.Unmarshal error:", err)
			return
		}
		name := cMsg.Name

		// 每个人的连接信息
		ConnMap[name] = model.Client{
			Conn: conn,
			Name: name,
		}

		switch cMsg.Op {
		case enum.Chat:
			Read(cMsg)
		case enum.Logout:
			Quit(cMsg)
		case enum.Login:
			ntyLogin(cMsg)
		case enum.UpdateUser:
			UpdUser(cMsg)

		default:
			fmt.Println("无效OP")
		}

	}

}

func Read(m model.Message) {
	fmt.Printf("%v 用户[%s]: %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, m.Msg)

	for _, client := range ConnMap {
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, m.Msg)
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Println("client Conn Error")
			return
		}
	}

}

// 提醒所有人新用户上线
func ntyLogin(m model.Message) {
	for _, client := range ConnMap {
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, "I Login")
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Println("new user Conn Error")
			continue
		}
	}
}

func Quit(m model.Message) {
	fmt.Printf("%v 用户[%s]: 退出 \n", time.Now().Format("2006-01-02 15:04:05"), m.Name)

	// 与上线通知同理，遍历所有链接进行离线通知
	for _, client := range ConnMap {
		// 当找到自己时，关闭与自身的连接且忽略给自己的离线通知
		if client.Name == m.Name {
			client.Conn.Close()
			continue
		}
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, "I Logout")
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Println("client Conn Error")
			return
		}
	}

}

func UpdUser(m model.Message) {
	fmt.Printf("%v 用户[%s]: 修改用户信息 \n", time.Now().Format("2006-01-02 15:04:05"), m.Name)

	var user model.User
	// 解码 msg
	err := json.Unmarshal([]byte(m.Msg), &user)
	if err != nil {
		return
	}

	ConnMap[m.Name] = model.Client{
		Conn:   ConnMap[m.Name].Conn,
		Name:   ConnMap[m.Name].Name,
		IsQuit: true,
		User:   user,
	}

	fmt.Printf("%v 用户[%s]: 用户信息 %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, user)
	msg := fmt.Sprintf("%v 用户[%s]: 用户信息 %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, ConnMap[m.Name].User)
	_, err = ConnMap[m.Name].Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
	if err != nil {
		fmt.Println("client Conn Error")
		return
	}
}
