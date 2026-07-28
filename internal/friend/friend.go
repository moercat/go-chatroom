// Package friend 提供好友与群组管理的 HTTP API。
//
// 端点：
//   - POST /api/friend/add           加好友
//   - POST /api/friend/remove        删好友
//   - GET  /api/friend/list          我的好友列表
//   - POST /api/group/join           加入群组
//   - POST /api/group/leave          离开群组
//   - GET  /api/group/list           群组列表
//   - GET  /api/group/:name/members  群成员
//   - GET  /api/message/private?user=xxx  与某人的私聊历史
//   - GET  /api/message/group?name=xxx   群消息历史
//   - GET  /api/message/unread       未读数
//   - POST /api/message/read         全部标已读
//
// 用法：把所有路径挂到 API 路由上。
package friend

import (
	"encoding/json"
	"go-chatroom/pkg/entity/model"
	"go-chatroom/pkg/store"
	"net/http"
	"strconv"
)

type req struct {
	User  string `json:"user"`
	Friend string `json:"friend"`
	Group string `json:"group"`
	Name   string `json:"name"`
}

// Register 注册路由（挂在传入 mux 上）
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/friend/add", handle(func(w http.ResponseWriter, r *http.Request, x req) {
		store.AddFriend(x.User, x.Friend)
		ok(w, "ok")
	}))
	mux.HandleFunc("/api/friend/remove", handle(func(w http.ResponseWriter, r *http.Request, x req) {
		store.RemoveFriend(x.User, x.Friend)
		ok(w, "ok")
	}))
	mux.HandleFunc("/api/friend/list", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		ok(w, store.ListFriends(user))
	})

	mux.HandleFunc("/api/group/join", handle(func(w http.ResponseWriter, r *http.Request, x req) {
		store.JoinGroup(x.Group, x.User)
		ok(w, "ok")
	}))
	mux.HandleFunc("/api/group/leave", handle(func(w http.ResponseWriter, r *http.Request, x req) {
		store.LeaveGroup(x.Group, x.User)
		ok(w, "ok")
	}))
	mux.HandleFunc("/api/group/list", func(w http.ResponseWriter, r *http.Request) {
		ok(w, store.ListGroups())
	})
	mux.HandleFunc("/api/group/members", func(w http.ResponseWriter, r *http.Request) {
		g := r.URL.Query().Get("name")
		ok(w, store.ListGroupMembers(g))
	})

	mux.HandleFunc("/api/message/private", func(w http.ResponseWriter, r *http.Request) {
		a := r.URL.Query().Get("user")
		b := r.URL.Query().Get("with")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		ok(w, store.PrivateHistory(a, b, limit))
	})
	mux.HandleFunc("/api/message/group", func(w http.ResponseWriter, r *http.Request) {
		g := r.URL.Query().Get("name")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		ok(w, store.GroupHistory(g, limit))
	})
	mux.HandleFunc("/api/message/unread", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		ok(w, map[string]int{"count": store.UnreadCount(user)})
	})
	mux.HandleFunc("/api/message/read", handle(func(w http.ResponseWriter, r *http.Request, x req) {
		store.MarkRead(x.User)
		ok(w, "ok")
	}))
}

// 把消息保存到 store（在 api/main.go 的 WebSocket handler 里调）
func SaveMessage(m model.Message) {
	store.Save(m)
}

// 处理 JSON POST 请求的 helper
func handle(fn func(http.ResponseWriter, *http.Request, req)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var x req
		if err := json.NewDecoder(r.Body).Decode(&x); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fn(w, r, x)
	}
}

func ok(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"code": 0,
		"data": data,
	})
}