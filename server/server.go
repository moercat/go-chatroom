package main

import (
	"encoding/json"
	"fmt"
	"go-chatroom/pkg/chat"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/enum"
	"go-chatroom/pkg/log"
	"net"
	"sync"
	"time"
)

var (
	ConnMap    = make(map[string]model.Client)
	GroupMap   = make(map[string][]string) // 群组成员映射
	groupMutex = sync.RWMutex{}            // 群组操作互斥锁
	mutex      = sync.RWMutex{}
	logger     *log.ChatLogger // 聊天日志记录器
)

func main() {
	// 初始化聊天日志记录器
	var err error
	logger, err = log.NewChatLogger("chat.log")
	if err != nil {
		fmt.Printf("无法创建聊天日志文件！error:%v", err)
		return
	}
	defer logger.Close()

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
	defer conn.Close()

	for {
		// 通过 Read 获取数据到 data中
		// ml 即为数据长度
		data := make([]byte, 1024) // Increased buffer size
		ml, err := conn.Read(data)
		if ml == 0 || err != nil {
			// Connection closed or error occurred
			fmt.Printf("Connection closed or error: %v\n", err)
			return
		}

		// 解析协议
		var cMsg model.Message
		err = json.Unmarshal(data[0:ml], &cMsg)
		if err != nil {
			fmt.Println("json.Unmarshal error:", err)
			continue
		}

		// 设置时间戳
		cMsg.Timestamp = time.Now().Unix()

		name := cMsg.Name

		// 每个人的连接信息 - thread-safe write
		mutex.Lock()
		existingClient, exists := ConnMap[name]
		if exists {
			// 保留现有的用户信息
			ConnMap[name] = model.Client{
				Conn: conn,
				Name: name,
				User: existingClient.User,
			}
		} else {
			ConnMap[name] = model.Client{
				Conn: conn,
				Name: name,
			}
		}
		mutex.Unlock()

		switch cMsg.Op {
		case enum.Chat:
			Read(cMsg)
		case enum.PrivateChat:
			SendPrivateMessage(cMsg)
		case enum.GroupChat:
			SendGroupMessage(cMsg)
		case enum.CreateGroup:
			CreateGroup(cMsg)
		case enum.ListGroups:
			ListGroups(cMsg)
		case enum.ListUsers:
			ListUsers(cMsg)
		case enum.Logout:
			Quit(cMsg)
			return // Exit the handler when client logs out
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

	// 记录公屏消息到日志
	if logger != nil {
		logger.LogMessage(m.Name, "", "", "Public", m.Msg, m.Timestamp)
	}

	mutex.RLock() // Use read lock since we're only reading
	defer mutex.RUnlock()

	// Create a copy of client connections to avoid holding the lock during writes
	clients := make(map[string]model.Client)
	for k, v := range ConnMap {
		clients[k] = v
	}

	for _, client := range clients {
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, m.Msg)
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(m.Area, msg)))
		if err != nil {
			fmt.Printf("client Conn Error for %s: %v\n", client.Name, err)
			// Don't return here, continue sending to other clients
		}
	}

}

// 提醒所有人新用户上线
func ntyLogin(m model.Message) {
	// 记录用户上线事件到日志
	if logger != nil {
		logger.LogMessage(m.Name, "", "", "System", "User Login", m.Timestamp)
	}

	mutex.RLock() // Use read lock since we're only reading
	defer mutex.RUnlock()

	// Create a copy of client connections to avoid holding the lock during writes
	clients := make(map[string]model.Client)
	for k, v := range ConnMap {
		clients[k] = v
	}

	for _, client := range clients {
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, "I Login")
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Printf("new user Conn Error for %s: %v\n", client.Name, err)
			continue
		}
	}
}

