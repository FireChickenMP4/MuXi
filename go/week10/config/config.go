package config

import (
	"context"
	"log"
	"net"
	"time"
	"week12/utils"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	var password string
	password = utils.GetEnv("REDIS_PASSWORD", "123456qwerty")
	opt := &redis.Options{
		Addr:     "localhost:6379",
		Password: password,
		DB:       0,
		// TLSConfig: &tls.Config{
		// 	MinVersion: tls.VersionTLS12,
		// },
		//要启用 TLS/SSL，你需要提供一个空的 tls.Config
		//如果你使用的是私有证书，你需要 指定 它们在 tls.Config 中
	}
	Rdb = redis.NewClient(opt)
	//或者用连续字符串
	// opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	// if err != nil {
	// 	panic(err)
	// }
	// rdb := redis.NewClient(opt)
}
func InitRedisBySSH() {
	//或者通过SSH通道连接
	var err error
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("")},
		// HostKeyCallBack: ssh.InsecureIgnoreHostKey(),
		// 好像不存在这个字段，但是go redis里面有
		Timeout: 15 * time.Second,
	}
	sshClient, err := ssh.Dial("tcp", "remoteIP:22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	opt := &redis.Options{
		Addr: net.JoinHostPort("127.0.0.1", "6379"),
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		},
		// Disable timeouts, because SSH does not support deadlines.
		ReadTimeout:  -1,
		WriteTimeout: -1,
	}

	Rdb = redis.NewClient(opt)
	//ssh好像说需要ssh跳板机
	//而且本身也有加密（redis加密要带tls）
	//用ssh密钥就能登录
	//这里直接不用tls
	//之后实践用到再说
}
