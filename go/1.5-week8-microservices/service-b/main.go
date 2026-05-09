package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"week8-microservices/common/registry"
)

var port = flag.String("port", "8082", "监听端口")

func main() {
	flag.Parse()

	// 1. 连接 etcd（连不上不 crash，降级运行）
	reg, err := registry.New([]string{"localhost:2379"})
	if err != nil {
		log.Printf("⚠️  连接 etcd 失败: %v（将以无注册模式运行）", err)
	}
	if reg != nil {
		defer reg.Close()
	}

	// 2. 注册到 etcd（失败仅警告）
	addr := ":" + *port
	if reg != nil {
		if err := reg.Register("service-b", "localhost"+addr, 10); err != nil {
			log.Printf("⚠️  注册 service-b 失败: %v", err)
		}
	}

	// 3. HTTP 接口
	http.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"service": "service-b",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
			"message": "Hello from Service B!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	fmt.Printf("🚀 Service B 启动在 %s\n", addr)

	// 4. 优雅退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\n🛑 Service B 正在退出...")
		if reg != nil {
			reg.Deregister()
		}
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(addr, nil))
}
