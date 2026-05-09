//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func isEtcdReachable() bool {
	client := http.Client{Timeout: 1 * time.Second}
	// etcd v3 无简单 HTTP 健康端点，尝试拨 TCP 判断
	conn, err := client.Get("http://localhost:2379/version")
	return err == nil && conn != nil
}

func buildService(t *testing.T, dir, name string) string {
	t.Helper()
	bin := filepath.Join(os.TempDir(), fmt.Sprintf("test-%s-%d.exe", name, time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", bin, fmt.Sprintf("../%s/", name))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s 失败: %v\n%s", name, err, out)
	}
	return bin
}

func TestIntegration(t *testing.T) {
	if !isEtcdReachable() {
		t.Skip("etcd 未运行在 localhost:2379，请执行 docker compose up -d")
	}

	// 1. 编译服务
	serviceBBin := buildService(t, "service-b", "service-b")
	serviceABin := buildService(t, "service-a", "service-a")
	gatewayBin := buildService(t, "gateway", "gateway")

	// 2. 启动服务
	serviceB := exec.Command(serviceBBin, "-port", "18082")
	serviceA := exec.Command(serviceABin, "-port", "18081")
	gateway := exec.Command(gatewayBin, "-port", "18080")

	serviceB.Stdout = os.Stdout
	serviceB.Stderr = os.Stderr
	serviceA.Stdout = os.Stdout
	serviceA.Stderr = os.Stderr
	gateway.Stdout = os.Stdout
	gateway.Stderr = os.Stderr

	mustStart := func(name string, cmd *exec.Cmd) {
		if err := cmd.Start(); err != nil {
			t.Fatalf("启动 %s 失败: %v", name, err)
		}
	}

	mustStart("service-b", serviceB)
	mustStart("service-a", serviceA)
	mustStart("gateway", gateway)

	// 3. 清理
	defer func() {
		for _, cmd := range []*exec.Cmd{gateway, serviceA, serviceB} {
			if cmd.Process != nil {
				cmd.Process.Signal(os.Interrupt)
			}
		}
		time.Sleep(500 * time.Millisecond)
		for _, cmd := range []*exec.Cmd{gateway, serviceA, serviceB} {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}
		for _, bin := range []string{serviceBBin, serviceABin, gatewayBin} {
			os.Remove(bin)
		}
	}()

	// 4. 等待注册
	time.Sleep(3 * time.Second)

	// 5. 测试三个端点
	tests := []struct {
		name     string
		url      string
		wantKeys []string
	}{
		{
			name:     "Gateway -> A: /hello",
			url:      "http://localhost:18080/hello",
			wantKeys: []string{"service", "message", "call_chain"},
		},
		{
			name:     "Gateway -> B: /time",
			url:      "http://localhost:18080/time",
			wantKeys: []string{"service", "time", "message"},
		},
		{
			name:     "Gateway -> A -> B: /hello-with-time",
			url:      "http://localhost:18080/hello-with-time",
			wantKeys: []string{"service", "message", "from_b_time", "status"},
		},
	}

	client := http.Client{Timeout: 5 * time.Second}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Get(tt.url)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("状态码 %d, 期望 %d", resp.StatusCode, http.StatusOK)
			}

			var body map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("JSON 解析失败: %v", err)
			}

			for _, key := range tt.wantKeys {
				if _, ok := body[key]; !ok {
					t.Errorf("缺少字段 %q, 响应: %v", key, body)
				}
			}

			t.Logf("响应: %v", body)
		})
	}
}
