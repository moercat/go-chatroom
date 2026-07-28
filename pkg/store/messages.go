// Package store 提供消息持久化（内存版）。
//
// 设计：
//   - 内存环形缓冲，重启清空
//   - 按 (from, to) 索引支持私聊历史
//   - 按 group 索引支持群聊历史
//   - 未读计数：每个用户对每条消息有 read 标志
//
// 生产可替换为：MySQL / SQLite / MongoDB。
package store

import (
	"go-chatroom/pkg/entity/model"
	"sync"
	"time"
)

// 内存存储
var (
	mu        sync.RWMutex
	messages  []model.Message           // 全部消息
	byPrivate map[string][]model.Message // key: "alice|bob"
	byGroup   map[string][]model.Message // key: "tech-team"
	readState map[string]map[int64]bool  // key: "user" -> msgID -> read
)

func init() {
	byPrivate = make(map[string][]model.Message)
	byGroup = make(map[string][]model.Message)
	readState = make(map[string]map[int64]bool)
}

// Save 存一条消息
func Save(m model.Message) {
	mu.Lock()
	defer mu.Unlock()

	m.Timestamp = time.Now().Unix()
	messages = append(messages, m)

	switch m.Area {
	case "private_chat":
		key := privateKey(m.Name, m.Target)
		byPrivate[key] = append(byPrivate[key], m)
	case "group_chat":
		byGroup[m.Group] = append(byGroup[m.Group], m)
	}
}

// PrivateHistory 两人之间的私聊记录
func PrivateHistory(a, b string, limit int) []model.Message {
	mu.RLock()
	defer mu.RUnlock()
	key := privateKey(a, b)
	list := byPrivate[key]
	if limit > 0 && limit < len(list) {
		list = list[len(list)-limit:]
	}
	return list
}

// GroupHistory 群组消息历史
func GroupHistory(group string, limit int) []model.Message {
	mu.RLock()
	defer mu.RUnlock()
	list := byGroup[group]
	if limit > 0 && limit < len(list) {
		list = list[len(list)-limit:]
	}
	return list
}

// UnreadCount 用户的未读消息数
func UnreadCount(user string) int {
	mu.RLock()
	defer mu.RUnlock()
	c := 0
	for _, m := range messages {
		if m.Name == user {
			continue // 自己发的
		}
		if isForUser(m, user) {
			if !readState[user][m.Timestamp] {
				c++
			}
		}
	}
	return c
}

// MarkRead 标记用户所有消息已读
func MarkRead(user string) {
	mu.Lock()
	defer mu.Unlock()
	if readState[user] == nil {
		readState[user] = make(map[int64]bool)
	}
	for _, m := range messages {
		if isForUser(m, user) && m.Name != user {
			readState[user][m.Timestamp] = true
		}
	}
}

func isForUser(m model.Message, user string) bool {
	switch m.Area {
	case "public_screen", "public":
		return true
	case "private_chat":
		return m.Name == user || m.Target == user
	case "group_chat":
		// 真正判断群成员资格（依赖 store.IsInGroup）
		if m.Group == "" {
			return false
		}
		return IsInGroup(m.Group, user)
	}
	return false
}

func privateKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}