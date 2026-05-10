# Week1 - NoSQL

MongoDB 和 PostgreSQL 双后端通用帖子服务，支持嵌套评论和自定义拓展字段。

## 快速开始

```bash
# 启动数据库
docker compose up -d

# 启动服务
go run main.go
```

服务默认监听 `:8080`，前端和 API 同源服务，无需额外配置。
可通过环境变量 `SERVER_PORT`、`MONGO_URI`、`PG_DSN` 配置。

## API

所有 API 通过 `/api/mongo/` 和 `/api/pg/` 前缀分别访问两个后端。

| Method | Path                         | Description            |
| ------ | ---------------------------- | ---------------------- |
| GET    | /api/{db}/posts              | 帖子列表               |
| POST   | /api/{db}/posts              | 创建帖子               |
| GET    | /api/{db}/posts/:id          | 帖子详情（含树形评论） |
| PUT    | /api/{db}/posts/:id          | 更新帖子               |
| DELETE | /api/{db}/posts/:id          | 删除帖子               |
| POST   | /api/{db}/posts/:id/comments | 添加评论               |
| DELETE | /api/{db}/comments/:id       | 删除评论               |

### 创建帖子示例

```json
{
  "title": "Hello",
  "content": "正文",
  "author": "Alice",
  "extensions": {
    "tags": ["go", "nosql"],
    "location": "Beijing"
  }
}
```

### 添加嵌套评论示例

根评论：

```json
{ "content": "好文", "author": "Bob" }
```

回复评论（`parent_id` 嵌套，`reply_to_*` 标记 @某人）：

```json
{
  "content": "谢谢！",
  "author": "Alice",
  "parent_id": "<根评论ID>",
  "reply_to_author": "Bob",
  "reply_to_comment_id": "<根评论ID>"
}
```

## 前端

浏览器打开 `http://localhost:8080/` 即可使用 Web 界面管理帖子（由 Go 服务直接托管，无 CORS 问题）。

## 测试

测试脚本会自动完成：启动 Docker 数据库 -> 启动 Go 服务 -> 运行 API 测试 -> 清理（停止服务 + 关闭 Docker）。

```bash
# PowerShell
./test-api.ps1

# Linux/Mac
./test-api.sh
```

## 设计

|          | MongoDB                                   | PostgreSQL                         |
| -------- | ----------------------------------------- | ---------------------------------- |
| 存储     | 单个文档内嵌 comments 数组                | posts 表 + comments 表（外键）     |
| 扩展字段 | `extensions` Map                          | `extensions` JSONB 列              |
| 评论嵌套 | 文档内 `parent_id`，读取时构建树          | 表自引用 `parent_id`，读取时构建树 |
| @回复    | `reply_to_author` + `reply_to_comment_id` | 同左                               |

## 项目结构

```
main.go          入口
index.html       前端页面
config/          配置
models/          数据模型 + Repository 接口
mongodb/         MongoDB 实现
postgres/        PostgreSQL 实现
handler/         HTTP 处理器
```
