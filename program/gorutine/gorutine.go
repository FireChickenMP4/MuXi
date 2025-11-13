package main

import (
	"fmt"
	"time"
)

type Empty interface{}
type semaphore chan Empty

var (
	MUTEX semaphore
	//这里还是nil，主函数里记得make给分配一下内存空间
	//并且需要有一个缓冲，要不然直接锁死
)

func main() {
	//通道的创建
	//var ch1 chan string
	//ch1 = make(chan string)
	//当然可以用短声明
	ch1 := make(chan string) //1缓冲也可以避免死锁 ch1 := make(chan string,1)
	//通道传递的类型可以是任何一种
	//空接口也可以，但我还不知道是不是需要类型断言
	//甚至传递通道的通道，传递函数的通道，都是可以的
	//通道实际上是类型化消息的队列：使数据得以传输。它是先进先出（FIFO）结构的所以可以保证发送给他们的元素的顺序
	//go中无缓冲通道是同步的完美工具,这样就实现了通过通讯来共享数据

	//通道操作符是<-
	//ch1 <- "A"
	//str1 := <-ch1 X错误的
	//emm，go的通道好像是为了用于gorutine之间的交流，所以没协程在的时候，无缓存传数据的话会fatal err，deadlock  并且！死锁不是真的panic，但是还是有可能把严重错误叫做panic注意！
	//但是go的话可以用缓冲解决
	//好像C++其实也是这样的，但是c++我记得没缓冲的话也不至于死锁
	//C++好像没有？？？？？？？我记忆从哪里来的

	//哦我知道了，无缓存的话，ch1 <- "A"由于发了一个，主gorutine阻塞，等待其他gorutine接受数据，但没有其他gorutine可以接受，main的主gorutine停了，彻底私了，所以就死锁了
	//然后说main函数结束，进程就结束了，协程也会挂掉，所以要让main函数等一下
	go func() {
		str1 := <-ch1
		fmt.Println(str1)
	}() //这样可以解决
	//或者发送放到闭包的gorutine里，让主gorutine不会阻塞
	ch1 <- "A"
	testText := []string{"qwq", "awa", "=w=", "1234567", "你好，世界！"}
	//另外，defer没用，因为defer需要正常返回才行，死锁彻底锁在那里了，并不会返回
	go sendData(ch1, testText...)
	go getData(ch1)
	time.Sleep(1e9)

	//然后说我们有的时候可能是需要有缓存的，缓存用完之前不会阻塞
	// ch1 :=make(chan string,buf) buf是缓存数量
	//如果容量大于 0，通道就是异步的了，拿现在这个程序举例，A就有可能不会被str1接收
	//也就是存在竞态
	//当执行顺序很重要时，必须使用sync.WaitGroup、sync.Mutex或额外的通道进行同步

	//也可以使用通道来达到同步的目的，这个很有效的用法在传统计算机中称为信号量（semaphore）。或者换个方式：通过通道发送信号告知处理已经完成（在协程中）。
	//让我们用带缓冲通道实现一个信号量
	//我们不关心通道里具体放了点啥，只关心其长度 所以用空接口就行

	//这里如果说同时有三个sendData就已经明显的错乱了
	go sendData(ch1, testText...)
	go sendData(ch1, testText...)
	go sendData(ch1, testText...)
	go getData(ch1)
	time.Sleep(3e9)
	//例如输出：
	//qwq =w= qwq qwq 1234567 awa awa 你好，世界！ =w= =w= 1234567 1234567 你好，世界！ 你好，世界！ awa
	//awa =w= qwq qwq 1234567 awa awa 你好，世界！ =w= =w= 1234567 1234567 你好，世界！ 你好，世界！ qwq
	//...我也不知道为什么甚至开头会有个awa，而不是qwq，可能对于getData调度来说先调度到了
	//感觉可能是因为第一组的影响，我把第一组注释掉就没影响了

	//然后试试自己写的互斥锁
	//s的容量取决于需要等多少，n个发信号的就是n
	//1个的可以无缓冲
	//貌似不对，不懂（（（
	//感觉好像是许可什么的，但是说如果是许可，P和V要反一下
	//相当于令牌
	//s用于主函数的wait
	//算了在week2里面搞了，顺便做作业了算是
}

func sendData[T any](ch chan T, data ...T) {
	for _, val := range data {
		ch <- val
	}
	//多个发的话不要关通道
}

func getData[T any](ch chan T) {
	// time.Sleep(1e9)
	for input := range ch { //或者 input=<-ch也可以传数据，但是是死循环
		fmt.Printf("%v ", input)
	}
}

// 然后直接对信号量进行操作：
func (s semaphore) P(n int) {
	e := new(Empty)
	for i := 0; i < n; i++ {
		s <- e
	}
}
func (s semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
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

/*func sendDataMutex[T any](s semaphore, ch chan T, data ...T) {
	MUTEX.Lock()
	sendData(ch, data...)
	MUTEX.UnLock()
	//s.Signal() 我目前示例里面不需要send发送信号
}*/
//这里过度锁定了，导致变成串发了
//不如说sendData只是往通道里输入数据，要是锁的话就失去并发的意义了
