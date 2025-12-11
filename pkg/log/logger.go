package log

import (
	"fmt"
	"os"
	"time"
)

// ChatLogger 聊天记录器
type ChatLogger struct {
	logFile *os.File
}

// NewChatLogger 创建一个新的聊天记录器
func NewChatLogger(filename string) (*ChatLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &ChatLogger{
		logFile: file,
	}, nil
}

// LogMessage 记录消息
func (cl *ChatLogger) LogMessage(sender, receiver, group, msgType, content string, timestamp int64) {
	t := time.Unix(timestamp, 0)
	logEntry := fmt.Sprintf("[%s] [%s] Sender: %s, Receiver: %s, Group: %s, Type: %s, Content: %s\n",
		t.Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
		sender,
		receiver,
		group,
		msgType,
		content)

	_, err := cl.logFile.WriteString(logEntry)
	if err != nil {
		fmt.Printf("写入日志失败: %v\n", err)
	}
}

// Close 关闭日志文件
func (cl *ChatLogger) Close() {
	if cl.logFile != nil {
		cl.logFile.Close()
	}
}
