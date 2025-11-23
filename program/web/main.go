package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	//还有个text/template 这个更适合处理非html的内容或处理已确认安全的数据
	//html版本具备强大的上下文感知能力和自动转义机制
	//能有效防止 XSS 攻击
)

type user struct {
	username string
	password string
	nickname string
	mutex    sync.RWMutex
}

var (
	users []user
)

func main() {
	users = make([]user, 100) //先来100个账户
	//go net/http 包已经给web做了非常好的支持
	//接下来学着搭一个web服务器
	http.Handle("./favicon.ico", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", sayhelloName) // 设置访问的路由
	http.HandleFunc("/pic", showPic)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//sayhelloName就是我们自己的实际功能实现的函数了
	//w是要发的响应，r是接收的请求
	//而且这个 Web 服务内部有支持高并发的特性

	//关于预防跨站脚本，也就是模板内html的转义
	//详情请见html/template包的main.go内容
	//阿里嘎多
}
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // 解析参数，默认是不会解析的
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme) //这个好像是说网页跳转app
	//还是说什么自定义url格式什么的
	//算了先不懂
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		//string.Join是字符串链接
	}
	fmt.Fprintf(w, "<h1>Hello astaxie!") // 这个写入到 w 的是输出到客户端的
	//效果的话就是显示Hello astaxie!
	//然后如果访问的url是http://localhost:9090/qwq/?url_long=111&url_long=222
	//终端输出是 out:
	/*
		map[url_long:[111 222]]
		path /qwq/
		scheme
		[111 222]
		key: url_long
		val: 111222
		// map[]
		// path /favicon.ico
		// scheme
		// []
	*/
	//这里后面还有一串map不知道是什么
	//哦，web服务会自动尝试请求/favicon.ico 图标文件
	//可以手动忽略（加个判断）
	// if r.URl.Path == "/favicon.ico" { return }
	//当然也可以自己手动加个图标
	//http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))
	//这一行可以避免请求出发sayhelloname
	//会直接读取文件
	//但好像不这么写也不会自动读图标啊
	//但是怎么感觉图标锁死了啊（（（（（
	//我自己换了一个又不能读了
	//哎就先这样吧
}
func showPic(w http.ResponseWriter, r *http.Request) {
	fileBytes, err := os.ReadFile("img.jpg")
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}
func login(w http.ResponseWriter, r *http.Request) {
	//写个简单的login的handler
	//然后也不考虑说存数据库啥的，先临时存一下就行
	//4. 不需要数据持久化，使用全局变量存储即可，也可以试试存文件里
	//要求是这么说的
	//账号密码用user结构体存
	//里面带了个RW互斥锁
	//然后既然是登录注册什么的
	//提交的一般是表单form
	//所以我们要先判断这个是什么方法传递过来 ，是POST还是GET
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		//这里是模板包
		log.Println(t.Execute(w, token))
	} else {
		err := r.ParseForm()
		// 解析 url 传递的参数
		// 对于 POST 则解析响应包的主体（request body）
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}
		//网页源代码这里token的value
		//比如是2550eaf8cfe9191235ca8806b6ed304a
		//这样就有标识了
		if err != nil {
			//handle error http.Error() for example
			log.Fatal("ParseForm:", err)
		}
		//请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Println("fruit:", r.Form["fruit"])
		//这里通过检查可以把fruit的选项value改掉
		//从而使得提交的表单是qwqawa
		//所以需要验证表单的输入
		fmt.Println("gender:", r.Form["gender"])
		fmt.Println("interest:", r.Form["interest"])
		//然后是防止多次提交表单
		//现在tempalte里面已经加了
		//<input type="hidden" name="token" value="{{.}}">
		//隐藏字段token，这个值我们通过 MD5 (时间戳) 来获取唯一值
		//然后我们把这个值存储到服务器端 (session 来控制，我们将在第六章讲解如何保存)
		//以方便表单提交时比对判定
	}
	//out :
	// method: GET
	// 2025/11/21 14:58:26 <nil>
	// method: POST
	// username: [FireChickenMP4]
	// password: [qwerty]
	//教程里说如果访问的是http://127.0.0.1:9090/login?username=astaxie
	//username输出是一个slice，存在多个值[qwq FireChickenMP4]

	//r.Form 里面包含了所有请求的参数
	//比如 URL 中 query-string、POST 的数据、PUT 的数据
	//所以当你在 URL 中的 query-string 字段和 POST 冲突时
	//会保存成一个 slice，里面存储了多个值
	//但是说Go之后版本会将POST和GET的数据分开
	//现在看来已经分开了
	//r.Form的类型是 url.Values
	//也就是map[string][]string类型
	//类似key=value吧，就是value可能存在多值，因为是切片嘛
	// r.Form.Set()	//设置成，应该是，但也还是只能设置两个
	//也就是覆盖原先key对应的值为新的value
	// r.Form.Add() //加
	//append实现
	// r.Form.Get("") //相当于找键值对 和直接r.Form[""]一样
	// r.Form.Encode() 这玩意弄出来就是?后面那一串，用&分隔的

	//Request本身也提供了FormValue()函数来获取用户提交的参数
	//r.Form["username"] 也可以写成r.FormValue("username")
	//后者也会自动调用r.ParseForm
	//但是FormValue只会返回同名参数的第一个 不存在则返回空字符串
}
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		//这里是模板包
		log.Println(t.Execute(w, token))
	} else {
		r.ParseMultipartForm(32 << 20)
		//最大内存32 <<20 ，应该就是32MB
		//如果文件大小超过了 maxMemory
		//那么剩下的部分将存储在系统的临时文件中
		//我们可以通过 r.FormFile 获取上面的文件句柄
		//然后实例中使用了 io.Copy 来存储文件。
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		//文件handler是multipart.FileHeader
		// 	type FileHeader struct {
		//     Filename string
		//     Header   textproto.MIMEHeader
		//     // contains filtered or unexported fields
		// }
		//页面返回信息
		//map[Content-Disposition:[form-data; name="uploadfile"; filename="0C64849BB077AF897EE0E1F842B5DEC9.jpg"] Content-Type:[image/jpeg]]
	}
}
