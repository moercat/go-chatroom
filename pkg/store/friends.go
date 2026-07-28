// Package store 提供好友关系与群组关系管理。
package store

import (
	"sync"
)

// 好友/群组关系（内存版，生产改 MySQL）
var (
	friendsMu sync.RWMutex
	friends   = make(map[string]map[string]bool) // user -> set of friends

	groupsMu sync.RWMutex
	groups   = make(map[string]map[string]bool) // group -> set of members
)

// AddFriend 加好友（双向）
func AddFriend(a, b string) {
	friendsMu.Lock()
	defer friendsMu.Unlock()
	if friends[a] == nil {
		friends[a] = make(map[string]bool)
	}
	if friends[b] == nil {
		friends[b] = make(map[string]bool)
	}
	friends[a][b] = true
	friends[b][a] = true
}

// RemoveFriend 删好友（双向）
func RemoveFriend(a, b string) {
	friendsMu.Lock()
	defer friendsMu.Unlock()
	delete(friends[a], b)
	delete(friends[b], a)
}

// IsFriend 是否好友
func IsFriend(a, b string) bool {
	friendsMu.RLock()
	defer friendsMu.RUnlock()
	return friends[a][b]
}

// ListFriends 列出用户所有好友
func ListFriends(user string) []string {
	friendsMu.RLock()
	defer friendsMu.RUnlock()
	out := make([]string, 0, len(friends[user]))
	for f := range friends[user] {
		out = append(out, f)
	}
	return out
}

// JoinGroup 加入群组
func JoinGroup(group, user string) {
	groupsMu.Lock()
	defer groupsMu.Unlock()
	if groups[group] == nil {
		groups[group] = make(map[string]bool)
	}
	groups[group][user] = true
}

// LeaveGroup 离开群组
func LeaveGroup(group, user string) {
	groupsMu.Lock()
	defer groupsMu.Unlock()
	delete(groups[group], user)
}

// IsInGroup 是否在群
func IsInGroup(group, user string) bool {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	return groups[group][user]
}

// ListGroupMembers 列出群成员
func ListGroupMembers(group string) []string {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	out := make([]string, 0, len(groups[group]))
	for u := range groups[group] {
		out = append(out, u)
	}
	return out
}

// ListGroups 列出所有群组
func ListGroups() []string {
	groupsMu.RLock()
	defer groupsMu.RUnlock()
	out := make([]string, 0, len(groups))
	for g := range groups {
		out = append(out, g)
	}
	return out
}