# AGENT.md

## 项目概述

FileFlow 是一个文件流转与整理服务，接收下载工具（如 qBittorrent、aria2 等）完成下载后的 webhook 回调，自动将文件按用户定义的规则分类、重命名并移动到指定目录。

## 技术栈

| 层面 | 技术 |
|------|------|
| 语言 | Go (1.22+) |
| HTTP 框架 | Echo (`github.com/labstack/echo/v4`) |
| 数据库 | SQLite (`github.com/mattn/go-sqlite3` 或 `modernc.org/sqlite`) |
| 部署 | Docker |

## 核心流程

```
下载工具完成 → HTTP POST webhook → FileFlow 收到回调
    → 提取 file_name + 事件类型
    → 扫描磁盘查找文件（配置的下载目录）
    → 遍历 SQLite 中的规则表进行匹配
    → 按匹配到的规则：重命名 → 移动到目标目录
    → 返回处理结果
```

## 目录结构（规划）

```
FileFlow/
├── main.go               # 入口
├── go.mod / go.sum
├── Dockerfile
├── config/
│   └── config.go         # 配置管理（下载目录、监听端口、数据库路径等）
├── handler/
│   ├── webhook.go        # webhook HTTP handler（接收通知）
│   └── rule.go           # 规则 CRUD handler（可选管理接口）
├── model/
│   ├── file.go           # 文件实体
│   └── rule.go           # 规则实体（表结构定义）
├── service/
│   ├── matcher.go        # 匹配逻辑：按规则匹配文件名
│   ├── mover.go          # 移动逻辑：重命名 + 移动到目标
│   └── scanner.go        # 扫描逻辑：在磁盘上定位文件
├── db/
│   ├── sqlite.go         # SQLite 初始化 & 连接
│   └── migrations.go     # 建表 & 迁移
├── api.http                # API 接口文档
└── AGENT.md              # 本文件
```

## API 接口

### Webhook 回调

**POST** `/api/webhook`

```json
{
  "file_name": "example.mp4",
  "event": "file_moved"
}
```

| 字段 | 说明 |
|------|------|
| `file_name` | 下载完成并移动到 HDD 后的文件名 |
| `event` | 固定为 `"file_moved"` 或 `"file_added"`（可扩展） |

返回：`200 OK` on success，`4xx/5xx` 含错误信息。

## 规则系统

规则存储在 SQLite 的 `rules` 表中，每条规则包含：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER PK | 自增主键 |
| `name` | TEXT | 规则名称（描述） |
| `pattern` | TEXT | 匹配模式（支持通配符/正则） |
| `target_dir` | TEXT | 目标目录路径 |
| `rename_template` | TEXT | 重命名模板（如 `{title}.{ext}`、`{category}/{title}`） |
| `priority` | INTEGER | 优先级（数字越小越优先） |
| `enabled` | BOOL | 是否启用 |
| `created_at` / `updated_at` | DATETIME | 时间戳 |

匹配规则时按优先级降序或升序遍历，首个匹配生效。

## 配置

通过环境变量或配置文件（如 `config.yaml`）提供：

```yaml
server:
  port: 8080

watch_dirs:
  - /data/downloads
  - /data/incoming

db_path: /data/fileflow.db

rules:
  auto_create: true      # 启动时自动创建默认规则库
```

## Docker 部署

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o fileflow .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/fileflow /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["fileflow"]
```

挂载卷：
- `/data` — 数据库文件 + 下载/目标目录
- `/etc/fileflow/config.yaml` — 配置文件

## 开发命令

```bash
go mod init FileFlow
go mod tidy
go run main.go
go build -o fileflow .
```

## 待实现功能

- [x] webhook 接收（`api.md` 已定义）
- [ ] SQLite 数据库初始化与规则迁移
- [ ] 文件系统扫描与定位
- [ ] 规则匹配引擎（glob / regex）
- [ ] 文件重命名与移动操作
- [ ] 错误处理与日志记录
- [ ] Docker 镜像构建与 Compose 编排
- [ ] 规则管理 CRUD API（可选）
