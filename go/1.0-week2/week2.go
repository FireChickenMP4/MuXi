package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type struct1 struct {
	name string
	id   int
}

type struct2 struct {
	id []int
}

type Empty struct{}
type semaphore chan Empty

func main() {
	//结构体比较
	a := new(struct1)
	b := new(struct1)
	c := new(struct1)
	a.name = "qwq"
	c.name = "qwq"
	b.name = "awa"
	a.id = 1
	c.id = 1
	b.id = 2
	if *a == *c {
		fmt.Println("ac")
	}
	if *a == *b {
		fmt.Println("ab")
	}
	//结构体可以直接比较
	//不同名称结构体的话可以强制转换比较
	//含有指针的话不行，指针地址不同
	//但是含有切片的不能
	//可以用reflect.DeepEqual比较
	d := new(struct2)
	e := new(struct2)
	d.id = []int{1, 2}
	e.id = []int{1, 2}
	if reflect.DeepEqual(d, e) {
		fmt.Println("de")
	}
	//注意，结构体不能make，然后map别new，new完没法放进去数据
	//map和slice扩容，slice是原本cap乘二，map是+1，所以map如果快速增长或者知道大概所需，可以make(map[xxx]xxx,cap)
	//哦还有，map不是cap，是len是容量，而且跟切片容量概念不同，大概
	//map键类型需要可以==和！=比较，或者需要能转化为数字，类似hash值一样的东西
	//如果要用结构体作为 key 可以提供 Key() 和 Hash() 方法，这样可以通过结构体的域计算出唯一的数字或者字符串的 key。
	m := make(map[string]int, 3) //我觉得是相当于减小扩容性能消耗了
	m["qwq"] = 1
	m["awa"] = 2
	fmt.Println(len(m))
	m["=w="] = 3
	m["-w-"] = 3
	fmt.Println(len(m))

	//利用闭包实现一个计数器
	mycount := func() func() {
		ans := 0
		return func() {
			ans++
			fmt.Printf("%d ", ans)
		}
	}()
	for range 10 {
		mycount()
	}

	//我们希望多个协程一块累加计数器（
	//从10之后继续搞
	for i := 0; i < 10; i++ {
		go mycount()
	}
	//如果这样会存在竞态，乱掉
	time.Sleep(1e9)
	//out:11 13 12 20 15 16 17 18 19 14
	fmt.Println()

	//加个互斥锁试试
	numGorutines := 5
	s := make(semaphore, numGorutines) //用于wait完成
	mycount = func() func() {
		mutex := make(semaphore, 1)
		counter := struct {
			ans   int
			mutex semaphore
		}{
			ans:   0,
			mutex: mutex,
		}
		counter.mutex.V(1)
		//先解锁
		return func() {
			counter.mutex.Lock()
			defer counter.mutex.UnLock()
			counter.ans++
			fmt.Printf("%d ", counter.ans)
		}
	}()
	for i := 0; i < numGorutines; i++ {
		go func(id int) {
			for j := 1; j <= 100; j++ {
				mycount()
			}
			s.Signal()
		}(i)
	}
	s.Wait(numGorutines)
	fmt.Println()

	//其实一般直接用sync的Mutex和WaitGroup就可以

	//然后做一下作业
	type idnum struct {
		num int
		id  int
	}
	sortedChans := make([]chan idnum, 20)
	for i := range sortedChans {
		sortedChans[i] = make(chan idnum, 1)
	}
	numGorutines = 20
	s = make(semaphore, numGorutines)
	randomNum := func() func(id int) {
		return func(id int) {
			rt := rand.Intn(1000)
			r := rand.Int()
			time.Sleep(time.Duration(rt))
			temp := idnum{
				num: r,
				id:  id,
			}
			fmt.Printf("num:%d,id:%d\n", r, id)
			sortedChans[id-1] <- temp
			s.Signal()
		}
	}()
	for i := 1; i <= 20; i++ {
		go randomNum(i)
	}
	s.Wait(numGorutines)
	fmt.Println()
	for i := 0; i < 20; i++ {
		d := <-sortedChans[i]
		fmt.Printf("num:%d,id:%d\n", d.num, d.id)
	}
	//非要说其实可以，emm，循环着读，读完再给塞回去，除非是对应序号的

	//然后是交替输出
	//交替每个输出俩，所以写个每次往通道里倒俩的互斥锁
	chstr := make(chan byte, 200)
	numGorutines = 2
	s = make(semaphore, numGorutines)

	sendDataStr := func() func(id int) {
		strA0 := []string{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789"}
		//要打印的字符串
		//最好的实践方法是把控制单独提出来做个中央控制的协程
		//我们只需保证发送一定是互锁的就可以了
		mutex := make([]semaphore, 2)
		for i := 0; i < 2; i++ {
			mutex[i] = make(semaphore, 1)
		}
		unfinished := []bool{true, true}
		mutex[0].V(1)
		//初始给字母字符串放一个令牌
		func1 := func(id int) {
			for i := 0; i < len(strA0[id]); i += 2 {
				mutex[id].Lock()

				if i < len(strA0[id]) {
					chstr <- strA0[id][i]
				}
				if i+1 < len(strA0[id]) {
					chstr <- strA0[id][i+1]
				}

				//检查是否还有字符
				unfinished[id] = (i+2 < len(strA0[id]))
				//智能释放令牌
				if unfinished[1-id] {
					//交给对方
					mutex[1-id].UnLock()
				} else if unfinished[id] && !unfinished[1-id] {
					//给自己放令牌
					mutex[id].UnLock()
				}
				//两者都没有更多，不放任何令牌
			}
			s.Signal() //告诉等待gorutine俺完事了
		}
		return func1
	}()
	go func() {
		go sendDataStr(0)
		go sendDataStr(1)

		s.Wait(numGorutines)
		close(chstr)
	}()

	for str := range chstr {
		fmt.Printf("%s", string(str))
	}
	fmt.Println()
}

// 然后直接对信号量进行操作：
// p获取n个令牌
func (s semaphore) P(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

// V释放n个令牌
func (s semaphore) V(n int) {
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
}

// 以及一个互斥锁
func (s semaphore) Lock() {
	s.P(1)
}
func (s semaphore) UnLock() {
	s.V(1)
}
func (s semaphore) Wait(n int) {
	s.P(n)
}
func (s semaphore) Signal() {
	s.V(1)
}
