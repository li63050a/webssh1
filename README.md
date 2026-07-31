# WebSSH - 轻量级 Web 终端管理器

一个开箱即用的 Web 端 SSH 客户端，使用 Go 语言编写，单二进制文件分发，无需安装任何数据库或运行时环境。

## ✨ 功能

- 🔐 **用户系统**：注册 / 登录 / 修改密码，JWT 认证
- 💻 **Web 终端**：基于 xterm.js，支持完整的终端交互
- 🗄️ **零配置数据库**：SQLite 自动创建（`data/data.db`），无需单独安装
- 🔒 **密码安全**：所有密码使用 bcrypt 哈希存储
- 👤 **预置账户**：首次运行自动生成管理员 `admin` (admin123) 和普通用户 `user` (user123)
- 🗑️ **一键重置**：删除 `data/data.db` 后重启，恢复出厂设置
- ⚙️ **外部配置**：`config.yaml` 自动生成，支持自定义端口和 SSL 证书
- 📦 **无依赖分发**：纯 Go 静态编译，单文件运行
- 🌍 **全架构支持**：提供一键编译脚本，覆盖 Linux/macOS/Windows/FreeBSD 等主流平台
- 🔌 **WebSocket 代理**：在浏览器中连接任意 SSH 服务器

## 🚀 快速开始

### 1. 下载二进制文件（推荐）

从 [Releases](https://github.com/li63050a/webssh1/releases) 页面下载对应平台的压缩包，解压后直接运行：

```bash
chmod +x webssh-linux-amd64
./webssh-linux-amd64
```

程序会自动创建 config.yaml 和 data/data.db，并启动 HTTP 服务（默认端口 8080）。

浏览器访问 http://你的服务器IP:8080，使用默认账号登录：

用户名 密码 角色
admin admin123 管理员
user user123 普通用户

2. 从源码编译

需要 Go 1.20+ 环境：

```bash
git clone https://github.com/li63050a/webssh1.git
cd webssh1
go mod tidy
CGO_ENABLED=0 go build -o webssh .
./webssh
```

📖 使用说明

1. 登录：输入默认账号或自行注册后登录
2. 连接 SSH：在终端面板输入目标主机的 IP、端口、用户名和密码，点击 Connect
3. 修改密码：登录后点击顶部 Change Password，输入旧密码和新密码即可更新
4. 注册新用户：在登录界面点击 Register 创建新账号
5. 重置系统：删除 data/data.db 文件后重启，所有用户数据清空，恢复默认账户
6. 启用 HTTPS：编辑 config.yaml，填写证书路径并将端口改为 :443

⚙️ 配置文件

首次运行会在当前目录生成 config.yaml：

```yaml
server:
  addr: ":8080"           # 监听地址
  cert_file: ""           # SSL 证书路径（留空则使用 HTTP）
  key_file: ""            # SSL 私钥路径
  debug: true             # 调试模式

database:
  path: "data/data.db"    # 数据库文件路径

jwt:
  secret: "change-me-in-production"  # JWT 密钥，生产环境请修改
```

🔨 全平台编译

运行 build-all.sh 一键生成所有架构的二进制文件：

```bash
chmod +x build-all.sh
./build-all.sh
```

支持的平台：linux/amd64、linux/arm64、linux/armv6、linux/armv7、linux/386、darwin/amd64、darwin/arm64、windows/amd64、windows/arm64、freebsd/amd64、freebsd/arm64、openbsd/amd64、openbsd/arm64、netbsd/amd64

编译产物存放在 build/ 目录下。

📁 项目结构

```
webssh1/
├── main.go
├── config/
│   └── config.go
├── db/
│   └── db.go
├── handlers/
│   ├── auth.go
│   └── ssh.go
├── static/
│   └── index.html
├── build-all.sh
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

🛠 技术栈

· Go + Gin + gorilla/websocket
· xterm.js
· golang.org/x/crypto (bcrypt + SSH)
· modernc.org/sqlite (纯 Go SQLite)
· golang-jwt/jwt

📜 开源许可

Apache License 2.0，详见 LICENSE 文件。

Copyright 2026 li63050a
