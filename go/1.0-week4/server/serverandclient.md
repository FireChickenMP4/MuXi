以下是根据要求优化后的完整服务端代码和配套客户端测试代码。优化重点包括：**修复核心逻辑漏洞**、**增强安全性**、**完善 HTTP 状态码**、**增加用户信息字段**，同时保持代码简洁易懂。

### 优化后的服务端代码 (`server.go`)

```go
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
				Name:    "session_user",
				Value:   username,
				Path:    "/",
				Expires: expiration,
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
				"message": "用户信息更新成功",
				"nickname": users[i].Nickname,
				"email":   users[i].Email,
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
```

### 客户端测试代码 (`client_test.go`)

```go
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
```

### 关键优化说明（小白友好版）：

#### 1. 用户数据结构

```go
type User struct {
	Username string `json:"username"`
	Password string `json:"-"` // 密码永不返回给客户端
	Nickname string `json:"nickname"` // 新增昵称
	Email    string `json:"email,omitempty"` // 可选邮箱
}
```

- **类比**：就像身份证包含姓名、昵称（小名）、邮箱（联系方式），但银行卡密码（Password）永远不会印在身份证上
- **错误做法**：原代码只存储用户名密码，没有基本信息

#### 2. 安全增强

```go
cookie := http.Cookie{
	...
	HttpOnly: true, // 阻止JavaScript访问Cookie
	Path: "/",      // 所有页面都能用
}
```

- **为什么重要**：防止黑客通过网页脚本窃取你的登录凭证
- **错误做法**：原代码没有设置 HttpOnly，存在安全风险

#### 3. HTTP 状态码规范

```go
// 正确使用状态码示例
http.Error(w, "用户名已被占用", http.StatusConflict) // 409
w.WriteHeader(http.StatusCreated) // 201 创建成功
```

- **200** = 成功 | **400** = 你填错了 | **401** = 你没登录 | **409** = 名字被占用了
- **错误做法**：原代码所有错误都返回 200，客户端无法区分成功失败

#### 4. 会话管理修复

```go
// 获取当前用户
cookie, err := r.Cookie("session_user") // 统一Cookie名称
```

- **关键修复**：使用固定 Cookie 名称`session_user`，避免路径问题
- **会话隔离**：修改用户名后，旧会话自动失效（需重新登录）

#### 5. 用户信息管理

```go
// 新增/profile接口
func updateProfile(w http.ResponseWriter, r *http.Request) {
	// 从Cookie获取当前用户，无需重复输入密码
}
```

- **设计逻辑**：登录后修改基本信息（昵称/邮箱）不需要反复输密码
- **错误做法**：原代码没有修改用户信息的接口

### 使用指南：

1. **启动服务端**：

   ```bash
   go run server.go
   ```

2. **运行测试客户端**：

   ```bash
   go run client_test.go
   ```

3. **预期成功输出**：

   ```
   [注册] 成功 | 状态码: 201
   >>> 响应内容: 注册成功！用户名: testuser

   [登录] 成功 | 状态码: 200
   >>> 响应内容: 登录成功！欢迎回来, 测试用户

   [获取用户信息] 成功 | 状态码: 200
   >>> 响应内容: {"username":"testuser","nickname":"测试用户","email":"test@example.com"}

   ...（其他步骤省略）...

   [最终验证] 成功 | 状态码: 200
   >>> 响应内容: {"username":"freshuser","nickname":"新昵称","email":"new@example.com"}
   ```

### 常见问题处理：

1. **如果看到 401 错误**：

   - 检查是否先登录再访问/whoami
   - 修改用户名后必须重新登录

2. **如果注册失败**：

   - 用户名可能已被占用（尝试换名字）
   - 密码少于 6 位字符

3. **测试后清理**：
   - 服务端使用内存存储，重启服务即重置
   - 如需持久化，可将`users`切片保存到 JSON 文件（扩展功能）

> 💡 **学习重点**：这个实现完整展示了**用户生命周期管理**（注册 → 登录 → 修改信息 → 验证），同时通过规范的 HTTP 状态码和安全措施，达到了生产级基础要求。实际项目中只需增加密码加密（bcrypt）和数据库存储即可。
