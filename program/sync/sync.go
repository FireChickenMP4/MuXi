package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println(1)
	//今天来试试sync包
	wg := &sync.WaitGroup{}
	mycounter := func() func(id int) {
		count := 0
		mutex := &sync.Mutex{}
		return func(id int) {
			mutex.Lock()
			count++
			fmt.Printf("id:%d,count:%d\n", id, count)
			mutex.Unlock()
			wg.Done() //其实就是wg.Add(-1)
		}
	}()
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go mycounter(i)
	}
	wg.Wait() //在归零前一直等待

	//对于mutex
	//必须指出的是，在第一次被使用后，不能再对sync.Mutex进行复制。
	// （sync包的所有原语都一样）。
	// 如果结构体具有同步原语字段，则必须通过指针传递它。

	//然后还有个sync.RWMutex
	//它实现了sync .Locker接口
	//所以都可以调用Lock和Unlock方法
	//但是它还允许使用RLock和RUnlock进行并发读取
	//sync.RWMutex允许至少一个读锁（即可以多个同时读）或一个写锁存在
	//而sync.Mutex允许一个读锁或一个写锁存在
	//当有大量读操作且很少写操作时，RWMutex性能更好
	//如果读写相近，mutex即可
	//rwmutex可以读读存在，但是读写，和写写不能同时存在
	//所以读取效率更高
	//这里我想不到什么场景，就不写示例了，大概是知道了

	//sync.Map是并发版本Go的map
	//使用Store(interface {}，interface {})添加元素。
	// 使用Load(interface {}) interface {}检索元素。
	// 使用Delete(interface {})删除元素。
	// 使用LoadOrStore(interface {}，interface {}) (interface {}，bool)检索或添加之前不存在的元素。如果键之前在map中存在，则返回的布尔值为true。
	// 使用Range遍历元素。
	m := &sync.Map{}

	m.Store(1, "one")
	m.Store(2, "two")
	//store添加元素

	val, contains := m.Load(1)
	temp := func() {
		if contains {
			fmt.Printf("%s\n", val.(string))
			//这里因为传入是空接口，所以要类型断言
		} else {
			fmt.Println("Not Found.")
		}
	}
	temp()
	m.Delete(1)
	val, contains = m.Load(1)
	temp()
	val, loaded := m.LoadOrStore(3, "three")
	if !loaded {
		fmt.Printf("%s\n", val.(string))
	}
	val, loaded = m.LoadOrStore(2, "二")
	if !loaded {
		fmt.Printf("%s\n", val.(string))
	}
	//所以这里loaded是是否已经存在
	m.Range(func(key, value any) bool {
		fmt.Printf("%d: %s\n", key.(int), value.(string))
		return true
	})
	//貌似这里map也还是无序的
	//out:
	/*
		one
		Not Found.
		three
		2: two
		3: three
		*or
		one
		Not Found.
		three
		3: three
		2: two
	*/
	//如你所见，Range方法接收一个类型为
	// func(key，value interface {})bool的函数参数。
	// 如果函数返回了false，则停止迭代。
	// 有趣的事实是，即使我们在恒定时间后返回false
	// 最坏情况下的时间复杂度仍为O(n)。
	//所以在sync .map和map+sync.Mutex之间，什么时候选前者？
	//如果我们对于map有频繁读取和不频繁写入的时候
	//当多个goroutine读取，写入和覆盖不相交的键时。
	//也就是说，如果我们有一个分片实现
	//其中一组由四个协程，每个协程负责25%的键（之间不冲突）
	//这种情况下sync.Map是首选
	//应该是因为其并发性

	//sync.Pool
	// sync.Pool是一个并发池，负责安全地保存一组对象。它有两个导出方法：

	// Get() interface{} 用来从并发池中取出元素。
	// Put(interface{}) 将一个对象加入并发池。
	type SendReq struct {
		id   int
		code int
		data string
	}
	NewConnection := func(id int) interface{} {
		a := SendReq{
			id,
			300,
			"你好",
		}
		return &a
		//这里pool需要返回的是指针(
	}
	pool := &sync.Pool{}

	pool.Put(NewConnection(1))
	pool.Put(NewConnection(2))
	pool.Put(NewConnection(3))

	getMes := func() {
		connection := pool.Get().(*SendReq) //顺便给类型断言了
		fmt.Printf("%d\n", connection.id)
	}
	repeat := func(times int, a func()) {
		for i := 0; i < times; i++ {
			a()
			//a是任意一个函数
		}
	}
	repeat(3, getMes)
	//out
	// 1
	// 3
	// 2
	//需要注意的是Get()方法会从并发池中随机取出对象
	//无法保证以固定的顺序获取并发池中存储的对象。

	//还可以为sync.Pool指定一个创建者方法:

	pool1 := &sync.Pool{
		New: func() any {
			return NewConnection(1)
		},
	}
	getMes = func() {
		connection := pool1.Get().(*SendReq) //顺便给类型断言了
		fmt.Printf("%d\n", connection.id)
	}
	getMes()
	getMes()
	//这样每次调用Get()时
	//将返回由在pool.New中指定的函数创建的对象（在本例中为指针）。
	//这里的意思是如果池中没有对象,但是发get请求的话
	//get会自动new一个默认的对象并返回
	//如果这里没有创建者方法，get会返回空值nil

	//什么时候用pool呢？
	//第一个是当我们必须重用共享的和长期存在的对象（例如，数据库连接）
	//第二个是用于优化内存分配。

	//让我们考虑一个写入缓冲区并将结果持久保存到文件中的函数示例
	//使用sync.Pool我们可以通过在不同的函数调用之间重用同一对象来重用为缓冲区分配的空间。
	//第一步是检索先前分配的缓冲区（如果是第一个调用，则创建一个缓冲区，但这是抽象的）。
	//然后，defer操作是将缓冲区放回sync.Pool中。
	//见下面writeFile函数
	//我也不太懂第一个，还有这个函数什么意思

	//sync.Once
	//我才发现有这么个好东西
	//sync.Once是一个简单而强大的原语，可确保一个函数仅执行一次。
	once := &sync.Once{}
	for i := 0; i < 4; i++ {
		i := i //这个示例搞什么玩意。。。,非整这局部变量
		//哦，goroutine异步
		//后面有读取操作
		//那我明白了
		go func() {
			once.Do(
				func() {
					fmt.Printf("first %d\n", i)
				}, //要加逗号，异行
				//同行可以不加})
			)
		}()
	}
	//我们使用了Do(func ())方法来指定只能被调用一次的部分。
	time.Sleep(1e9)
	//这里first不一定是1，可能是3，取决于协程之间

	//sync.cond
	//最后一个了
	//sync.Cond可能是sync包提供的同步原语中最不常用的一个
	//它用于发出信号（一对一）或广播信号（一对多）到goroutine。
	//让我们考虑一个场景，我们必须向一个goroutine指示共享切片的第一个元素已更新
	//创建sync.Cond需要sync.Locker对象（sync.Mutex或sync.RWMutex）：
	cond := sync.NewCond(&sync.Mutex{})
	//然后，让我们编写负责显示切片的第一个元素的函数：
	printFirstElement := func(s []int, cond *sync.Cond) {
		// defer wg.Done()
		cond.L.Lock()
		cond.Wait()
		//这里的wait是
		//先释放锁
		//等signal
		//等到了重新拿锁然后修改数据
		//所以wait方法需要持有锁的时候才可以调用
		//否则会panic
		//这与wg的wait并不是一个东西
		//我的理解的话，相当于换了条道阻塞，负责通知
		fmt.Printf("%d\n", s[0])
		cond.L.Unlock()
	}
	//我们可以使用cond.L访问内部的互斥锁
	//一旦获得了锁，我们将调用cond.Wait()
	//这会让当前goroutine在收到信号前一直处于阻塞状态

	//让我们回到main goroutine
	//我们将通过传递共享切片和先前创建的sync.Cond来创建printFirstElement池
	//然后我们调用get()函数，将结果存储在s[0]中并发出信号：
	s := make([]int, 1)
	// wg.Add(runtime.NumCPU())
	fmt.Println("逻辑 CPU 核心数：", runtime.NumCPU()) //获得cpu数这是
	for i := 0; i < runtime.NumCPU(); i++ {
		go printFirstElement(s, cond)
	}
	cond.L.Lock()
	s[0] = 1
	//发signal之前也要持有锁
	cond.Signal()
	//并且这个只叫醒一个
	cond.L.Unlock()
	// md这样等待的也太慢,不一定发信号的时候有人在等。。。
	// 无语了，好麻烦

	//这个信号会解除一个goroutine的阻塞状态
	//解除阻goroutine将会显示s[0]中存储的值。
	//但是，有的人可能会争辩说我们的代码破坏了Go的最基本原则之一：
	//不要通过共享内存进行通信；而是通过通信共享内存。
	//确实，在这个示例中，最好使用channel来传递get()返回的值
	//但是我们也提到了sync.Cond也可以用于广播信号
	//我们修改一下上面的示例，把Signal()调用改为调用Broadcast()。
	//在这种情况下，所有goroutine都将被触发

	// cond.L.Lock()
	// s[0] = 1
	// cond.Broadcast()
	// cond.L.Unlock()
	// wg.Wait()
	// md这样等待的也太慢，发broadcast的不一定都在等,然后就导致有部分没听到
	// 然后就会死锁

	//众所周知，channel里的元素只会由一个goroutine接收到
	//通过channel模拟广播的唯一方法是关闭channel。
	//当一个channel被关闭后，channel中已经发送的数据都被成功接收后
	//后续的接收操作将不再阻塞，它们会立即返回一个零值。
	//但是这种方法智能广播一次
	//因此尽管争议很大，但是这无疑是cond的一个有趣功能
}
func writeFile(pool *sync.Pool, filename string) error {
	buf := pool.Get().(*bytes.Buffer)

	defer pool.Put(buf)

	// Reset 缓存区，不然会连接上次调用时保存在缓存区里的字符串foo
	// 编程foofoo 以此类推
	buf.Reset()

	buf.WriteString("foo")

	return os.WriteFile(filename, buf.Bytes(), 0644)
}
