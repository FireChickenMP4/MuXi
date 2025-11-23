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
	// 创建带自动Cookie管理的客户端（关键！）
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar, // 启用标准Cookie管理
		Timeout: 10 * time.Second,
	}

	// 1. 注册
	username := "testuser"
	password := "testpass"
	registerData := url.Values{"username": {username}, "password": {password}}
	resp := sendFormRequest(client, "POST", "http://localhost:9090/register", registerData)
	printResponse("注册", resp)
	resp.Body.Close()

	// 2. 初始登录
	loginData := url.Values{"username": {username}, "password": {password}}
	resp = sendFormRequest(client, "POST", "http://localhost:9090/login", loginData)
	printResponse("初始登录", resp)
	resp.Body.Close()

	// 3. 修改密码
	newPassword := "newpass123"
	passData := url.Values{
		"username":    {username},
		"password":    {password},
		"newpassword": {newPassword},
	}
	resp = sendFormRequest(client, "POST", "http://localhost:9090/password", passData)
	printResponse("修改密码", resp)
	resp.Body.Close()

	// 4. 修改用户名
	newUsername := "newuser"
	userData := url.Values{
		"username":    {username},
		"password":    {newPassword}, // 必须用新密码！
		"newusername": {newUsername},
	}
	resp = sendFormRequest(client, "POST", "http://localhost:9090/username", userData)
	printResponse("修改用户名", resp)
	resp.Body.Close()

	// 5. 重要！退出当前会话（清除旧Cookie）
	client.Jar = nil // 清除旧Cookie
	jar, _ = cookiejar.New(nil)
	client.Jar = jar // 创建新会话

	// 6. 用新凭证重新登录
	newLoginData := url.Values{
		"username": {newUsername},
		"password": {newPassword},
	}
	resp = sendFormRequest(client, "POST", "http://localhost:9090/login", newLoginData)
	printResponse("新账号登录", resp)
	resp.Body.Close()

	// 7. 验证当前用户（自动携带新Cookie）
	resp = sendFormRequest(client, "GET", "http://localhost:9090/whoami", nil)
	printResponse("当前用户", resp)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf(">>> 服务器返回: '%s'\n", string(body)) // 应该返回 newuser
	resp.Body.Close()
	//out :
	// [注册] 服务器响应: 状态码=200
	// [初始登录] 服务器响应: 状态码=200
	// [修改密码] 服务器响应: 状态码=200
	// [修改用户名] 服务器响应: 状态码=200
	// [新账号登录] 服务器响应: 状态码=200
	// [当前用户] 服务器响应: 状态码=200
	// >>> 服务器返回: 'newuser'
}

func sendFormRequest(client *http.Client, method, url string, data url.Values) *http.Response {
	var body io.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal("创建请求失败:", err)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	return resp
}

func printResponse(action string, resp *http.Response) {
	fmt.Printf("[%s] 服务器响应: 状态码=%d\n", action, resp.StatusCode)
}
