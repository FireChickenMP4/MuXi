package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"week8-microservices/common/registry"
)

var (
	mu            sync.RWMutex
	serviceAAddrs []string
	serviceBAddrs []string
	reg           *registry.Registry
	port          = flag.String("port", "8080", "监听端口")
)

func setServiceAAddrs(addrs []string) {
	mu.Lock()
	serviceAAddrs = addrs
	mu.Unlock()
}

func setServiceBAddrs(addrs []string) {
	mu.Lock()
	serviceBAddrs = addrs
	mu.Unlock()
}

func getServiceAAddrs() []string {
	mu.RLock()
	defer mu.RUnlock()
	return serviceAAddrs
}

func getServiceBAddrs() []string {
	mu.RLock()
	defer mu.RUnlock()
	return serviceBAddrs
}

func main() {
	flag.Parse()

	// 1. 连接 etcd（连不上不 crash，用缓存兜底）
	var err error
	reg, err = registry.New([]string{"localhost:2379"})
	if err != nil {
		log.Printf("⚠️  连接 etcd 失败: %v（将以纯缓存模式运行）", err)
	}
	if reg != nil {
		defer reg.Close()
	}

	// 2. 初始发现 + 监听变化
	refreshAddrs()

	if reg != nil {
		reg.Watch("service-a", func(addrs []string) {
			setServiceAAddrs(addrs)
			fmt.Printf("🔄 service-a 实例变更: %v\n", addrs)
		})
		reg.Watch("service-b", func(addrs []string) {
			setServiceBAddrs(addrs)
			fmt.Printf("🔄 service-b 实例变更: %v\n", addrs)
		})
	}

	// 3. 路由
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Gateway→A] 收到请求 /hello")
		proxyToA(w, r, "/api/hello")
	})
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Gateway→B] 收到请求 /time")
		proxyToB(w, r)
	})
	http.HandleFunc("/hello-with-time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Gateway→A→B] 收到请求 /hello-with-time")
		proxyToA(w, r, "/api/hello-with-time")
	})

	addr := ":" + *port
	fmt.Printf("🚀 Gateway 启动在 %s\n", addr)
	fmt.Println("可用路径:")
	fmt.Println("  GET /hello           → Gateway → A")
	fmt.Println("  GET /time            → Gateway → B")
	fmt.Println("  GET /hello-with-time → Gateway → A → B")

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\n🛑 Gateway 正在退出...")
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(addr, nil))
}

func refreshAddrs() {
	if reg == nil {
		fmt.Println("📋 etcd 不可用，无缓存")
		return
	}
	a := reg.DiscoverCached("service-a")
	b := reg.DiscoverCached("service-b")
	setServiceAAddrs(a)
	setServiceBAddrs(b)
	fmt.Printf("📋 当前服务: service-a=%v, service-b=%v\n", a, b)
}

// proxyToA 转发到 Service A，支持多实例重试
func proxyToA(w http.ResponseWriter, r *http.Request, path string) {
	addrs := getServiceAAddrs()
	if len(addrs) == 0 {
		refreshAddrs()
		addrs = getServiceAAddrs()
	}
	for _, addr := range addrs {
		if tryForward(w, addr, path) == nil {
			return
		}
	}
	http.Error(w, `{"error":"Service A 不可用"}`, http.StatusServiceUnavailable)
}

// proxyToB 转发到 Service B，支持多实例重试
func proxyToB(w http.ResponseWriter, r *http.Request) {
	addrs := getServiceBAddrs()
	if len(addrs) == 0 {
		refreshAddrs()
		addrs = getServiceBAddrs()
	}
	for _, addr := range addrs {
		if tryForward(w, addr, "/api/time") == nil {
			return
		}
	}
	http.Error(w, `{"error":"Service B 不可用"}`, http.StatusServiceUnavailable)
}

// tryForward 尝试转发到一个实例，成功返回 nil，失败返回 error
func tryForward(w http.ResponseWriter, target string, path string) error {
	url := "http://" + target + path
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[Gateway] 实例 %s 不可达: %v\n", target, err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	return nil
}
