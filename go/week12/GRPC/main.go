package main

import (
	"context"
	"grpc/helloworld"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go rpcserver()
	client()
}

func client() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败：%v", err)
	}
	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		name := c.DefaultQuery("name", "World")
		req := &helloworld.HelloRequest{Name: name}

		resp, err := client.SayHello(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": resp.Message})
	})
	r.Run(":8080")
}

func rpcserver() {
	//这里直接服务端和客户端放一块了
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("服务器监听失败 %v", err)
	}
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve：%v", err)
	}
}

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello" + req.Name}, nil
}
