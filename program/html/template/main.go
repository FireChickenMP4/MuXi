package main

import (
	"fmt"
	"html/template"
	"net/url"
)

func main() {
	//这里是从web转过来的一个实例演示
	//预防跨站脚本
	//虽然我感觉它叽里咕噜的说的不清楚
	//但是我大概知道它可能想表达类似SQL注入一样的东西
	//所以html/template包里面有html转义
	//可以把html有关的转义掉，防止其在浏览器端生效
	v := url.Values{}
	v.Add("username", "<script>alert()</script>")
	v.Add("password", "qwerty")
	fmt.Println("username:", template.HTMLEscapeString(v.Get("username")))
	fmt.Println("password:", template.HTMLEscapeString(v.Get("password")))
	// template.HTMLEscape(w, []byte(v.Get("username")))
	//输出到客户端这是
	//out :
	// username: &lt;script&gt;alert()&lt;/script&gt;
	// password: qwerty
	//有点像strings包那个Replacer的感觉
	//反正换掉了是，也是html的转义
	//Go 的 html/template 包默认帮你过滤了 html 标签
	//但是有时候你只想要输出这个 <script>alert()</script> 看起来正常的信息
	//该怎么处理？请使用text/template
	// t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	// err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")
	//out :
	//Hello, <script>alert('you have been pwned')</script>!
	//或者使用template.HTML类型
	// t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	// err = t.ExecuteTemplate(w, "T", template.HTML("<script>alert('you have been pwned')</script>"))
	//转换成 template.HTML 后，变量的内容也不会被转义
	//转义的例子：
	// t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	// err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")
	//out :
	//Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
}
