package learn

import (
	"context"
	"fmt"
	"log"
	"week12/config"

	"github.com/redis/go-redis/v9"
)

func RedisLearn() {
	//每个 Redis 命令都接受一个上下文
	//你可以使用它来设置 超时 或传播一些信息
	ctx := context.Background()
	val, err := config.Rdb.Get(ctx, "key").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("key does not exists")
			return
		}
		log.Fatal(err)
	}
	fmt.Println(val)
	//或者可以保存命令，并在稍后分别访问值和错误
	get := config.Rdb.Get(ctx, "key")
	fmt.Println(get.Val(), get.Err())
	// Administrator in MuXi\go\week10 main*​​​ ⇡
	// ❯ go run main.go
	// 2026/02/02 22:38:10 redis: nil
	// exit status 1
	// Administrator in MuXi\go\week10 main*​​​ ⇡
	// ❯ go run main.go
	// 2
	// 2 <nil>

	//执行不支持的命令
	//也就是执行任意或者自定义命令
	val2, err := config.Rdb.Do(ctx, "get", "key").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("key does not exists")
			return
		}
		log.Fatal(err)
	}
	fmt.Println(val2.(string))
	//Do 返回一个Cmd，它有一堆助手来处理interface{}值
	// s, err := cmd.Text()
	// flag, err := cmd.Bool()

	// num, err := cmd.Int()
	// num, err := cmd.Int64()
	// num, err := cmd.Uint64()
	// num, err := cmd.Float32()
	// num, err := cmd.Float64()

	// ss, err := cmd.StringSlice()
	// ns, err := cmd.Int64Slice()
	// ns, err := cmd.Uint64Slice()
	// fs, err := cmd.Float32Slice()
	// fs, err := cmd.Float64Slice()
	// bs, err := cmd.BoolSlice()

	//然后关于redis.Nil刚才已经用过了，这是一个示例来大概区分一下
	// val, err := rdb.Get(ctx, "key").Result()
	// switch {
	// case err == redis.Nil:
	// 	fmt.Println("key does not exist")
	// case err != nil:
	// 	fmt.Println("Get failed", err)
	// case val == "":
	// 	fmt.Println("value is empty")
	// }

	//Conn 代表单个Redis连接，而不是连接池。
	//优先从client运行命令，除非有对连续单个Redis连接的特殊需求
	//有的话就用Conn
	//比如订阅持久连接
	//事务
	//连接特定的配置
	//然后需要注意的是必须显式关闭连接
	//不要滥用独立链接，为每个请求创建独立链接不如直接用连接池
	//以及并发访问同一个连接会出错！
	cn := config.Rdb.Conn()
	defer cn.Close()
}
