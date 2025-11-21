package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	//go net/http 包已经给web做了非常好的支持
	//接下来学着搭一个web服务器
	http.Handle("./favicon.ico", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", sayhelloName) // 设置访问的路由
	http.HandleFunc("/pic", showPic)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//sayhelloName就是我们自己的实际功能实现的函数了
	//w是要发的响应，r是接收的请求
	//而且这个 Web 服务内部有支持高并发的特性
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
