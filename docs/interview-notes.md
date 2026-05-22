# 后端服务骨架

```aiignore
Go module
    ↓
项目目录分层
    ↓
.env 配置
    ↓
Docker Compose 启动 MySQL
    ↓
GORM 连接 MySQL
    ↓
Gin 注册路由
    ↓
http.Server 启动服务
    ↓
goroutine 后台运行
    ↓
channel 接收错误
    ↓
signal.NotifyContext 监听退出信号
    ↓
context.WithTimeout 控制关闭超时
    ↓
优雅退出
```