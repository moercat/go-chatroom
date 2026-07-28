// sidebar.js — 侧边栏增强：未读数 + 好友列表
// 接入 pkg/store 提供的 REST API
// 通过 fetch 调用 /api/message/unread, /api/friend/list

(function () {
    'use strict';

    // 拉取未读数并更新徽章
    async function refreshUnread(user) {
        try {
            const resp = await fetch(`/api/message/unread?user=${encodeURIComponent(user)}`);
            const data = await resp.json();
            const count = data.data?.count || 0;
            updateUnreadBadge(count);
        } catch (e) {
            console.error('fetch unread failed', e);
        }
    }

    function updateUnreadBadge(count) {
        const badge = document.getElementById('unread-badge');
        if (!badge) return;
        if (count > 0) {
            badge.textContent = count > 99 ? '99+' : count;
            badge.style.display = 'inline-block';
        } else {
            badge.style.display = 'none';
        }
    }

    // 拉取好友列表并填充私聊下拉框
    async function refreshFriends(user) {
        try {
            const resp = await fetch(`/api/friend/list?user=${encodeURIComponent(user)}`);
            const data = await resp.json();
            const friends = data.data || [];
            const select = document.getElementById('private-target');
            if (!select) return;
            // 清空旧选项（保留第一个提示）
            while (select.options.length > 1) select.remove(1);
            friends.forEach(f => {
                const opt = document.createElement('option');
                opt.value = f;
                opt.textContent = f;
                select.appendChild(opt);
            });
        } catch (e) {
            console.error('fetch friends failed', e);
        }
    }

    // 拉取群组列表
    async function refreshGroups() {
        try {
            const resp = await fetch('/api/group/list');
            const data = await resp.json();
            const groups = data.data || [];
            const select = document.getElementById('group-target');
            if (!select) return;
            while (select.options.length > 1) select.remove(1);
            groups.forEach(g => {
                const opt = document.createElement('option');
                opt.value = g;
                opt.textContent = g;
                select.appendChild(opt);
            });
        } catch (e) {
            console.error('fetch groups failed', e);
        }
    }

    // 加载私聊历史
    async function loadPrivateHistory(user, withUser) {
        try {
            const resp = await fetch(`/api/message/private?user=${encodeURIComponent(user)}&with=${encodeURIComponent(withUser)}&limit=50`);
            const data = await resp.json();
            const messages = data.data || [];
            const container = document.getElementById('private-messages');
            if (!container) return;
            container.innerHTML = '';
            messages.forEach(m => appendMessage(container, m));
        } catch (e) {
            console.error('load private history failed', e);
        }
    }

    // 加载群消息历史
    async function loadGroupHistory(group) {
        try {
            const resp = await fetch(`/api/message/group?name=${encodeURIComponent(group)}&limit=50`);
            const data = await resp.json();
            const messages = data.data || [];
            const container = document.getElementById('group-messages');
            if (!container) return;
            container.innerHTML = '';
            messages.forEach(m => appendMessage(container, m));
        } catch (e) {
            console.error('load group history failed', e);
        }
    }

    // 工具：在容器里追加一条消息
    function appendMessage(container, m) {
        const div = document.createElement('div');
        div.className = 'message';
        div.innerHTML = `<span class="msg-name">${escapeHtml(m.name)}</span>: <span class="msg-content">${escapeHtml(m.msg)}</span> <span class="msg-time">${formatTime(m.timestamp)}</span>`;
        container.appendChild(div);
        container.scrollTop = container.scrollHeight;
    }

    function escapeHtml(s) {
        const d = document.createElement('div');
        d.textContent = s || '';
        return d.innerHTML;
    }

    function formatTime(ts) {
        if (!ts) return '';
        const d = new Date(ts * 1000);
        return d.toLocaleTimeString();
    }

    // 暴露到全局（供 script.js 调用）
    window.ChatSidebar = {
        refreshUnread,
        refreshFriends,
        refreshGroups,
        loadPrivateHistory,
        loadGroupHistory,
        markRead: async function (user) {
            await fetch('/api/message/read', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ user }),
            });
            updateUnreadBadge(0);
        },
    };

    // 定时刷新未读数
    setInterval(() => {
        if (window.currentUser) refreshUnread(window.currentUser);
    }, 5000);
})();