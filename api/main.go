package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/enum"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 在生产环境中，应该限制允许的源
			return true
		},
	}

	// 存储WebSocket连接
	connections = make(map[string]*ConnectionPair)
	connMutex   = sync.RWMutex{}

	// TCP服务器地址
	originalServerAddr = "127.0.0.1:8000"
)

// ConnectionPair 存储WebSocket和TCP连接的配对
type ConnectionPair struct {
	WebSocket   *websocket.Conn
	TCPConn     net.Conn
	Name        string
	MessageChan chan model.Message
}

func main() {
	port := flag.String("port", "8080", "HTTP server port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 提供静态文件
		http.ServeFile(w, r, "../web/index.html")
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/style.css")
	})

	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/script.js")
	})

	http.HandleFunc("/ws", handleWebSocket)

	fmt.Printf("Web聊天室服务器启动在端口 %s\n", *port)
	fmt.Printf("请访问 http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 持续监听来自WebSocket的消息
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("读取WebSocket消息失败: %v", err)
			break
		}

		// 解析消息
		var msg model.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("解析JSON消息失败: %v", err)
			continue
		}

		// 根据消息类型处理
		if err := processMessage(&msg, ws); err != nil {
			log.Printf("处理消息失败: %v", err)

			// 发送错误消息回客户端
			errorMsg := model.Message{
				Name:      "SYSTEM",
				Op:        enum.Chat,
				Msg:       fmt.Sprintf("错误: %v", err),
				Area:      enum.PublicScreen,
				Timestamp: time.Now().Unix(),
			}
			ws.WriteJSON(errorMsg)
			continue
		}
	}

	ws.Close()
}

// 处理消息
func processMessage(msg *model.Message, ws *websocket.Conn) error {
	// 获取或创建TCP连接
	connPair, err := getOrCreateConnection(msg.Name, ws)
	if err != nil {
		return err
	}

	// 将消息转发到TCP服务器
	data, err := json.Marshal(*msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	_, err = connPair.TCPConn.Write(data)
	if err != nil {
		return fmt.Errorf("写入TCP连接失败: %v", err)
	}

	return nil
}

// 获取或创建连接对
func getOrCreateConnection(name string, ws *websocket.Conn) (*ConnectionPair, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	connPair, exists := connections[name]
	if exists && connPair.TCPConn != nil {
		// 检查连接是否还有效
		if isValidConn(connPair.TCPConn) {
			// 更新WebSocket连接
			connPair.WebSocket = ws
			return connPair, nil
		}
		// 如果连接无效，关闭并重新创建
		connPair.TCPConn.Close()
		delete(connections, name)
	}

	// 创建新的TCP连接
	tcpConn, err := net.Dial("tcp", originalServerAddr)
	if err != nil {
		return nil, fmt.Errorf("连接TCP服务器失败: %v", err)
	}

	// 创建连接对
	connPair = &ConnectionPair{
		WebSocket:   ws,
		TCPConn:     tcpConn,
		Name:        name,
		MessageChan: make(chan model.Message, 100),
	}
	connections[name] = connPair

	// 启动TCP读取协程
	go startTCPReader(connPair)

	// 发送登录消息
	loginMsg := model.Message{
		Name: name,
		Op:   enum.Login,
		Msg:  "",
		Area: enum.PublicScreen,
	}
	loginData, _ := json.Marshal(loginMsg)
	tcpConn.Write(loginData)

	return connPair, nil
}

// 检查连接是否有效
func isValidConn(conn net.Conn) bool {
	// 尝试写入少量数据来检测连接状态
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err := conn.Write([]byte{0}) // 发送一个字节
	// 如果写入失败，连接可能已断开
	return err == nil
}

// 开始监听TCP连接
func startTCPReader(connPair *ConnectionPair) {
	defer func() {
		connPair.TCPConn.Close()

		connMutex.Lock()
		if existing, ok := connections[connPair.Name]; ok && existing == connPair {
			delete(connections, connPair.Name)
		}
		connMutex.Unlock()
	}()

	reader := bufio.NewReader(connPair.TCPConn)
	for {
		// 读取直到换行符
		line, err := reader.ReadString('\n')
		if err != nil {
			if connPair.WebSocket != nil {
				// 发送断开连接消息到客户端
				disconnectMsg := model.Message{
					Name:      "SYSTEM",
					Op:        enum.Chat,
					Msg:       "服务器连接已断开",
					Area:      enum.PublicScreen,
					Timestamp: time.Now().Unix(),
				}
				connPair.WebSocket.WriteJSON(disconnectMsg)
			}
			break
		}

		// 解析服务器返回的消息（格式为：【区域】时间 [用户名]: 内容）
		parsedMsg := parseServerMessage(line, connPair.Name)

		// 发送消息到WebSocket
		if connPair.WebSocket != nil {
			err := connPair.WebSocket.WriteJSON(parsedMsg)
			if err != nil {
				log.Printf("发送消息到WebSocket失败: %v", err)
				break
			}
		}
	}
}

// 解析服务器返回的消息格式
// 格式如："【公屏】2023-10-01 12:00:00 [Alice]: Hello World"
func parseServerMessage(serverLine, senderName string) model.Message {
	// 去除可能的换行符
	serverLine = strings.Trim(serverLine, "\n\r")

	var area enum.Area
	var op enum.Operation
	var name, content string
	var timestamp int64

	// 确定消息区域
	if strings.HasPrefix(serverLine, "【公") {
		area = enum.PublicScreen
		op = enum.Chat
	} else if strings.HasPrefix(serverLine, "【群") {
		area = enum.GroupArea
		op = enum.GroupChat
	} else if strings.HasPrefix(serverLine, "【私") {
		area = enum.PrivateArea
		op = enum.PrivateChat
	} else {
		area = enum.PublicScreen
		op = enum.Chat
	}

	// 解析内容
	// 格式："【区域】时间 [用户名]: 内容"
	parts := strings.SplitN(serverLine, "] ", 2)
	if len(parts) == 2 {
		// 提取用户名
		headerPart := parts[0] // 例如 "【公屏】2023-10-01 12:00:00 [Alice"
		content = parts[1]     // 消息内容

		// 提取用户名
		nameStart := strings.LastIndex(headerPart, "[")
		if nameStart != -1 {
			name = headerPart[nameStart+1:] // 提取用户名
		} else {
			name = "SYSTEM"
		}

		// 提取时间（格式：2006-01-02 15:04:05）
		timeStart := strings.Index(headerPart, "】") + 1
		timeEnd := strings.LastIndex(headerPart, "[")
		if timeStart != 0 && timeEnd != -1 && timeEnd > timeStart {
			timeStr := strings.TrimSpace(headerPart[timeStart:timeEnd])
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
				timestamp = parsedTime.Unix()
			} else {
				timestamp = time.Now().Unix() // 如果解析失败，使用当前时间
			}
		} else {
			timestamp = time.Now().Unix()
		}
	} else {
		// 如果格式不匹配，使用整个行作为内容
		name = "SYSTEM"
		content = serverLine
		timestamp = time.Now().Unix()
	}

	// 创建消息对象
	return model.Message{
		Name:      name,
		Op:        op,
		Msg:       content,
		Area:      area,
		Timestamp: timestamp,
	}
}
