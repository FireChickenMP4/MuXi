# Week8：注册中心与微服务

基于 etcd 实现服务注册与发现的微服务系统，演示三种调用路径：Gateway→A、Gateway→B、A→B。

## 架构

```
              [API Gateway :8080]
                    │
          ┌─────────┴─────────┐
          │                   │
     [Service A :8081]   [Service B :8082]
          │
     [Service B :8082]   ← A 内部通过 etcd 发现 B 并调用
```

服务启动后各自向 etcd 注册；Gateway 从 etcd 发现所有服务并转发请求。

## 快速开始

### 1. 启动 etcd

```powershell
# 推荐：Docker 一行启动
docker compose up -d
```

### 2. 启动三个服务（各开一个终端）

```powershell
# 方式 A：go run（首次会自动编译）
go run ./service-b/
go run ./service-a/
go run ./gateway/

# 方式 B：直接运行编译好的 .exe
.\service-b.exe
.\service-a.exe
.\gateway.exe

# 如需换端口
.\service-b.exe -port 8083
```

### 3. 测试

**手动测试：**

```powershell
# Gateway → A
curl http://localhost:8080/hello
# {"call_chain":"Gateway → A","message":"Hello from Service A!","service":"service-a"}

# Gateway → B
curl http://localhost:8080/time
# {"message":"Hello from Service B!","time":"...","service":"service-b"}

# Gateway → A → B（A 内部调用 B）
curl http://localhost:8080/hello-with-time
# {"call_chain":"A → B","from_b_time":"...","message":"Hello from Service A!","status":"成功"}
```

**自动集成测试（一键启动 etcd + 跑测试 + 清理）：**

```powershell
# Windows
.\run-test.ps1

# Linux / Git Bash
./run-test.sh

# 保留 etcd 容器反复调试
.\run-test.ps1 -KeepEtcd
# 或
KEEP_ETCD=1 ./run-test.sh
```

集成测试位于 `integration/integration_test.go`（build tag: `integration`），自动编译三个服务为临时二进制，使用独立端口 18080-18082 验证全部三条调用链。需安装 Go 1.25+ 和 Docker。

## API

| 路径                   | 说明                  | 调用链          |
| ---------------------- | --------------------- | --------------- |
| `GET /hello`           | 问候（纯 A）          | Gateway → A     |
| `GET /time`            | 当前时间              | Gateway → B     |
| `GET /hello-with-time` | 问候 + 时间（A 调 B） | Gateway → A → B |

## 项目结构

```
1.5-week8-microservices/
├── go.mod
├── go.sum
├── docker-compose.yml            # Docker 部署 etcd
├── README.md
├── AGENTS.md                     # AI 编码助手指引
├── run-test.ps1                  # Windows 一键测试脚本
├── run-test.sh                   # Linux/macOS 一键测试脚本
├── common/registry/
│   └── etcd.go                   # 注册中心模块
├── integration/
│   └── integration_test.go       # 集成测试
├── service-a/main.go             # Service A
├── service-b/main.go             # Service B
└── gateway/main.go               # API Gateway
```

## 清理

```powershell
# 停止容器
docker compose down

# 停止并清空 etcd 所有数据
docker compose down -v

# 一键测试脚本会自动清理，无需手动操作
```

## 常见问题

**启动报 `bind: address already in use`**

etcd 重启后旧的服务进程可能还占着端口但未重新注册，成为僵尸进程：

```powershell
# 一把杀掉所有相关进程
taskkill /F /IM gateway.exe
taskkill /F /IM service-a.exe
taskkill /F /IM service-b.exe
```

然后重新依次启动 etcd → Service B → Service A → Gateway 即可。

## 关键设计

### 服务注册

- 每个服务启动时调用 `registry.Register(name, addr, ttl)`
- 底层在 etcd 中写入 key `/services/{name}/{addr}`，附带租约
- 后台 goroutine 周期性续约（KeepAlive），TTL=10s；进程退出时注销

### 服务发现

- `registry.Discover(name)` — 从 etcd 实时查询，超时 1s
- `registry.DiscoverCached(name)` — 优先查 etcd，失败则降级到本地缓存，避免 etcd 宕机时阻塞
- 使用前缀查询 `/services/{name}/` 匹配所有注册的实例
- `registry.Watch(name, callback)` 监听实例变化，实时更新本地缓存

### etcd 容灾

- etcd 宕机时，各服务**不会 crash**，仅打印警告并降级运行
- Gateway 和 Service A 通过 `DiscoverCached` 使用上次缓存的地址继续转发
- etcd 恢复后，Watch 自动重连并刷新缓存，服务无缝恢复
- 首次启动时若无 etcd 也无缓存，Gateway 返回 503（无路由表无法转发）

### 多实例重试

同一服务可启动多个实例（不同端口），etcd 会注册全部地址：

```powershell
# 实例 1（默认 8082）
.\service-b.exe

# 实例 2（另一个终端）
.\service-b.exe -port 8083
```

Gateway 和 Service A 转发时**依次尝试所有实例**，一个挂掉自动跳下一个，避免单点故障。

### A → B 服务间调用

Service A 的 `/api/hello-with-time` 处理流程：

1. 调用 `reg.DiscoverCached("service-b")` 获取 B 的地址（优先 etcd，失败用缓存）
2. 通过 HTTP 调用 B 的 `/api/time`
3. 将 B 返回的时间嵌入自己的响应

> 这演示了微服务中**服务间通信**的标准模式：通过注册中心发现下游，而非硬编码地址。

## Week8 学习要点对应

| 课程内容         | 代码体现                                   |
| ---------------- | ------------------------------------------ |
| etcd 注册与发现  | `registry.Register` / `registry.Discover`  |
| 服务心跳（租约） | `client.Grant` + `KeepAlive`               |
| 故障摘除         | 租约过期自动删除 key                       |
| 注册中心容灾     | `DiscoverCached` 缓存兜底，etcd 宕机不阻塞 |
| 微服务拆分       | Gateway、A、B 各自独立部署                 |
| Gateway 模式     | Gateway 作为统一入口转发请求               |
| 配置中心概念     | etcd 作为集中式服务目录                    |
