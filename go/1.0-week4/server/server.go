package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// 用户结构体（增加昵称等基本信息）
type User struct {
	Username string `json:"username"`
	Password string `json:"-"` // 不返回给客户端
	Nickname string `json:"nickname"`
	Email    string `json:"email,omitempty"` // 可选字段
}

var (
	users = []User{} // 初始化为空切片
)

func main() {
	// 初始化测试用户（方便测试，实际应移除）
	users = append(users, User{
		Username: "admin",
		Password: "secure123",
		Nickname: "管理员",
		Email:    "admin@example.com",
	})

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/password", updatePassword)
	http.HandleFunc("/username", updateUsername)
	http.HandleFunc("/profile", updateProfile) // 新增：修改用户信息
	http.HandleFunc("/whoami", getUser)

	fmt.Println("服务器启动: http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

// 注册新用户
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	password := strings.TrimSpace(r.Form.Get("password"))
	nickname := strings.TrimSpace(r.Form.Get("nickname"))
	email := strings.TrimSpace(r.Form.Get("email"))

	// 基础校验
	if username == "" || password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}
	if len(password) < 6 {
		http.Error(w, "密码至少需要6位字符", http.StatusBadRequest)
		return
	}

	// 检查用户名是否存在
	for _, u := range users {
		if u.Username == username {
			http.Error(w, "用户名已被占用", http.StatusConflict)
			return
		}
	}

	// 创建新用户（默认昵称=用户名）
	if nickname == "" {
		nickname = username
	}

	newUser := User{
		Username: username,
		Password: password, // 实际应用应加密存储
		Nickname: nickname,
		Email:    email,
	}
	users = append(users, newUser)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "注册成功！用户名: %s", username)
}

// 用户登录
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	password := strings.TrimSpace(r.Form.Get("password"))

	// 验证凭证
	var found bool
	for _, u := range users {
		if u.Username == username && u.Password == password {
			found = true

			// 设置登录Cookie（有效期1天）
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{
				Name:     "session_user",
				Value:    username,
				Path:     "/",
				Expires:  expiration,
				HttpOnly: true, // 增强安全性
			}
			http.SetCookie(w, &cookie)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "登录成功！欢迎回来, %s", u.Nickname)
			return
		}
	}

	if !found {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
	}
}

// 修改密码
func updatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	oldPass := strings.TrimSpace(r.Form.Get("password"))
	newPass := strings.TrimSpace(r.Form.Get("newpassword"))

	// 基础校验
	if newPass == "" || len(newPass) < 6 {
		http.Error(w, "新密码至少需要6位字符", http.StatusBadRequest)
		return
	}

	// 验证旧密码
	for i, u := range users {
		if u.Username == username && u.Password == oldPass {
			// 更新密码
			users[i].Password = newPass
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "密码修改成功")
			return
		}
	}

	http.Error(w, "用户名或旧密码错误", http.StatusUnauthorized)
}

// 修改用户名
func updateUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效表单数据", http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	password := strings.TrimSpace(r.Form.Get("password"))
	newUsername := strings.TrimSpace(r.Form.Get("newusername"))

	// 检查新用户名
	if newUsername == "" {
		http.Error(w, "新用户名不能为空", http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.Username == newUsername {
			http.Error(w, "新用户名已被占用", http.StatusConflict)
			return
		}
	}

	// 验证当前凭证
	for i, u := range users {
		if u.Username == username && u.Password == password {
			// 更新用户名
			users[i].Username = newUsername
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "用户名修改成功! 新用户名: %s", newUsername)
			return
		}
	}

	http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
}

// 修改用户信息（昵称/邮箱）
func updateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 从Cookie获取当前登录用户
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "无效表单数据", http.StatusBadRequest)
		return
	}

	nickname := strings.TrimSpace(r.Form.Get("nickname"))
	email := strings.TrimSpace(r.Form.Get("email"))

	// 更新用户信息
	for i, u := range users {
		if u.Username == cookie.Value {
			if nickname != "" {
				users[i].Nickname = nickname
			}
			if email != "" {
				users[i].Email = email
			}

			response := map[string]string{
				"message":  "用户信息更新成功",
				"nickname": users[i].Nickname,
				"email":    users[i].Email,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "用户不存在", http.StatusNotFound)
}

// 获取当前用户信息
func getUser(w http.ResponseWriter, r *http.Request) {
	// 从Cookie获取用户名
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}

	// 查找用户
	for _, u := range users {
		if u.Username == cookie.Value {
			// 创建安全响应（不返回密码）
			safeUser := struct {
				Username string `json:"username"`
				Nickname string `json:"nickname"`
				Email    string `json:"email,omitempty"`
			}{
				Username: u.Username,
				Nickname: u.Nickname,
				Email:    u.Email,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(safeUser)
			return
		}
	}

	http.Error(w, "用户不存在", http.StatusNotFound)
}
