# AGENTS.md -- Week8 Microservices (Go + etcd)

## 项目概述

基于 etcd 的服务注册与发现微服务演示项目。Go 1.25.5 模块化单体仓库。

## 构建与运行

### 常用命令

```powershell
# 编译所有服务
go build ./gateway/
go build ./service-a/
go build ./service-b/

# 运行（自动编译）
go run ./gateway/
go run ./service-a/
go run ./service-b/

# 指定端口运行
go run ./gateway/ -port 8080
go run ./service-a/ -port 8081
go run ./service-b/ -port 8082

# 启动 etcd
docker compose up -d
docker compose down            # 停止
docker compose down -v         # 停止并清空数据

# 编译并生成 .exe
go build -o gateway.exe ./gateway/
go build -o service-a.exe ./service-a/
go build -o service-b.exe ./service-b/

# 获取依赖
go mod tidy

# 检查依赖
go mod verify
```

### 测试

```powershell
# 一键方式：启动 etcd -> 测试 -> 停止 etcd
.\run-test.ps1

# 保留 etcd 容器以便反复调试
.\run-test.ps1 -KeepEtcd

# 运行集成测试（需先启动 etcd）
go test -tags=integration -v -timeout 60s ./integration/

# 跳过集成测试（仅跑单元测试）
go test -v ./...

# 竞态检测
go test -race -tags=integration ./integration/

# 运行单个包的所有测试
go test ./common/registry/
```

集成测试位于 `integration/integration_test.go` 中，使用 build tag `integration` 隔离：
- 自动编译 service-a/service-b/gateway 为临时二进制
- 使用独立端口（18080-18082）避免冲突
- 测试三个端点：/hello, /time, /hello-with-time
- 验证 JSON 响应的关键字段完整性
- defer 清理进程和临时文件

**注意：** 集成测试依赖 etcd，未检测到 etcd 时会自动 Skip。

### 代码质量

```powershell
# 格式化所有 Go 文件
gofmt -w .

# 静态检查（需要安装：go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest）
go vet ./...

# 更多检查（需要安装 staticcheck）
# go install honnef.co/go/tools/cmd/staticcheck@latest
# staticcheck ./...
```

## 项目结构

```
week8-microservices/
├── go.mod                      # Go 模块定义 (module week8-microservices)
├── go.sum                      # 依赖锁文件
├── docker-compose.yml          # etcd 容器配置
├── README.md
├── common/registry/
│   └── etcd.go                 # Registry 封装 (etcd 注册/发现/监听/缓存)
├── service-a/
│   └── main.go                 # 服务 A :8081
├── service-b/
│   └── main.go                 # 服务 B :8082
└── gateway/
    └── main.go                 # API Gateway :8080
```

## 代码风格指南

### 导入

- 标准库、第三方、内部包的导入分组，每组之间空行
- 内部包使用模块前缀 `week8-microservices/`
- 使用 `gofmt` 自动排序导入

```go
import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	clientv3 "go.etcd.io/etcd/client/v3"

	"week8-microservices/common/registry"
)
```

### 命名

| 类别 | 风格 | 示例 |
|------|------|------|
| 包名 | 小写、无下划线 | `registry`, `main` |
| 导出类型 | PascalCase | `Registry`, `New` |
| 未导出类型/变量 | camelCase | `reg`, `leaseID`, `cancel`, `mu` |
| 常量 | PascalCase | (项目无常量) |
| 函数 | PascalCase (导出), camelCase (未导出) | `New`, `Discover`, `tryForward` |
| 接收者 | 1-2 字母缩写 | `r *Registry` |
| 文件命名 | snake_case | `etcd.go` |

### 格式化

- 必须使用 `gofmt` 标准化
- 缩进：tab（Go 标准）
- 行宽：无硬限制，建议不超过 120 字符
- 结构体字段对齐（gofmt 处理）

### 类型

- 优先使用 `interface{}` 而非 `any`（项目惯例）
- 错误作为返回值的最后一个值
- 结构体用 `sync.RWMutex` 保护并发访问

### 错误处理

```go
// 正确：错误链式包装
if err != nil {
    return fmt.Errorf("do thing: %w", err)
}

// 错误：裸返回
// return err  // 不提供上下文
```

- 所有 error 必须被检查或显式忽略（`_ =`）
- 错误消息以大写字母开头或具体描述（项目惯例：中文错误说明）
- `defer` 按需使用，确保资源释放

### 并发安全

```go
type Registry struct {
    mu        sync.RWMutex
    instances map[string][]string
}

// 读锁 - 多个 goroutine 可同时读
r.mu.RLock()
cached := r.instances[name]
r.mu.RUnlock()

// 写锁 - 独占
r.mu.Lock()
r.instances[name] = addrs
r.mu.Unlock()
```

### HTTP Handler 模式

```go
// 函数式 handler
http.HandleFunc("/api/hello", helloHandler)

// 内联 handler（短逻辑）
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ok"))
})

// 始终设置 Content-Type
w.Header().Set("Content-Type", "application/json")
```

### 容错设计（项目约定）

- etcd 连接失败：不 crash，仅 log 警告，降级运行
- 服务发现失败：用本地缓存 (`DiscoverCached`) 兜底
- 服务调用失败：遍历多实例重试
- 指针为 nil 时用 `if reg != nil` 守卫

### flag 使用

```go
var port = flag.String("port", "8081", "监听端口")

func main() {
    flag.Parse()
    addr := ":" + *port
}
```

### goroutine 管理

```go
// 启动 goroutine
go func() {
    for range keepAliveCh {
        // 消费通道
    }
}()

// 优雅退出：信号监听
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
<-sigCh
// 执行清理...
os.Exit(0)
```

### 依赖

| 依赖 | 用途 |
|------|------|
| `go.etcd.io/etcd/client/v3` | etcd 客户端 |
| `google.golang.org/grpc` | (传递) etcd gRPC 通信 |
| `go.uber.org/zap` | (传递) etcd 内部日志 |

添加新依赖前确认是否真正必要，优先使用标准库。

### JSON 响应

```go
resp := map[string]interface{}{
    "service": "service-a",
    "message": "Hello from Service A!",
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(resp)
```

不要手动拼接 JSON 字符串，始终使用 `encoding/json`。

### 全局变量

仅在 `main` 包中使用全局变量（服务级别的共享状态），`registry` 等库包中使用结构体封装。