func Quit(m model.Message) {
	fmt.Printf("%v 用户[%s]: 退出 \n", time.Now().Format("2006-01-02 15:04:05"), m.Name)

	// 记录用户下线事件到日志
	if logger != nil {
		logger.LogMessage(m.Name, "", "", "System", "User Logout", m.Timestamp)
	}

	// 与上线通知同理，遍历所有链接进行离线通知
	mutex.RLock() // Use read lock since we're only reading
	// Create a copy of client connections to avoid holding the lock during writes
	clients := make(map[string]model.Client)
	for k, v := range ConnMap {
		clients[k] = v
	}
	mutex.RUnlock()

	for _, client := range clients {
		// 当找到自己时，关闭与自身的连接且忽略给自己的离线通知
		if client.Name == m.Name {
			client.Conn.Close()
			continue
		}
		msg := fmt.Sprintf("%v [%s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, "I Logout")
		_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Printf("client Conn Error for %s: %v\n", client.Name, err)
			// Don't return here, continue sending to other clients
		}
	}

	// Remove the client from the connection map
	mutex.Lock()
	delete(ConnMap, m.Name)
	mutex.Unlock()
}

func UpdUser(m model.Message) {
	fmt.Printf("%v 用户[%s]: 修改用户信息 \n", time.Now().Format("2006-01-02 15:04:05"), m.Name)

	var user model.User
	// 解码 msg
	err := json.Unmarshal([]byte(m.Msg), &user)
	if err != nil {
		return
	}

	// Update client info in a thread-safe manner
	mutex.Lock()
	if existingClient, exists := ConnMap[m.Name]; exists {
		ConnMap[m.Name] = model.Client{
			Conn:   existingClient.Conn,
			Name:   existingClient.Name,
			IsQuit: existingClient.IsQuit,
			User:   user,
		}
	}
	mutex.Unlock()

	fmt.Printf("%v 用户[%s]: 用户信息 %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, user)

	// Get the updated client info to send message back
	mutex.RLock()
	client, exists := ConnMap[m.Name]
	mutex.RUnlock()

	if exists {
		msg := fmt.Sprintf("%v 用户[%s]: 用户信息 %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, client.User)
		_, err = client.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, msg)))
		if err != nil {
			fmt.Printf("client Conn Error for %s: %v\n", client.Name, err)
		}
	}
}

// 发送私聊消息
func SendPrivateMessage(m model.Message) {
	fmt.Printf("%v 用户[%s] -> [%s]: %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Name, m.Target, m.Msg)

	// 记录私聊消息到日志
	if logger != nil {
		logger.LogMessage(m.Name, m.Target, "", "Private", m.Msg, m.Timestamp)
	}

	mutex.RLock()
	sender, senderExists := ConnMap[m.Name]
	target, targetExists := ConnMap[m.Target]
	mutex.RUnlock()

	if !senderExists {
		fmt.Printf("发送者 %s 不存在\n", m.Name)
		return
	}

	if !targetExists {
		replyMsg := fmt.Sprintf("用户 %s 不在线或不存在\n", m.Target)
		_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PrivateArea, replyMsg)))
		if err != nil {
			fmt.Printf("向发送者 %s 返回错误信息失败: %v\n", m.Name, err)
		}
		return
	}

	// 构建私聊消息
	privateMsg := fmt.Sprintf("%v [%s -> %s]: %v", time.Now().Format("2006-01-02 15:04:05"), m.Name, m.Target, m.Msg)

	// 发送给目标用户
	_, err := target.Conn.Write([]byte(chat.ShowInOneArea(enum.PrivateArea, privateMsg)))
	if err != nil {
		fmt.Printf("发送私聊消息给 %s 失败: %v\n", m.Target, err)
	}

	// 发送给发送者确认
	_, err = sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PrivateArea, privateMsg)))
	if err != nil {
		fmt.Printf("发送私聊消息给 %s 失败: %v\n", m.Name, err)
	}
}

