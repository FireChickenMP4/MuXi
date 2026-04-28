## GO MODULE

**如何开启go mod/go work**

Go mod 在Golang 1.17引入，使用go mod对包进行管理时需要确保开启Go Module

`go env GO111MODULE` 可以查看是否开启

若未开启，输入 `go env -w GO111MODULE=on` 启用

Go有一个官方的代理仓库，它会根据版本及模块名缓存开发者下载过的模块。不过由于其服务器部署在国外，访问速度对于国内的用户不甚友好，我们需要修改默认的模块代理地址，运行以下指令设置七牛云代理 `go env -w GOPROXY=https://goproxy.cn,direct`

运行`go env`可以查看所有环境配置



**如何使用go mod/go work**

我们在编写程序时，把所有代码都写到一个文件中，会显得十分没有条理，那么如果将某些辅助函数放在一个单独的文件夹中，应该如何实现这些函数或者结构之间的调用呢？这里就需要go mod对函数和结构进行管理和调用。在GitHub中也有许多开源项目框架和工具，我们也可以使用go mod对开源的工具进行管理和使用。

在你想要创建go.mod文件的文件夹下运行指令`go mod init <模块名>`

```
demo1
├── model
│   ├── model.go
│   └── foo
│       └── model.go
├── test
│   └── test.go
├── go.mod
└── main.go
```

```go
package model

type Student struct {
	Name   string
	Age    int
	Scores []int
}
```

```go
package test

import (
	"demo/model"
	"fmt"
)

func ShowScores() {
	Stu := model.Student{
		Name:   "小明",
		Age:    18,
		Scores: []int{90, 80, 70},
	}

	fmt.Printf("姓名：%s，年龄：%d，成绩：%v\n", Stu.Name, Stu.Age, Stu.Scores)
}
```

```go
package foo

import "fmt"

func Foo() {
	fmt.Printf("Hello,world!")
}
```

```go
package main

import (
	"demo/model/foo"
	"demo/test"
)

func main() {
	test.ShowScores()
	foo.Foo()
}
```

可以看出go.mod相当于一个管理者，负责管理**所在文件目录及子目录**中所有结构、函数的调用



那么如果将某些函数再进行细分，分成多个go mod文件管理的项目，该如何对这些结构、函数进行调度呢？于是引入了go work工作区

在这些go mod的父文件夹运行指令`go work init`来创建一个工作区

在工作区所在文件夹运行指令`go work use <模块名>`将go mod添加到工作区

```
demo2
├── foo
│   ├── go.mod
│   ├── go.sum
│   └── foo.go
├── main
│   ├── main.go
│   └── go.mod
└── go.work
```

```go
package foo

import "fmt"

func Greet(name string) {
	fmt.Printf("Hello %s, welcome to Golang!", name)
}
```

```go
package main

import "foo"

func main() {
	foo.Greet("xxx")
}
```

go work相当于go mod的父级，负责不同go mod之间的调度



**go mod的其他用法**

在go.mod所在文件夹输入指令`go mod tidy`可以将你每个文件import的包清点并下载

例如以下情况，ants是外部导入的包，直接cv这份代码是无法运行的，会告诉你缺少包引入，那么这时可以在它所在的文件夹输入，它就会自动下载缺失的依赖，前提是能找到go.mod文件

```go
package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

func main() {
	p, _ := ants.NewPool(2, ants.WithNonblocking(true))
	defer p.Release()

	var wg sync.WaitGroup
	// 连续提交3个任务
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		err := p.Submit(wrapper(i, &wg))
		if err != nil {
			fmt.Printf("example:%d err:%v\n", i, err)
			wg.Done()
		}
	}

	wg.Wait()
}

func wrapper(i int, wg *sync.WaitGroup) func() {
	return func() {
		fmt.Printf("hello from example:%d\n", i)
		time.Sleep(1 * time.Second)
		wg.Done()
	}
}
```

