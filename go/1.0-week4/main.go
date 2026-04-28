package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type user struct {
	username string
	password string
}

var (
	users []user
)

func main() {
	//要写作业力
	//顺便说学一下怎么存储用户状态好了
	//然后前端页面也就懒得用template写了
	//光写个api，然后用client模拟请求试一下就得了
	users = make([]user, 100)

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/password", password)
	http.HandleFunc("/username", username)
	http.HandleFunc("/whoami", getUser)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//handle error http.Error() for example
		log.Fatal("ParseForm:", err)
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	for _, user := range users {
		if user.username == username {
			fmt.Println("用户已存在")
			return
		}
	}
	if r.Form.Get("username") != "" && r.Form.Get("password") != "" {
		users = append(users, user{
			username,
			password,
		})
	}
	fmt.Println("注册成功，用户名：", r.Form.Get("username"), "密码：", r.Form.Get("password"))
}
func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//handle error http.Error() for example
		log.Fatal("ParseForm:", err)
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	for _, user := range users {
		if user.username == username && user.password == password {
			fmt.Println("登录成功，欢迎回来：", username)
			//cookie通过http包的SetCookie来设置
			// type Cookie struct {
			//     Name       string
			//     Value      string
			//     Path       string
			//     Domain     string
			//     Expires    time.Time
			//     RawExpires string

			// // MaxAge=0 means no 'Max-Age' attribute specified.
			// // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
			// // MaxAge>0 means Max-Age attribute present and given in seconds
			//     MaxAge   int
			//     Secure   bool
			//     HttpOnly bool
			//     Raw      string
			//     Unparsed []string // Raw text of unparsed attribute-value pairs
			// }

			expiration := time.Now()
			expiration = expiration.AddDate(1, 0, 0)
			cookie := http.Cookie{Name: "username", Value: username, Expires: expiration, Path: "/"}
			//这里path要设置成根目录
			//否则只会在/login生效
			http.SetCookie(w, &cookie)

			return
		}
	}
	fmt.Println("登录失败，用户名或密码错误")
}
func password(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//handle error http.Error() for example
		log.Fatal("ParseForm:", err)
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	newpassword := r.Form.Get("newpassword")
	canLogin := false
	foundedIndex := -1
	for index, user := range users {
		if user.username == username && user.password == password {
			canLogin = true
			foundedIndex = index
		}
	}
	if canLogin {
		users[foundedIndex].password = newpassword
	} else {
		fmt.Println("用户名或密码错误")
	}
}
func username(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//handle error http.Error() for example
		log.Fatal("ParseForm:", err)
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	newusername := r.Form.Get("newusername")
	canLogin := false
	foundedIndex := -1
	for index, user := range users {
		if user.username == username && user.password == password {
			canLogin = true
			foundedIndex = index
		}
	}
	if canLogin {
		users[foundedIndex].username = newusername
		fmt.Println("旧密码：", password, "新密码：", newusername)
	} else {
		fmt.Println("用户名或密码错误")
	}
}
func getUser(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	fmt.Fprint(w, cookie.Value)
	// or:
	// for _, cookie := range r.Cookies() {
	// 	fmt.Fprint(w, cookie.Name)
	// }
}
