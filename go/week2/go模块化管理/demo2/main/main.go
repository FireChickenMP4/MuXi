package main

import (
	"fmt"
	"foo"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func main() {
	foo.Greet("xxx")
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
