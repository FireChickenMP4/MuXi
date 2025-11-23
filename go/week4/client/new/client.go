package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func main() {
	// 创建带Cookie管理的HTTP客户端
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	// ===== 1. 测试注册 =====
	registerData := url.Values{
		"username": {"testuser"},
		"password": {"securepass123"},
		"nickname": {"测试用户"},
		"email":    {"test@example.com"},
	}
	testAPI(client, "注册", "POST", "http://localhost:9090/register", registerData)

	// ===== 2. 测试登录 =====
	loginData := url.Values{
		"username": {"testuser"},
		"password": {"securepass123"},
	}
	testAPI(client, "登录", "POST", "http://localhost:9090/login", loginData)

	// ===== 3. 测试获取用户信息 =====
	testAPI(client, "获取用户信息", "GET", "http://localhost:9090/whoami", nil)

	// ===== 4. 测试修改密码 =====
	passData := url.Values{
		"username":    {"testuser"},
		"password":    {"securepass123"},
		"newpassword": {"newpass456"},
	}
	testAPI(client, "修改密码", "POST", "http://localhost:9090/password", passData)

	// ===== 5. 测试修改用户名 =====
	userData := url.Values{
		"username":    {"testuser"},
		"password":    {"newpass456"}, // 使用新密码
		"newusername": {"freshuser"},
	}
	testAPI(client, "修改用户名", "POST", "http://localhost:9090/username", userData)

	// ===== 6. 重新登录新用户名 =====
	newLoginData := url.Values{
		"username": {"freshuser"},
		"password": {"newpass456"},
	}
	testAPI(client, "新账号登录", "POST", "http://localhost:9090/login", newLoginData)

	// ===== 7. 测试修改用户信息 =====
	profileData := url.Values{
		"nickname": {"新昵称"},
		"email":    {"new@example.com"},
	}
	testAPI(client, "修改用户信息", "POST", "http://localhost:9090/profile", profileData)

	// ===== 8. 最终验证 =====
	testAPI(client, "最终验证", "GET", "http://localhost:9090/whoami", nil)
}

// 通用API测试函数
func testAPI(client *http.Client, action, method, url string, data url.Values) {
	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("[%s] 创建请求失败: %v", action, err)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("[%s] 请求失败: %v", action, err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	status := "成功"
	if resp.StatusCode >= 400 {
		status = "失败"
	}

	fmt.Printf("[%s] %s | 状态码: %d\n", action, status, resp.StatusCode)
	fmt.Printf(">>> 响应内容: %s\n\n", string(result))
}
