package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/howeyc/gopass"
)

type User struct {
	StudentId string `json:"id"`
	Name      string `json:"name"`
	Grade     string `json:"grade,omitempty"`
	//omitempty属性是如果没加，就不显示
}
type UserS []User

var (
	jar    *cookiejar.Jar
	client *http.Client
	wg     *sync.WaitGroup
	ch     chan User
	signal chan struct{}
	limit  int
	first  int
)

const (
	loginUrl   string = "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page="
	reserveUrl string = "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx"
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

	// loginBody, err := io.ReadAll(loginResp.Body)
	// erf(err)
	// fmt.Println(string(loginBody))
	// file, err := os.Create("loginRespBody.html")
	// erf(err)
	// defer file.Close()
	// io.Copy(file, loginResp.Body)
	//http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=2025211351&_=1764435994060
	//我去搜某个研讨间的成员的时候
	//以我为例子，是这样的
	//type ReservaApply看来不变
	//term是你要查的
	//然后有个-一样的东西，应该是什么编码
	//unix毫秒数应该
	//但是好像也不用带（
	//不带数据也能返回
	//服务器好像也不在乎这一小会的差距？
	//就算传输耗了点时间

	// queryForm := url.Values{}
	// queryForm.Add("type", "logonname")
	// queryForm.Add("ReservaApply", "ReservaApply")
	// queryForm.Add("term", "2025211351")
	// now := time.Now()
	// queryForm.Add("_", fmt.Sprint(now.UnixMilli()))
	// fullUrl := reserveUrl + "?" + queryForm.Encode()
	// queryReq, err := http.NewRequest("GET", fullUrl, nil)
	// //这里不要带请求体
	// erf(err)
	// queryResp, err := client.Do(queryReq)
	// erf(err)
	// queryBody, err := io.ReadAll(queryResp.Body)
	// erf(err)
	// fmt.Println(string(queryBody))

	wg = &sync.WaitGroup{}
	limit = 25 //大于21
	first = 211000
	wg.Add(limit)
	ch = make(chan User, limit*200+5)
	signal = make(chan struct{})
	go func() {
		ok := true
		for i := range 20 {
			go queryName(i)
		}
		i := 20
		for ok {
			_, ok = <-signal
			go queryName(i)
			i++
		}
	}()
	var users []User
	go func() {
		// encoder := json.NewEncoder(file)
		for user := range ch {
			// err = encoder.Encode(val)
			// erf(err)
			// 为了好排序，用个users切片存一下
			users = append(users, user)
			if len(users)%100 == 0 {
				fmt.Printf("已收集 %d 条记录...\n", len(users))
			}
		}
	}()
	wg.Wait()
	close(signal)
	close(ch)
	//学了个encoder的写法
	//用Fprintln也可以
	// 可以直接用写入函数
	// io.Copy(file, bytes.NewReader([]byte("你好")))  // 复制数据到文件
	// fmt.Fprintln(file, "这是文本")                 // 格式化写入
	// json.NewEncoder(file).Encode(user)            // JSON编码到文件

	// // 读取时同样适用
	// data, _ := io.ReadAll(file)                    // 读取所有内容
	// scanner := bufio.NewScanner(file)              // 逐行扫描文件
	//回头用数据库搞一下

	//这里因为前面太多不是的了，所以
	//直接从211开始爬了

	//这里说一下
	//我本来打算去想办法处理json数组的问题
	//但是发现，ndjson或者说jsonl
	//其实对于这类流数据更适合
	//感觉也更适合存这种数据
	//所以就换用ndjson了
	//ndjson每一行都是一个独立的json对象
	//每添加一个对象只需要添加对应的一行
	//所以适合日志记录、大规模数据处理和实时数据流等用途
	sort.Sort(UserS(users))
	//这里我们发现需要传递一个data Interface类型的参数
	//然后Interface需要实现这三个参数
	//Len() int
	//Less(i,j int) bool
	//Swap(i,j int)
	file, err := os.OpenFile("Student.ndjson", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//os.O_APPEND 追加模式
	//os.O_CREATE 创建模式
	//os.O_WRONLY 只写模式
	erf(err)
	defer file.Close()
	encoder := json.NewEncoder(file)
	for _, user := range users {
		err = encoder.Encode(user)
		erf(err)
	}

	//然后对于长期来说，还是说最后转成json比较好
	//新建文件夹......
}

func queryName(start int) {
	startid := start*200 + first
	if start >= limit {
		return
	}
	for i := range 200 {
		queryForm := url.Values{}
		queryForm.Add("type", "logonname")
		queryForm.Add("ReservaApply", "ReservaApply")
		queryForm.Add("term", "2025"+fmt.Sprintf("%06d", startid+i))
		now := time.Now()
		queryForm.Add("_", fmt.Sprint(now.UnixMilli()))
		fullUrl := reserveUrl + "?" + queryForm.Encode()
		queryReq, err := http.NewRequest("GET", fullUrl, nil)
		erf(err)
		queryReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		queryReq.Header.Set("Referer", "http://kjyy.ccnu.edu.cn/clientweb/xcus/ic2/Default.aspx")
		queryReq.Header.Set("X-Requested-With", "XMLHttpRequest")
		queryResp, err := client.Do(queryReq)
		erf(err)
		defer queryResp.Body.Close()
		queryBody, err := io.ReadAll(queryResp.Body)
		erf(err)
		var users []User
		err = json.Unmarshal(queryBody, &users)
		erf(err)
		if len(users) != 0 {
			user := users[0]
			user.Grade = "2025"
			user.StudentId = "2025" + fmt.Sprintf("%06d", startid+i)
			// fmt.Println(user)
			ch <- user
		}
	}
	signal <- struct{}{}
	wg.Done()
}
func (users UserS) Len() int {
	return len(users)
}
func (users UserS) Swap(i, j int) {
	users[i], users[j] = users[j], users[i]
}
func (users UserS) Less(i, j int) bool {
	return users[i].StudentId < users[j].StudentId
}
