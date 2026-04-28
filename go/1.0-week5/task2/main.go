package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/howeyc/gopass"
)

var (
	jar    *cookiejar.Jar
	client *http.Client
)

const (
	loginUrl string = "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
)

func erf(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {

	// const libUrl string = "http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx"
	var err error
	jar, err = cookiejar.New(nil)
	erf(err)
	client = &http.Client{
		Jar: jar,
	}
	getReq, err := http.NewRequest("GET", loginUrl, nil)
	erf(err)
	getReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	getReq.Header.Set("Accept", "text/html")
	getResp, err := client.Do(getReq)
	erf(err)
	defer getResp.Body.Close()
	getBody, err := io.ReadAll(getResp.Body)
	erf(err)
	// fmt.Println(string(getBody))
	html := string(getBody)
	doc := soup.HTMLParse(html)
	ltEle := doc.Find("input", "name", "lt")
	ltVal := ltEle.Attrs()["value"]
	executionEle := doc.Find("input", "name", "execution")
	executionVal := executionEle.Attrs()["value"]

	username, err := gopass.GetPasswdMasked()
	erf(err)
	password, err := gopass.GetPasswdMasked()
	erf(err)
	fromData := url.Values{}
	fromData.Add("username", string(username))
	fromData.Add("password", string(password))
	fromData.Add("lt", ltVal)
	fromData.Add("execution", executionVal)
	fromData.Add("_eventId", "submit")
	fromData.Add("submit", "登录")

	loginReq, err := http.NewRequest("POST", loginUrl, strings.NewReader(fromData.Encode()))
	erf(err)
	loginReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginReq.Header.Set("Referer", loginUrl)

	loginResp, err := client.Do(loginReq)
	erf(err)
	defer loginResp.Body.Close()
	//直接用我写好的模拟登录

	//cookiejar已经帮我们自动管理好cookie了
	//接下来只需要学类似用学号爬姓名的地址就好了

	//说在前面，我不想被学校图书馆拉黑，所以没测试
	//现在还没去过图书馆
	//但是已经有过舍友被神秘图书馆
	//取消座位制度坑过一次了
	//就当我这套逻辑可以吧
	//我觉得应该可以
	//逻辑应该都是一样的，朝着一个url发请求就可以了
	//哎哎，这种交互都是这样的

	//以南湖分馆一楼开敞座位区为例
	//如果是可选的也比较好做，直接去能工智人f12查一下就可以
	//这里以这个为例

	//座位号基本上都是Nxxxx
}
