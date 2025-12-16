// 聊天室Web客户端
class ChatClient {
    constructor() {
        this.ws = null;
        this.currentUser = '';
        this.currentTab = 'public';
        this.users = [];
        this.groups = [];
        
        // 绑定方法
        this.initEventListeners();
    }

    initEventListeners() {
        // 登录事件
        document.getElementById('login-btn').addEventListener('click', () => this.login());
        document.getElementById('username').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.login();
        });

        // 聊天界面事件
        document.getElementById('send-btn').addEventListener('click', () => this.sendMessage());
        document.getElementById('message-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.sendMessage();
        });

        // 退出事件
        document.getElementById('logout-btn').addEventListener('click', () => this.logout());

        // 标签页切换
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => this.switchTab(e.target.dataset.tab));
        });

        // 选择私聊对象
        document.getElementById('private-target').addEventListener('change', (e) => {
            if (this.currentTab === 'private') {
                document.getElementById('message-input').focus();
            }
        });

        // 选择群聊组
        document.getElementById('group-target').addEventListener('change', (e) => {
            if (this.currentTab === 'group') {
                document.getElementById('message-input').focus();
            }
        });

        // 创建群组
        document.getElementById('create-group-btn').addEventListener('click', () => this.showCreateGroup());

        // 刷新群组
        document.getElementById('refresh-groups').addEventListener('click', () => this.requestGroups());

        // 修改资料按钮
        document.getElementById('update-profile-btn').addEventListener('click', () => this.showProfileModal());
        document.getElementById('save-profile-btn').addEventListener('click', () => this.updateProfile());
        document.getElementById('cancel-profile-btn').addEventListener('click', () => this.hideProfileModal());

        // 输入焦点
        document.getElementById('new-age').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.updateProfile();
        });
        document.getElementById('new-sex').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.updateProfile();
        });
    }

    async login() {
        const username = document.getElementById('username').value.trim();
        if (!username) {
            alert('请输入用户名');
            return;
        }

        this.currentUser = username;

        try {
            // 连接到WebSocket服务器
            this.connectWebSocket(username);
            
            // 切换到聊天界面
            document.getElementById('login-screen').classList.add('hidden');
            document.getElementById('chat-screen').classList.remove('hidden');
            document.getElementById('current-user').textContent = username;
        } catch (error) {
            console.error('登录失败:', error);
            alert('连接服务器失败，请稍后再试');
        }
    }

    connectWebSocket(username) {
        // WebSocket服务器地址 - 我们将在下一步实现这个后端
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        this.ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

        this.ws.onopen = () => {
            console.log('已连接到聊天服务器');

            // 发送登录消息
            this.sendWsMessage({
                name: username,
                op: 3, // enum.Login
                msg: '',
                area: "public_screen"
            });
        };

        this.ws.onmessage = (event) => {
            // 服务器发送的是格式化的消息，如：【公屏】时间 [用户名]: 内容
            // 我们需要在这里解析这些消息
            const rawMessage = event.data;
            
            // 根据服务器格式解析消息
            const parsedMsg = this.parseServerMessage(rawMessage);
            this.handleReceivedMessage(parsedMsg);
        };

        this.ws.onclose = () => {
            console.log('与服务器断开连接');
            alert('与服务器断开连接');
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
        };
    }

    parseServerMessage(rawMessage) {
        // 解析服务器返回的消息格式
        // 格式如："【公屏】2023-10-01 12:00:00 [Alice]: Hello World"
        let message = rawMessage.trim();
        
        let area = 'public_screen';
        let op = 1; // 默认是聊天消息
        let name = 'SYSTEM';
        let content = message;
        let timestamp = Math.floor(Date.now() / 1000);

        // 确定区域类型
        if (message.startsWith('【公')) {
            area = 'public_screen';
            op = 1;
        } else if (message.startsWith('【群')) {
            area = 'group_chat';
            op = 6; // GroupChat
        } else if (message.startsWith('【私')) {
            area = 'private_chat';
            op = 5; // PrivateChat
        }

        // 解析内容部分
        const parts = message.split('] ');
        if (parts.length >= 2) {
            const header = parts[0]; // 例如 "【公屏】2023-10-01 12:00:00 [Alice"
            content = parts.slice(1).join('] '); // 余下的内容

            // 提取用户名
            const nameStart = header.lastIndexOf('[');
            if (nameStart !== -1) {
                name = header.substring(nameStart + 1); // 提取用户名

                // 提取时间戳
                const timeStart = header.indexOf('】') + 1;
                const timeEnd = header.lastIndexOf('[');
                if (timeStart !== 0 && timeEnd !== -1 && timeEnd > timeStart) {
                    const timeStr = header.substring(timeStart, timeEnd).trim();
                    const date = new Date(timeStr);
                    if (!isNaN(date.getTime())) {
                        timestamp = Math.floor(date.getTime() / 1000);
                    }
                }
            }
        }

        return {
            name: name,
            op: op,
            msg: content,
            area: area,
            timestamp: timestamp
        };
    }

    sendWsMessage(data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.warn('WebSocket未连接');
        }
    }

    sendMessage() {
        const input = document.getElementById('message-input');
        const message = input.value.trim();
        if (!message) return;

        let target = '';
        let group = '';

        switch (this.currentTab) {
            case 'public':
                // 公屏聊天
                this.sendWsMessage({
                    name: this.currentUser,
                    op: 1, // enum.Chat
                    msg: message,
                    area: "public_screen",
                    timestamp: Math.floor(Date.now() / 1000)
                });
                break;
            case 'private':
                // 私聊
                target = document.getElementById('private-target').value;
                if (!target) {
                    alert('请选择私聊对象');
                    return;
                }
                this.sendWsMessage({
                    name: this.currentUser,
                    op: 5, // enum.PrivateChat
                    msg: message,
                    target: target,
                    area: "private_chat",
                    timestamp: Math.floor(Date.now() / 1000)
                });
                break;
            case 'groups':
                // 群聊
                group = document.getElementById('group-target').value;
                if (!group) {
                    alert('请选择群组');
                    return;
                }
                this.sendWsMessage({
                    name: this.currentUser,
                    op: 6, // enum.GroupChat
                    msg: message,
                    group: group,
                    area: "group_chat",
                    timestamp: Math.floor(Date.now() / 1000)
                });
                break;
        }

        input.value = '';
        input.focus();
    }

    handleReceivedMessage(data) {
        // 检查是否是系统消息（包含用户列表或群组列表）
        if (data.msg.includes('当前用户') || data.msg.includes('当前群组')) {
            this.handleListResponse(data);
            return;
        }
        
        // 显示消息
        let chatType = 'public';
        if (data.area === 'private_chat' || data.op === 5) {
            chatType = 'private';
        } else if (data.area === 'group_chat' || data.op === 6) {
            chatType = 'group';
        }
        
        this.displayMessage(data, chatType);
    }

    handleListResponse(data) {
        // 解析响应消息，例如 "当前在线用户列表:\n- user1\n- user2"
        if (data.msg.includes('当前在线用户列表:')) {
            const users = this.parseUserList(data.msg);
            this.updateUsers(users);
        } else if (data.msg.includes('当前群组列表:')) {
            const groups = this.parseGroupList(data.msg);
            this.updateGroups(groups);
        }
    }

    parseUserList(response) {
        // 解析用户列表 "当前在线用户列表:\n- user1\n- user2"
        const lines = response.split('\n');
        const users = [];
        for (let i = 1; i < lines.length; i++) {  // 从第2行开始是用户列表
            const line = lines[i].trim();
            if (line.startsWith('- ')) {
                const username = line.substring(2).trim();
                if (username && username !== this.currentUser) {
                    users.push(username);
                }
            }
        }
        return users;
    }

    parseGroupList(response) {
        // 解析群组列表 "当前群组列表:\n- group1\n- group2"
        const lines = response.split('\n');
        const groups = [];
        for (let i = 1; i < lines.length; i++) {  // 从第2行开始是群组列表
            const line = lines[i].trim();
            if (line.startsWith('- ')) {
                const groupName = line.substring(2).trim();
                if (groupName) {
                    groups.push(groupName);
                }
            }
        }
        return groups;
    }

    displayMessage(data, chatType) {
        let containerId = `${chatType}-chat`;
        if (chatType === 'public') containerId = 'public-chat';
        else if (chatType === 'private') containerId = 'private-chat';
        else if (chatType === 'group') containerId = 'group-chat';

        const container = document.getElementById(containerId);
        if (!container) return;

        const div = document.createElement('div');
        div.className = `message ${data.name === this.currentUser ? 'own' : ''}`;
        
        const timeStr = new Date(data.timestamp * 1000).toLocaleString();
        div.innerHTML = `
            <div class="message-header">
                <span class="message-sender">${this.escapeHtml(data.name)}</span>
                <span class="message-time">${timeStr}</span>
            </div>
            <div class="message-content">${this.escapeHtml(data.msg)}</div>
        `;

        container.appendChild(div);
        container.scrollTop = container.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    switchTab(tabName) {
        // 更新标签页激活状态
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tab === tabName);
        });

        // 显示对应的聊天面板
        document.querySelectorAll('.chat-messages').forEach(panel => {
            panel.classList.toggle('hidden', !panel.id.includes(tabName));
            if (!panel.classList.contains('hidden')) {
                panel.scrollTop = panel.scrollHeight;
            }
        });

        this.currentTab = tabName;

        // 根据当前标签页可能需要执行特定操作
        if (tabName === 'private') {
            this.requestUsers();
        } else if (tabName === 'groups') {
            this.requestGroups();
        }
    }

    requestUsers() {
        this.sendWsMessage({
            name: this.currentUser,
            op: 9, // enum.ListUsers
            msg: '',
            area: "public_screen",
            timestamp: Math.floor(Date.now() / 1000)
        });
    }

    updateUsers(users) {
        this.users = users;
        const userList = document.getElementById('users-list');
        userList.innerHTML = '';
        
        users.forEach(user => {
            const li = document.createElement('li');
            li.textContent = user;
            li.addEventListener('click', () => {
                // 在私聊标签页点击用户，自动选中该用户
                if(this.currentTab === 'private') {
                    document.getElementById('private-target').value = user;
                }
            });
            userList.appendChild(li);
        });

        // 更新私聊下拉列表
        const privateSelect = document.getElementById('private-target');
        privateSelect.innerHTML = '<option value="">选择私聊对象</option>';
        users.filter(u => u !== this.currentUser).forEach(user => {
            const option = document.createElement('option');
            option.value = user;
            option.textContent = user;
            privateSelect.appendChild(option);
        });
    }

    requestGroups() {
        this.sendWsMessage({
            name: this.currentUser,
            op: 8, // enum.ListGroups
            msg: '',
            area: "public_screen",
            timestamp: Math.floor(Date.now() / 1000)
        });
    }

    updateGroups(groups) {
        this.groups = groups;
        const groupsList = document.getElementById('groups-list');
        // 保留刷新按钮
        groupsList.innerHTML = '<li><button id="refresh-groups" class="btn btn-small">刷新群组</button></li>';
        
        groups.forEach(group => {
            const li = document.createElement('li');
            li.textContent = group;
            li.addEventListener('click', () => {
                // 在群聊标签页点击群组，自动选中该群组
                if(this.currentTab === 'groups') {
                    document.getElementById('group-target').value = group;
                }
            });
            groupsList.appendChild(li);
        });

        // 更新群聊下拉列表
        const groupSelect = document.getElementById('group-target');
        groupSelect.innerHTML = '<option value="">选择群组</option>';
        groups.forEach(group => {
            const option = document.createElement('option');
            option.value = group;
            option.textContent = group;
            groupSelect.appendChild(option);
        });

        // 重新绑定刷新按钮事件
        document.getElementById('refresh-groups').addEventListener('click', () => this.requestGroups());
    }

    showCreateGroup() {
        const groupName = prompt('请输入群组名称:');
        if (groupName && groupName.trim()) {
            this.sendWsMessage({
                name: this.currentUser,
                op: 7, // enum.CreateGroup
                msg: groupName.trim(),
                area: "public_screen",
                timestamp: Math.floor(Date.now() / 1000)
            });
        }
    }

    logout() {
        if (this.ws) {
            this.sendWsMessage({
                name: this.currentUser,
                op: 2, // enum.Logout
                msg: '',
                area: "public_screen",
                timestamp: Math.floor(Date.now() / 1000)
            });
            this.ws.close();
        }
        
        // 返回登录界面
        document.getElementById('chat-screen').classList.add('hidden');
        document.getElementById('login-screen').classList.remove('hidden');
        document.getElementById('username').value = this.currentUser;
        this.currentUser = '';
    }

    showProfileModal() {
        document.getElementById('profile-modal').classList.remove('hidden');
        document.getElementById('new-age').focus();
    }

    hideProfileModal() {
        document.getElementById('profile-modal').classList.add('hidden');
        document.getElementById('new-age').value = '';
        document.getElementById('new-sex').value = '';
    }

    updateProfile() {
        const age = document.getElementById('new-age').value.trim();
        const sex = document.getElementById('new-sex').value.trim();

        if (!age && !sex) {
            alert('请至少填写一项资料');
            return;
        }

        // 创建用户对象
        const user = {
            age: age,
            sex: sex
        };

        this.sendWsMessage({
            name: this.currentUser,
            op: 4, // enum.UpdateUser
            msg: JSON.stringify(user),
            area: "public_screen",
            timestamp: Math.floor(Date.now() / 1000)
        });

        this.hideProfileModal();
    }

    updateUserProfile(username, userInfo) {
        if (username === this.currentUser) {
            document.getElementById('user-age').textContent = userInfo.age || '-';
            document.getElementById('user-sex').textContent = userInfo.sex || '-';
        }
    }
}

// 页面加载完成后初始化聊天客户端
document.addEventListener('DOMContentLoaded', () => {
    new ChatClient();
});