// 发送群聊消息
func SendGroupMessage(m model.Message) {
	fmt.Printf("%v 群组[%s] 用户[%s]: %v \n", time.Now().Format("2006-01-02 15:04:05"), m.Group, m.Name, m.Msg)

	// 记录群聊消息到日志
	if logger != nil {
		logger.LogMessage(m.Name, "", m.Group, "Group", m.Msg, m.Timestamp)
	}

	// 获取群组成员
	groupMutex.RLock()
	members, exists := GroupMap[m.Group]
	groupMutex.RUnlock()

	if !exists {
		mutex.RLock()
		sender, senderExists := ConnMap[m.Name]
		mutex.RUnlock()

		if senderExists {
			replyMsg := fmt.Sprintf("群组 %s 不存在\n", m.Group)
			_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.GroupArea, replyMsg)))
			if err != nil {
				fmt.Printf("向发送者 %s 返回错误信息失败: %v\n", m.Name, err)
			}
		}
		return
	}

	// 构建群聊消息
	groupMsg := fmt.Sprintf("%v [%s]-%s: %v", time.Now().Format("2006-01-02 15:04:05"), m.Group, m.Name, m.Msg)

	// 发送给群组内所有成员
	mutex.RLock()
	for _, memberName := range members {
		if client, ok := ConnMap[memberName]; ok {
			_, err := client.Conn.Write([]byte(chat.ShowInOneArea(enum.GroupArea, groupMsg)))
			if err != nil {
				fmt.Printf("发送群聊消息给 %s 失败: %v\n", memberName, err)
			}
		}
	}
	mutex.RUnlock()
}

// 创建群组
func CreateGroup(m model.Message) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	// 检查群组是否已存在
	if _, exists := GroupMap[m.Msg]; exists {
		mutex.RLock()
		sender, senderExists := ConnMap[m.Name]
		mutex.RUnlock()

		if senderExists {
			replyMsg := fmt.Sprintf("群组 %s 已存在\n", m.Msg)
			_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
			if err != nil {
				fmt.Printf("向发送者 %s 返回错误信息失败: %v\n", m.Name, err)
			}
		}
		return
	}

	// 创建新群组，将创建者加入该群组
	GroupMap[m.Msg] = []string{m.Name}

	mutex.RLock()
	sender, senderExists := ConnMap[m.Name]
	mutex.RUnlock()

	if senderExists {
		replyMsg := fmt.Sprintf("成功创建群组 %s 并加入该群组\n", m.Msg)
		_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
		if err != nil {
			fmt.Printf("向发送者 %s 返回成功信息失败: %v\n", m.Name, err)
		}
	}
}

// 列出所有群组
func ListGroups(m model.Message) {
	groupMutex.RLock()
	defer groupMutex.RUnlock()

	mutex.RLock()
	sender, senderExists := ConnMap[m.Name]
	mutex.RUnlock()

	if !senderExists {
		return
	}

	if len(GroupMap) == 0 {
		replyMsg := "当前没有可用的群组\n"
		_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
		if err != nil {
			fmt.Printf("向发送者 %s 返回信息失败: %v\n", m.Name, err)
		}
		return
	}

	replyMsg := "当前群组列表:\n"
	for groupName := range GroupMap {
		replyMsg += fmt.Sprintf("- %s\n", groupName)
	}

	_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
	if err != nil {
		fmt.Printf("向发送者 %s 返回群组列表失败: %v\n", m.Name, err)
	}
}

// 列出所有在线用户
func ListUsers(m model.Message) {
	mutex.RLock()
	defer mutex.RUnlock()

	sender, senderExists := ConnMap[m.Name]

	if !senderExists {
		return
	}

	if len(ConnMap) == 0 {
		replyMsg := "当前没有在线用户\n"
		_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
		if err != nil {
			fmt.Printf("向发送者 %s 返回信息失败: %v\n", m.Name, err)
		}
		return
	}

	replyMsg := "当前在线用户列表:\n"
	for userName := range ConnMap {
		replyMsg += fmt.Sprintf("- %s\n", userName)
	}

	_, err := sender.Conn.Write([]byte(chat.ShowInOneArea(enum.PublicScreen, replyMsg)))
	if err != nil {
		fmt.Printf("向发送者 %s 返回用户列表失败: %v\n", m.Name, err)
	}
}
