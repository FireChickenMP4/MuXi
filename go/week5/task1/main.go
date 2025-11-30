package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/howeyc/gopass"
)

func main() {
	erf := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	const loginUrl string = "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	// const libUrl string = "http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx"
	jar, err := cookiejar.New(nil)
	erf(err)
	client := &http.Client{
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

	// loginBody, err := io.ReadAll(loginResp.Body)
	// erf(err)
	// fmt.Println(string(loginBody))
	file, err := os.Create("loginRespBody.html")
	erf(err)
	defer file.Close()
	io.Copy(file, loginResp.Body)
	//http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=2025211351&_=1764435994060
	//我去搜某个研讨间的成员的时候
	//以我为例子，是这样的
	//type ReservaApply看来不变
	//term是你要查的
	//然后有个-一样的东西，应该是什么编码
	//unix毫秒数应该
	const reserveUrl string = "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx"
	queryForm := url.Values{}
	queryForm.Add("type", "logonname")
	queryForm.Add("ReservaApply", "ReservaApply")
	queryForm.Add("term", "2025211351")
	now := time.Now()
	queryForm.Add("_", string(now.UnixMilli()))
	queryReq, err := http.NewRequest("GET", reserveUrl, strings.NewReader(queryForm.Encode()))
	erf(err)
	queryResp, err := client.Do(queryReq)
	erf(err)
	queryBody, err := io.ReadAll(queryResp.Body)
	erf(err)
	fmt.Println(string(queryBody))
}
