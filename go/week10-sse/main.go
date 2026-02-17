package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/events", sseHandler)
	r.GET("/", index)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	//然后效果就是，对于events路由的访问，就是一条条推送
	//根路由下前端通过EventSource挂载了api之后就会自动解析数据流
	//再格式化，就实现了动态时间界面
}

func index(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", []byte(`
            <!DOCTYPE html>
            <html>
            <head><title>SSE 测试</title></head>
            <body>
                <div id="time"></div>
                <script>
                    const source = new EventSource('/events');
                    source.onmessage = function(event) {
                        document.getElementById('time').innerHTML = event.data;
                    };
                </script>
            </body>
            </html>
        `))
	//ds老师写的
}

func sseHandler(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	//这个是SSE标准的MIME类型
	c.Header("Connection", "keep-alive")
	c.Header("Cache-Control", "no-cache")
	//sse需要的请求头
	//这里c.Writer.Header().Set()也行

	clientGone := c.Done() //这是断开连接的chan

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-clientGone:
			log.Println("Client disconnected")
			return
		case t := <-ticker.C:
			//sse是以文本块形式发送的，也就是结尾是\n\n
			//以及格式要求开头是data:
			_, err := fmt.Fprintf(c.Writer, "data: The time is %s\n\n", t.Format(time.UnixDate))
			//这里调用Writer的Write方法也可以，但是也是要格式化msg的，然后发[]byte
			if err != nil {
				return
			}
			c.Writer.Flush()
			//这个是把缓冲区的数据理解发送给客户端
			//SSE实时推送要有这个
		}
	}
}
