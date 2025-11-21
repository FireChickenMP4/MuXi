package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

func main() {
	//来来，让我们系统学一下strings包

	//1.复制字符串
	ori := "hello 世界！"
	copys := strings.Clone(ori)
	fmt.Println(ori, copys)
	fmt.Println(&ori, &copys)
	//out :
	// hello 世界！ hello 世界！
	// 0xc00011c030 0xc00011c040

	//2.比较字符串
	//将a,b按照字典顺序比较（中文怎办？好吧，中文比的是Unicode编码顺序）
	//a>b ret 1 a<b ret -1 a==b ret 0
	fmt.Println(strings.Compare("abc", "abe"))
	fmt.Println(strings.Compare("abcd", "abe"))
	fmt.Println(strings.Compare("abijk", "abe"))
	fmt.Println(strings.Compare("abe", "abe"))
	//out :
	// -1
	// -1
	// 1
	// 0

	//3.包含字符串
	//连续包含
	fmt.Println(strings.Contains("abcdefg", "a"))
	fmt.Println(strings.Contains("abcdefg", "abc"))
	fmt.Println(strings.Contains("abcdefg", "ba"))
	//out :
	// true
	// true
	// false
	//不连续包含（就是存在有这几个单个的）
	fmt.Println(strings.ContainsAny("abcdefg", "bac"))
	fmt.Println(strings.ContainsAny("abcdefg", "gfdecba"))
	fmt.Println(strings.ContainsAny("abcdefg", "h"))
	//out :
	// true
	// true
	// false
	//字符
	fmt.Println(strings.ContainsRune("abcedf", 'a'))
	fmt.Println(strings.ContainsRune("abcedf", '0'))
	fmt.Println(strings.ContainsRune("你好世界", '你'))
	//out :
	// true
	// false
	// true
	//其实还有个ContainsFunc

	//4.子串出现次数
	fmt.Println(strings.Count("3.1415926", "1"))
	fmt.Println(strings.Count("there is a girl", "e"))
	fmt.Println(strings.Count("there is a girl", " "))
	fmt.Println(strings.Count("there is a girl", ""))
	//  					   1234567890123456
	fmt.Println(len("there is a girl"))
	//out :
	// 2
	// 2
	// 3
	// 16 len+1相当于是说，应该是算上结束符了，或者说检索的时候到索引end了
	// 15

	//5.删除指定子串
	//func Cut(s, sep string) (before, after string, found bool)
	//这里展示一下方法 删除 s 内第一次出现的字串 sep, 并返回删除后的结果
	//before-被删除字串位置前面的字符串
	//after-被删除字串位置后面的字符串
	//found-是否找到子串
	show := func(s, sep string) {
		before, after, found := strings.Cut(s, sep)
		fmt.Printf("Cut(%q, %q) = %q, %q, %v\n", s, sep, before, after, found)
		//%q该值对应的双引号括起来的go语法字符串字面值，必要时会采用安全的转义表示
	}
	show("Hello world world", " ")
	show("Hello world world", "world")
	show("Hello world world", "Hello")
	show("Hello world world", "Hello world world")
	//out :
	// Cut("Hello world world", " ") = "Hello", "world world", true
	// Cut("Hello world world", "world") = "Hello ", " world", true
	// Cut("Hello world world", "Hello") = "", " world world", true
	// Cut("Hello world world", "Hello world world") = "", "", true
	myCutAll := func(s, sep string) (str string) {
		if sep == "" {
			return s
		}
		needCut := true
		var before, after string
		for needCut {
			before, after, needCut = strings.Cut(s, sep)
			s = before + after //简单相加了就，简单粗暴
			//虽然时间复杂度贼拉高 o(n^2)
			//好像是说每次加都要重新分配一次内存啥的
			//然后再拷贝到新的，相当于1加到n
			//而且转换这里也费事
			//用buffer好一些
			//但是strings包有自己的builder，连接用的
			//两者的对比后面我也会说
		}
		str = before
		return str
	}
	fmt.Println(myCutAll("Hello world world!", "world"))
	fmt.Println(myCutAll("Hello world world!", " "))
	fmt.Println(myCutAll("Hello world world!", ""))
	//out :
	// Hello  !
	// Helloworldworld!
	// Hello world world!
	//这里我一开始没判断sep为空，死循环了，Cut本身没这个问题

	//6.忽略大小写相等
	fmt.Println(strings.EqualFold("你好", "你好"))
	fmt.Println(strings.EqualFold("Hello", "Hello"))
	fmt.Println(strings.EqualFold("Hello", "hELLO"))
	//out :
	// true
	// true
	// true

	//7.分割字符串
	//有两个这个
	//func Fields(s string) []string
	//func FieldsFunc(s string, f func(rune) bool) []string
	//前者是根据空格来分割字符串
	//后者是函数 f 的返回值来决定是否分割字符串
	fmt.Printf("%q\n", strings.Fields(" a b c d e f g "))
	fmt.Printf("%q\n", strings.FieldsFunc("a,b,c,d,e,f,g", func(r rune) bool {
		return r == ','
	}))
	fmt.Printf("%q\n", strings.FieldsFunc("abcdefg", func(r rune) bool {
		return true
	}))
	fmt.Printf("%q\n", strings.FieldsFunc("abcdefg", func(r rune) bool {
		return r == ' '
	}))
	//out :
	// ["a" "b" "c" "d" "e" "f" "g"]
	// ["a" "b" "c" "d" "e" "f" "g"]
	// []
	// ["abcdefg"]

	//8.寻找前后缀
	// func HasPrefix(s, prefix string) bool
	// func HasSuffix(s, suffix string) bool
	//前者是寻找前缀，后者是寻找后缀
	//感兴趣可以去看看这里的源码实现，比较巧妙
	fmt.Println(strings.HasPrefix("abbc cbba", "abb"))
	fmt.Println(strings.HasSuffix("abbc cbba", "bba"))
	//out :
	// true
	// true

	//9.子串的位置
	//找第一次出现的一共有三个
	//func Index(s, substr string) int
	//func IndexAny(s, chars string) int
	//func IndexRune(s string, r rune) int
	//这里与Contains同理		0123456
	fmt.Println(strings.Index("abcdefg", "bc"))
	fmt.Println(strings.IndexAny("abcdefg", "eb"))
	fmt.Println(strings.IndexRune("abcdefg", 'g'))
	fmt.Println(strings.IndexRune("abcdefg", '\x00'))
	//out :
	// 1
	// 1
	// 6
	// -1
	//这里说一下，Go语言的字符串没有特定的结束符
	//。。。和其他语言的输出比较的话
	//会出现 c的aaa和go的aaa不相等，因为c的有\0
	//最后一个出现的同理，只不过前缀改成了LastIndex
	fmt.Println(strings.LastIndex("abcdefga", "a"))
	fmt.Println(strings.LastIndexAny("abcdefghisa", "ba"))
	//out :
	// 7
	// 10
	//然后还有个Func版本

	//10.遍历替换字符串
	//func Map(mapping func(rune) rune, s string) string
	// Map 返回字符串 s 的副本
	//并根据映射函数修改字符串 s 的所有字符
	//如果映射返回负值，则从字符串中删除该字符，不进行替换
	fmt.Println(strings.Map(func(r rune) rune {
		return r - 32
	}, "abcdefghijk"))
	fmt.Println(strings.Map(func(r rune) rune {
		return r + 32
	}, "ABCDEFGHIJK"))
	fmt.Println(strings.Map(func(r rune) rune {
		if r < 'F' {
			return -1
		} else {
			return r
		}
	}, "ABCDEFGHIJK"))
	fmt.Println(strings.Map(func(r rune) rune {
		if r != '!' {
			return '\x20'
			//x20是空格
		} else {
			return r
		}
	}, "qwq!"))
	//out :
	// ABCDEFGHIJK
	// abcdefghijk
	// FGHIJK
	//    !

	//11.重复复制字符串
	//根据给定的 Count 复制字符串，如果为负数会导致panic
	fmt.Println(strings.Repeat("a", 10))
	fmt.Println(strings.Repeat("abc", 10))
	tempStr := strings.Repeat("abc", 0)
	if tempStr == "" {
		fmt.Println("Empty")
	}
	//out :
	// aaaaaaaaaa
	// abcabcabcabcabcabcabcabcabcabc
	// Empty

	//12.替换字符串
	//func Replace(s, old, new string, n int) string
	//s 为源字符串，old 指要被替换的部分，new 指 old 的替换部分
	//n 指的是替换次数，n 小于 0 时表示不限制替换次数
	fmt.Println(strings.Replace("Hello this is golang", "golang", "c++", 1))
	fmt.Println(strings.Replace("Hello this is golang", "o", "c", -1))
	fmt.Println(strings.Replace("Hello this is golang", "o", "c", 1))
	fmt.Println(strings.Replace("Hello this is golang", "l", "c", 2))
	// out :
	// Hello this is c++
	// Hellc this is gclang
	// Hellc this is golang
	// Hecco this is golang
	//有个Replace的方便函数
	//ReplaceAll 相当于最后是-1，不用输入了

	//13.分隔字符串
	//根据子串 sep 将字符串 s 分隔成一个字符串切片
	//func Split(s, sep string) []string
	//根据子串 sep 将字符串 s 分隔成一个字符串切片，其分隔次数由 n 决定
	//func SplitN(s, sep string, n int) []string
	//根据子串 sep 将字符串 s 分隔成包含 sep 的字符串元素组成的字符串切片
	//func SplitAfter(s, sep string) []string
	//在后面切，所以是After
	//根据子串 sep 将字符串 s 分隔成包含 sep 的字符串元素组成的字符串切片，其分隔次数由 n 决定
	//func SplitAfterN(s, sep string, n int) []string
	//然后教程里没写的两个SplitSeq和SplitSeqN 以及对应的After版本
	//1.23新加的
	//func SplitSeq(s, sep string) iter.Seq[string]
	//这个返回一个迭代器，要放到for里面循环输出
	//这个的好处是，是一下一下的出（当然也是一次性的，第二次循环就没了）
	//然后Split好像都有个特殊设计
	//空字符串会直接每个字符分开，你好变成你和好
	fmt.Printf("%q\n", strings.Split("this is go language", " "))
	fmt.Printf("%q\n", strings.Split(" this is go language", " "))
	fmt.Printf("%q\n", strings.SplitN("this is go language", " ", 2))
	fmt.Printf("%q\n", strings.SplitAfter("this is go language", " "))
	fmt.Printf("%q\n", strings.SplitAfterN("this is go language", " ", 2))
	tempIter := strings.SplitSeq("苹果，香蕉，凤梨", "，")
	for fruit := range tempIter {
		fmt.Println(fruit)
	}
	//out :
	// ["this" "is" "go" "language"]
	// ["" "this" "is" "go" "language"]
	// ["this" "is go language"]
	// ["this " "is " "go " "language"]
	// ["this " "is go language"]
	// 苹果
	// 香蕉
	// 凤梨
	//普通Split调用genSplit(s, sep, 0, -1)
	//SplitN调用genSplit(s, sep, 0, n)
	//然后对于genSplit，如果n小于0的话，会自动转化n为
	//Count+1（也就是找到对应子串的数量加一）
	//然后创建大小为n的切片
	//所以这个n的含义是切成多少份，而不是切多少次
	//所以SplitN(...,2)只会切一次
	//以及Split切空格，如果开头就有空格，那么索引0是空字符串

	//14.大小写转换
	//换成小写
	//func ToLower(s string) string
	//换成大写
	//func ToUpper(s string) string
	//两者均有个
	//func ToXXXSpecial(c unicode.SpecialCase, s string) string
	//根据传入对应语言的unicode.SpecialCase，转换成对应语言的大写字符串
	//个人觉得这个相当于是说映射表了，感觉甚至能自己自定义
	//unicode.SpecialCase 对，可以的，扒了源码了
	//SpecialCase定制（通过覆盖）标准映射
	/*
		type CaseRange struct {
		Lo    uint32
		Hi    uint32
		Delta d
		}
		type SpecialCase []CaseRange
		type d [MaxCase]rune
		const (
		UpperLower = MaxRune + 1 // (Cannot be a valid delta.)
		)
	*/
	//Lo 和 Hi 是起始和终止，表示特殊规则适用范围
	//然后Delta分为正常模式和交替模式
	//正常模式存储的是位移，0说明不变， 1 2 3位置分别是变成大写和小写以及标题需要的位移
	//Delta只有三个字（rune）
	//然后如果Delta三个数都是UpperLower(无效的数，超出Unicode范围了)
	//这样范围内就是大小写交替
	//1000到1003就是 1000大写-1001小写 1002大写-1003小写
	//行，偏完题了已经
	fmt.Println(strings.ToLower("My name is jack,Nice to meet you!"))
	fmt.Println(strings.ToLowerSpecial(unicode.TurkishCase, "Önnek İş"))
	fmt.Println(strings.ToUpper("My name is jack,Nice to meet you!"))
	fmt.Println(strings.ToUpperSpecial(unicode.TurkishCase, "örnek iş"))
	//这里就是土耳其语言的特殊情况，但是这几个字好像不在其中也
	/*
		var _TurkishCase = SpecialCase{
		CaseRange{0x0049, 0x0049, d{0, 0x131 - 0x49, 0}},
		CaseRange{0x0069, 0x0069, d{0x130 - 0x69, 0, 0x130 - 0x69}},
		CaseRange{0x0130, 0x0130, d{0, 0x69 - 0x130, 0}},
		CaseRange{0x0131, 0x0131, d{0x49 - 0x131, 0, 0x49 - 0x131}},
		}
		//TurkishCase重定向到_TurkishCase
	*/
	//out :
	// my name is jack,nice to meet you!
	// önnek iş
	// MY NAME IS JACK,NICE TO MEET YOU!
	// ÖRNEK İŞ

	//15.修剪字符串
	//小熟人，之前muxi考试用过 Trim类
	//一共是五个
	// func Trim(s, cutset string) string
	// func TrimLeft(s, cutset string) string
	// func TrimPrefix(s, suffix string) string
	// func TrimRight(s, cutset string) string
	// func TrimSuffix(s, suffix string) string
	//看名字就能知道干啥的，Trim是两端cutset任意匹配子串删除
	//Left是剪左端cutset任意匹配子串（R同理）
	//Prefix是将 cutset 匹配的子串删除（Suffix同理）
	//在处理字符串时，cutset 是一个用于删除特定字符集合的工具
	//然而，当 cutset 包含重复字符时，可能会导致意外行为或错误
	//cutset 是相当于any的作用
	//下面演示刻意重复是为了展现这一点
	//而suffix前后缀的话就是完整匹配
	fmt.Println(strings.Trim("!!this is a test statement!!", "!!!"))
	fmt.Println(strings.TrimLeft("!!this is a test statement!!", "!!!"))
	fmt.Println(strings.TrimRight("!!this is a test statement!!", "!!!"))
	fmt.Println(strings.TrimPrefix("!!this is a test statement!!", "!!!"))
	fmt.Println(strings.TrimSuffix("!!this is a test statement!!", "!!!"))
	//所以这里的话，前后缀更像是连续
	//而Left Right更像是any
	//然后Trim是全部的any
	//那怎么删全图的所有字符串呢
	//myAllCut:欸你是不是在找我 //自己写的
	//Fields,Split:我也可以啊
	//ReplaceAll:那我不是更方便吗
	//怎么删中间的字符串呢？ Cut:?
	//任意的字符呢？ CutAll:?
	//out :
	// this is a test statement
	// this is a test statement!!
	// !!this is a test statement
	// !!this is a test statement!!
	// !!this is a test statement!!

	//最后三个，我觉得算一些搭建的高效方法？
	//16.字符串 Builder
	//字符串 Builder 比起直接操作字符串更加节省内存。
	//strings.Builder :
	/*
		type Builder struct {
		addr *Builder // of receiver, to detect copies by value

		// External users should never get direct access to this buffer, since
		// the slice at some point will be converted to a string using unsafe, also
		// data between len(buf) and cap(buf) might be uninitialized.
		buf []byte
		}
	*/
	//一会讲对比和原理？
	builder := strings.Builder{}
	builder.WriteString("hello") //我好像还写过bytes.buffer类似的方法
	builder.WriteString(" world!")
	fmt.Println(builder.Len())
	fmt.Println(builder.String())
	//out :
	// 12
	// hello world!
	//原理：首先用的字节切片，没长度限制
	//然后扩容策略是类似c++ vector ,每次不够的时候乘以2
	//（builder.Grow()方法）
	//写时复制保护：addr字段用于检测非法并发访问
	/*
		// copyCheck implements a dynamic check to prevent modification after
		// copying a non-zero Builder, which would be unsafe (see #25907, #47276).
		//
		// We cannot add a noCopy field to Builder, to cause vet's copylocks
		// check to report copying, because copylocks cannot reliably
		// discriminate the zero and nonzero cases.
		func (b *Builder) copyCheck() {
			if b.addr == nil {
				// This hack works around a failing of Go's escape analysis
				// that was causing b to escape and be heap allocated.
				// See issue 23382.
				// TODO: once issue 7921 is fixed, this should be reverted to
				// just "b.addr = b".
				b.addr = (*Builder)(abi.NoEscape(unsafe.Pointer(b)))
			} else if b.addr != b {
				panic("strings: illegal use of non-zero Builder copied by value")
			}
		}
		//这个相当于给自己贴了个自己的地址标签
		//如果检查发现这个贴的标签不是自己的地址，说明是非法复制了一个非0的builder
		//这个设计是为了防止因为意外复制 Builder 而导致的难以发现的 bug
		//b.addr = (*Builder)(abi.NoEscape(unsafe.Pointer(b)))
		//这里使用特殊技巧避免内存问题
		//1.unsafe.Pointer(b)
		//将指针b转换为Go中最底层的指针类型
		//unsafe.Pointer是Go中唯一可以指向任何类型内存的指针
		//它会绕过类型安全检查
		//2.abi.NoEscape(...)
		//特殊的编译器提示，表示"这个指针不会逃出当前函数"
		//这是Go编译器的内部函数，用来覆盖编译器的逃逸分析机制
		//3.(*Builder)(...)
		//强制类型转换
		//这样就不会大费周折把指针放到堆里（因为编译器认为有可能会发到函数外使用）
		//直接放在栈里就好
		//注释中提到，这只是临时方案，未来Go编译器改进后（issue #7921）会恢复为简写法
		//但是23年到现在也没改咯。。。
	*/
	//然后写操作优化，主要是因为切片以及避免类型转换（buf也是如此）
	//直接操作字节切片而不是string
	//以及批量内存写入 s...语法实现内存块复制
	//然后底层原理上优化在于
	//1.逃逸分析优化
	//主要是栈，提前分配栈内存，以及说整个builder对象和缓冲区
	//都可能分配在栈上
	//2.字符串转换优化
	/*
		// strings/builder.go
		func (b *Builder) String() string {
		    return *(*string)(unsafe.Pointer(&b.buf))
		}
	​*/
	//unsafe转[]byte到string
	//避免常规转换时的内存拷贝
	//依赖Go底层字符串与切片的相同头部结构：
	/*
		type StringHeader struct {
		    Data uintptr
		    Len  int
		}
		type SliceHeader struct {
		    Data uintptr
		    Len  int
		    Cap  int
		}
	*/
	//然后实际使用的时候容量可以预分配一些
	//以及说可以手写封装并发安全
	//对比其他方案，比如bytes.buffer
	//strings.Builder更激进的内存分配策略（带来内存占用的更小？我觉得是这个因果）
	//方法集的最小化以及字符串转换无需拷贝
	//让其与bytes.Buffer的比较完整的庞然大物来对比
	//在需要频繁构建字符串的场景（如模板渲染、协议组装等）成为最佳选择
	//然后最后提醒，builder不要复制，不要复制，尤其是函数值传递

	//17.字符串 Replacer
	//Replacer replaces a list of strings with replacements.
	//It is safe for concurrent use by multiple goroutines.
	//这玩意自带并发安全
	//Replacer 转用于替换字符串
	//自动优化底层实现（根据替换规则数量选择算法）
	//一会扒了它 hiahiahia 好吧感觉水平还不够理解源码
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	//注意要传入一对，奇数会panic
	//替换关系< -> &lt; > -> &gt;
	fmt.Println(r.Replace("This is <b>HTML</b>!"))
	//out : This is &lt;b&gt;HTML&lt;/b&gt;!
	//这个也是很实用了qwq
	//这个替过程产是一次性完成的，有点像转义的感觉
	//特点:
	// 零内存分配（对于已初始化的 Replacer）
	// 自动处理嵌套替换（避免无限循环）
	//然后它还有个方法
	//func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
	//相当于将替换结果直接写入io.Writer了
	//特点：
	// 避免中间字符串分配（适合大文本处理）
	// 返回写入的字节数和可能的错误
	// 一般用于流式处理日志文件
	//优势的话，多次高频固定（不动态变化的）需要线程安全，用replacer
	//然后就是大量情况下，性能和内存效率极高，比replace快了3.7倍
	//内存零分配
	//然后建议全局复用，性能更好

	//18.字符串 Reader
	//Reader 实现了 io.Reader, io.ReaderAt,
	//io.ByteReader, io.ByteScanner,
	//io.RuneReader, io.RuneScanner,
	//io.Seeker, 和 io.WriterTo interfaces
	//io包感觉又是大工程啊，io包我还不是很懂
	//所以估计方法一堆堆的我也暂时学不来，所以暂时搁置！
	//优点应该就是，快！然后可以跟io对接
	//把字符串包装成一个高效的io.Reader
	//而且也是零分配内存，高性能的解决方案
	//可以使用各种基于io.Reader的流式处理方法
	reader := strings.NewReader("abcdefghijk")
	buffer := make([]byte, 20, 20)
	read, err := reader.Read(buffer)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(read)
	fmt.Println(string(buffer))
	//out :
	// 11
	// abcdefghijk
	// 12345678901

	//然后和string转换关系比较密切的包是strconv，详情请见strconv
	//在新建文件夹了
}
