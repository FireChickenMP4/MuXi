package main

import (
	"gin/code"
)

func main() {
	//今儿个弄gin玩玩
	// r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run()
	//这里gin有一些warning
	//其中一个是有关模式的问题
	//debug会显示所有细节
	//而release只显示重要事件
	//一个是防止日志被利用，一个是节省性能
	//gin.SetMode(gin.ReleaseMode) // 设置为发布模式
	//另外两个，一个是端口默认8080.读取系统变量PORT未成
	//另一个是可信代理列表
	//r.SetTrustedProxies([]string{"192.168.1.1", "10.0.0.0/8", "172.16.0.0/12"})
	//这样风险会小一些
	//无所谓，自己写小代码，这些都可以不管

	// [GIN] 2025/12/15 - 22:15:31 | 404 |            0s |       127.0.0.1 | GET      "/"
	// [GIN] 2025/12/15 - 22:15:31 | 404 |            0s |       127.0.0.1 | GET      "/favicon.ico"
	// [GIN] 2025/12/15 - 22:15:37 | 200 |            0s |       127.0.0.1 | GET      "/ping"
	//可以看到GIN的这个非常，emm，高级啊
	//啥都告诉你了
	//然后咱个也是收到json message:pong 了啊

	//接下来看官方完整一些的实例喵
	//在code.go里面
	r := code.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	_ = r.Run(":8080")
}
