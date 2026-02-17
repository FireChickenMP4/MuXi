package main

import (
	"log"
	"net/http"
	"time"
	"week12/config"
	"week12/handle"
	"week12/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitRedis()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // 允许的请求源
			"http://localhost:5173", // 前端vite的默认启动地址
			"http://localhost:3000", // 前端自己定义的启动地址
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                   // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的请求头
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	routes.RegisterRoutes(r)
	go handle.RunHub()
	log.Println("服务器启动在 http://localhost:8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务器启动失败: %v\n", err)
	}
}
