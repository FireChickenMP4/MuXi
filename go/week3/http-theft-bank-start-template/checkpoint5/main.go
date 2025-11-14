package main

import (
	// "encoding/json"
	"bytes"
	"fmt"
	"strings"

	"github.com/Grand-Theft-Auto-In-CCNU-MUXI/hacker-support/encrypt"
	"github.com/Grand-Theft-Auto-In-CCNU-MUXI/hacker-support/httptool"
)

func erf(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	req1, err := httptool.NewRequest(
		httptool.GETMETHOD,
		"https://gtainmuxi.muxixyz.com/api/v1/organization/code",
		"",
		httptool.DEFAULT, // 这里可能不是 DEFAULT，自己去翻阅文档
	)
	erf(err)
	// fmt.Println(req)

	resp1, err := req1.SendRequest()
	erf(err)
	// resp1.ShowHeader()
	//Passport : eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
	// .
	// eyJjb2RlIjoiMSIsImlhdCI6MTc2MzAzODc2NSwibmJmIjoxNzYzMDM4NzY1fQ
	// .
	// HJnwPZJAguVRDODmCT9wMqYoMNHOTXIussjmUb5DP9c
	//应该是base64
	//JWT这是，JSON Web Token
	//{"alg":"HS256","typ":"JWT"}
	//{"code":"1","iat":1763038765,"nbf":1763038765}
	//HJnwPZJAguVRDODmCT9wMqYoMNHOTXIussjmUb5DP9c
	val, err := resp1.GetHeader("Passport")
	erf(err)
	// fmt.Println(val)
	req2, err := httptool.NewRequest(
		httptool.GETMETHOD,
		"https://gtainmuxi.muxixyz.com/api/v1/organization/code",
		"",
		httptool.DEFAULT, // 这里可能不是 DEFAULT，自己去翻阅文档
	)
	erf(err)
	req2.AddHeader("passport", val[0])
	resp2, err := req2.SendRequest()
	erf(err)
	// resp2.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		访问成功后，网站会给你返回信息，在header中找到你的passport。
		将passport加入到你以后的每次请求头中。
		完成上述步骤后，用代码访问 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/organization/secret_key，注意查收其中response的信息。
		3.ExtraInfo:
	*/
	// resp2.ShowHeader()
	//Passport : eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb2RlIjoiMSIsImlhdCI6MTc2MzA0MjU3MCwibmJmIjoxNzYzMDQyNTcwfQ.A1iou7SAXLuujC7Fu1CyqqQWfiipjdmiTAz86dwyJy0
	//令牌会变，所以要用之前结果
	val, err = resp2.GetHeader("Passport")
	erf(err)
	req3, err := httptool.NewRequest(
		httptool.GETMETHOD,
		"https://http-theft-bank.gtainccnu.muxixyz.com/api/v1/organization/secret_key",
		"",
		httptool.DEFAULT, // 这里可能不是 DEFAULT，自己去翻阅文档
	)
	erf(err)
	req3.AddHeader("passport", val[0])
	resp3, err := req3.SendRequest()
	erf(err)
	// resp3.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		恭喜你拿到了 passport，现在你可以着手准备骇入银行。
		银行的第一道门是代码安全门，我们计划将错误代码写入此门来破解它。
		但是这个门具有识别明文代码的功能，所以我们还需要一个密钥加密我们的错误代码，再写入至此门。
		不需要担心，两者我们都为你提供了，你只需要解析 response 中的密文（在 ExtraInfo 中）来得到它们。
		你现在的任务：
		解析密文，获取 error_code 和 secret_key
		使用 secret_key 加密 error_code
		然后将加密过的 error_code 放入 请求body 并以 "正确的请求方法" 发送至 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/gate , 同时注意response的 信息。
		3.ExtraInfo:
		c2VjcmV0X2tleTpNdXhpU3R1ZGlvMjAzMzA0LCBlcnJvcl9jb2RlOmZvciB7Z28gZnVuYygpe3RpbWUuU2xlZXAoMSp0aW1lLkhvdXIpfSgpfQ==
	*/
	//示例base64解码 secret_key:MuxiStudio203304, error_code:for {go func(){time.Sleep(1*time.Hour)}()}
	//貌似其实不会变（
	//通过开盒ShowBody，或者说分析结构体，HttpResponse层层追
	/*
		func (r *HttpResponse) ShowBody() {
		if len(r.Raw) != 0 {
			fmt.Println("Body:")
			fmt.Println(r.Raw)
			return
		}

		fmt.Println("Message:")
		fmt.Println("response body:")
		fmt.Println("1.Message:")
		fmt.Println(r.Body.Message)
		fmt.Println("2.Text:")
		fmt.Println(r.Body.Data.Text)
		fmt.Println("3.ExtraInfo:")
		fmt.Println(r.Body.Data.ExtraInfo)
		}
	*/
	//r.Body.Data.ExtraInfo是我们想要的
	extraInfo := resp3.Body.Data.ExtraInfo
	psw, err := encrypt.Base64Decode(extraInfo)
	erf(err)
	// fmt.Println(psw)
	//secret_key:MuxiStudio203304, error_code:for {go func(){time.Sleep(1*time.Hour)}()}
	//现在我要把secret_key和error_code切开
	tempStr := strings.Split(psw, ":")
	// for i := range tempStr {
	// 	fmt.Println(tempStr[i])
	// }
	/*
		secret_key
		MuxiStudio203304, error_code
		for {go func(){time.Sleep(1*time.Hour)}()}
	*/
	secret_key := strings.Split(tempStr[1], ",")[0]
	error_code := tempStr[2]
	data, err := encrypt.AESEncryptOutInBase64([]byte(error_code), []byte(secret_key))
	erf(err)
	/*
		req4, err := httptool.NewRequest(
			httptool.POSTMETHOD,
			"http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/gate",
			string(data),
			httptool.DEFAULT, // 这里可能不是 DEFAULT，自己去翻阅文档
		)
		erf(err)
		req4.AddHeader("passport", val[0])
		resp4, err := req4.SendRequest()
		erf(err)
		resp4.ShowBody()
		resp4.ShowHeader()
	*/
	//这里。。。emm
	//你因为使用 错误的方法 写入病毒，被银行安全系统识别。很遗憾，你被逮捕了。（请尝试使用其他方法重新访问)
	//被气笑了
	//但是想来其实既然是code的话，应该提交文件，用post传上去注入才对（，经典注入
	//文件怎么办呢。。。
	//用os包创建一个好了
	// file, err := os.Create("Data.txt")
	// erf(err)
	// defer file.Close()
	// file.WriteString(string(data))
	req4, err := httptool.NewRequest(
		httptool.PUTMETHOD,
		"http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/gate",
		string(data),
		httptool.DEFAULT,
	)
	erf(err)
	req4.AddHeader("passport", val[0])
	resp4, err := req4.SendRequest()
	erf(err)
	// resp4.ShowBody()
	//还是不行
	//但是我们不妨思考，POST发新建可能会被警惕，那PUT呢？
	//修改代码的话应该可以的吧
	//直接PUT字符串可以，POST不再张贴
	/*
		1.Message:
		OK
		2.Text:
		干的漂亮！你已经突破了第一扇门，请继续访问 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/iris_recognition_gate 。
		3.ExtraInfo:
	*/
	// req5, err := httptool.NewRequest(
	// 	httptool.GETMETHOD,
	// 	"http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/iris_recognition_gate",
	// 	"",
	// 	httptool.DEFAULT,
	// )
	// erf(err)
	// req5.AddHeader("passport", val[0])
	// resp5, err := req5.SendRequest()
	// erf(err)
	// resp5.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		你现在已经到第二扇门了，是虹膜识别安全门。
		你需要向组织请求已准备好的虹膜样本，访问 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/organization/iris_sample 下载图片。
		再将此图片上传至 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/iris_recognition_gate 以破解此门，加油！
		3.ExtraInfo:
	*/
	req6, err := httptool.NewRequest(
		httptool.GETMETHOD,
		"http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/organization/iris_sample",
		"",
		httptool.DEFAULT,
	)
	erf(err)
	req6.AddHeader("passport", val[0])
	resp6, err := req6.SendRequest()
	erf(err)
	resp6.Save("temp.jpg")
	req7, err := httptool.NewRequest(
		httptool.POSTMETHOD,
		"http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/bank/iris_recognition_gate",
		"temp.jpg",
		httptool.FILE,
	)
	erf(err)
	req7.AddHeader("passport", val[0])
	resp7, err := req7.SendRequest()
	erf(err)
	// resp7.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		还剩最后一道门了！
		我们需要银行结构图碎片，这些碎片就隐藏在前面某四个路由的响应头中，位于 map-fragments 字段。
		将它们用"/"拼起来就是最后一道门的所在位置！注意response的信息。
		3.ExtraInfo:
	*/

	tempStr1 := []string{}
	tempM, err := resp2.GetHeader("Map-Fragments")
	erf(err)
	tempStr1 = append(tempStr1, tempM...)
	tempM, err = resp3.GetHeader("Map-Fragments")
	erf(err)
	tempStr1 = append(tempStr1, tempM...)
	tempM, err = resp4.GetHeader("Map-Fragments")
	erf(err)
	tempStr1 = append(tempStr1, tempM...)
	tempM, err = resp7.GetHeader("Map-Fragments")
	erf(err)
	tempStr1 = append(tempStr1, tempM...)
	// fmt.Println(tempStr1)
	//[muxi backend computer examination]
	var buf bytes.Buffer
	buf.WriteString("http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/")
	for _, val := range tempStr1 {
		buf.WriteString(val)
		buf.WriteString("/")
	}
	// fmt.Println(buf.String())
	// http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/backend/computer/examination/
	req8, err := httptool.NewRequest(
		httptool.GETMETHOD,
		buf.String(),
		"",
		httptool.DEFAULT,
	)
	erf(err)
	req8.AddHeader("passport", val[0])
	resp8, err := req8.SendRequest()
	erf(err)
	resp8.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		1，真亏你能来到这里！在你眼前的就是最后的密码门了。
		但是密码位数未知，看来只能通过全排列程序去暴力破解。

		>示例如下：
		============================================
		输入：
		3
		1 2 3
		输出：
		[[1 2 3][1 3 2][2 1 3][2 3 1][3 1 2][3 2 1]]
		============================================

		>代码模板:

		func permute(nums []int) [][]int {
		    // insert your code

		}

		func main() {
		    var n int
		        fmt.Scanf("%d", &n)

		        testSlice := make([]int, n)
		    // 标准输入n个不重复的数字

		    res := permute(testSlice)
		    fmt.Println(res)
		}

		请在完成此程序后上传至 http://http-theft-bank.gtainccnu.muxixyz.com/api/v1/muxi/backend/computer/examination 来破解最后的密码门
		3.ExtraInfo:
	*/
	//我也不知道是完整程序结构还是单独就这俩函数
	//我也不知道是用POST还是PUT甚至PATCH
	//算了还是真的自己写到一个文件里吧
	req9, err := httptool.NewRequest(
		httptool.POSTMETHOD,
		buf.String(),
		"temp/main.go",
		httptool.FILE,
	)
	erf(err)
	req9.AddHeader("passport", val[0])
	resp9, err := req9.SendRequest()
	erf(err)
	resp9.ShowBody()
	/*
		1.Message:
		OK
		2.Text:
		END
		我就知道你能成功！Backend组织欢迎你！
		3.ExtraInfo:
	*/
}
