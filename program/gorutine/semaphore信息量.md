#### 14.2.9 用带缓冲通道实现一个信号量

> 信号量是实现互斥锁（排外锁）常见的同步机制，限制对资源的访问，解决读写问题，比如没有实现信号量的 sync 的 Go 包，使用带缓冲的通道可以轻松实现：

带缓冲通道的容量和要同步的资源容量相同
通道的长度（当前存放的元素个数）与当前资源被使用的数量相同
容量减去通道的长度就是未处理的资源个数（标准信号量的整数值）
不用管通道中存放的是什么，只关注长度；因此我们创建了一个长度可变但容量为 0（字节）的通道：

```go
type Empty interface {}
type semaphore chan Empty
```

将可用资源的数量 N 来初始化信号量 semaphore：`sem = make(semaphore, N)`

然后直接对信号量进行操作：

```go
// acquire n resources
func (s semaphore) P(n int) {
    e := new(Empty)
    for i := 0; i < n; i++ {
        s <- e
    }
}

// release n resouces
func (s semaphore) V(n int) {
    for i:= 0; i < n; i++{
        <- s
    }
}
```

可以用来实现一个互斥的例子：

```go
/* mutexes */
func (s semaphore) Lock() {
    s.P(1)
}

func (s semaphore) Unlock(){
    s.V(1)
}

/* signal-wait */
func (s semaphore) Wait(n int) {
    s.P(n)
}

func (s semaphore) Signal() {
    s.V(1)
}
```
