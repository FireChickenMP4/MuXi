package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"week8-microservices/common/registry"
)

var (
	reg  *registry.Registry
	port = flag.String("port", "8081", "监听端口")
)

func main() {
	flag.Parse()

	// 1. 连接 etcd（连不上不 crash，降级运行）
	var err error
	reg, err = registry.New([]string{"localhost:2379"})
	if err != nil {
		log.Printf(" 连接 etcd 失败: %v（将以无注册模式运行）", err)
	}
	if reg != nil {
		defer reg.Close()
	}

	// 2. 注册到 etcd（失败仅警告）
	addr := ":" + *port
	if reg != nil {
		if err := reg.Register("service-a", "localhost"+addr, 10); err != nil {
			log.Printf(" 注册 service-a 失败: %v", err)
		}
	}

	// 3. HTTP 接口
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/hello-with-time", helloWithTimeHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	fmt.Printf(" Service A 启动在 %s\n", addr)

	// 4. 优雅退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\n Service A 正在退出...")
		if reg != nil {
			reg.Deregister()
		}
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(addr, nil))
}

// helloHandler /api/hello（纯 A，不调 B）
func helloHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"service":    "service-a",
		"message":    "Hello from Service A!",
		"call_chain": "Gateway → A",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// helloWithTimeHandler /api/hello-with-time（A 通过 etcd 发现 B 并调用）
func helloWithTimeHandler(w http.ResponseWriter, r *http.Request) {
	timeFromB := "未获取到"
	callStatus := "失败"

	if reg != nil {
		addrs := reg.DiscoverCached("service-b") // 优先 etcd，失败用缓存
		if len(addrs) > 0 {
			fmt.Printf("[A→B] 发现 Service B 实例: %v\n", addrs)

			// 遍历所有实例，直到有一个响应成功
			for _, addr := range addrs {
				resp, err := http.Get("http://" + addr + "/api/time")
				if err != nil {
					fmt.Printf("[A→B] 实例 %s 不可达: %v\n", addr, err)
					continue
				}
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)

				var data map[string]interface{}
				json.Unmarshal(body, &data)
				if t, ok := data["time"]; ok {
					timeFromB = fmt.Sprintf("%v", t)
					callStatus = "成功"
				}
				break
			}
		} else {
			fmt.Println("[A→B] 未发现 Service B 实例（etcd 不可用且无缓存）")
		}
	}

	resp := map[string]interface{}{
		"service":     "service-a",
		"message":     "Hello from Service A!",
		"from_b_time": timeFromB,
		"call_chain":  "A → B",
		"status":      callStatus,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
