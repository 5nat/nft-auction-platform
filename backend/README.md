# Backend

Go 后端服务，负责 NFT 拍卖系统的链下能力，包括 HTTP API、数据库存储、链上事件索引和后续交易服务。

## 本地启动

### 1. 启动 MySQL

```bash
docker compose up -d
```

### 2. 启动后端
```bash
go run ./cmd/server
```

### 3. 健康检查
```bash
curl http://localhost:8080/health
```

正常返回：
```json
{
    "status": "ok",
    "database": "ok"
}
```

### 常用命令
```aiignore
docker compose ps
docker compose logs mysql
docker compose down
docker compose down -v
go mod tidy
go run ./cmd/server
```
