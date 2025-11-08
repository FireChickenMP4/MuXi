package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4} //初始化也可以用make,new
	//a := make([]int,len,cap)
	//a := new([cap]int)[0:len]
	//个人理解其实是创建一个cap的底层数组
	//然后截取了个切片返回
	//make和new使用上没啥区别，但是底层上有些不同
	//make可能是在栈上分配，但好像说其实也是在堆，而new总是在堆上
	//并且make直接返回切片类型
	//new跟我理解一样，是先返回的数组指针，然后截取得到切片
	/*
		a := new([100]int)
		fmt.Printf("%T", a)
		输出：*[100]int
		所以真要输出数组要用解引用
	*/
	//推荐使用make，go官方给的好像
	//数组的话[...]int{xxx} ...是编译器自行推断大小
	str := [...]string{3: "awa", 5: "qwq"} //index:val可以赋值对应索引的
	fmt.Println(len(str), cap(str))        //大小开正好的
	str[1] = "aaa"
	for index, val := range str {
		fmt.Println("index:", index, "value:", val)
	}
	a = append(a, 5) //append追加，空切片也可以
	b := make([]int, len(a))
	//cap(a)是上限容量,我个人理解可以说是
	//底层数组的长度相当于是
	copy(b, a)
	c := a[:]
	a[0] = 2
	for _, val := range a {
		fmt.Print(val)
	}
	fmt.Println()
	for _, val := range b {
		fmt.Print(val)
	}
	fmt.Println()
	for _, val := range c {
		fmt.Print(val)
	}
	fmt.Println()
	//copy跟截取不同，就是不相关了

	//关于append,cap够就仍用原先底层数组
	//不够就创个新的更大的底层数组
	//把数据复制过去
	x := []int{1, 2, 3}
	y := x[1:] //2之后 2 3
	y = append(y, 4)
	fmt.Println(x) //输出1 2 3 而没有4
	fmt.Println(y) //输出2 3 4 而没有1
	//因为append在cap不够会自动扩容，导致y指向新的底层数组
	//这过程中只是数据的拷贝
	w := make([]int, 3, 5)
	w[0] = 1
	w[1] = 2
	w[2] = 3
	z := w[1:3]
	z = append(z, 4)
	fmt.Println(w)
	fmt.Print("\n")
	fmt.Println(z)
	fmt.Print("\n")
	//此时append容量是足够的，所以a,b仍然是同一个底层数组
	//但是a却还是123，是因为len此时是3，如果将a扩展到4就可以看到了
	//底层是[1 2 3 4 0]
	aex := w[:4]
	fmt.Println(aex)
	/*[1 2 3]

	[2 3 4]

	[1 2 3 4]
	*/
	p := []int{1, 2}
	fmt.Println(len(p), cap(p))
	p = append(p, 3)
	fmt.Println(len(p), cap(p))
	//append的扩容和通常的扩容类似
	//即cap翻倍

	//所以说go的数组是值
	//而go的切片更像是数组的指针
	//与c++概念相似
	/*type slice struct {
	    ptr uintptr   // 指向底层数组某个元素的指针
	    len int       // 长度
	    cap int       // 容量
		}*/

	//多维切片
	twod := make([][]int, 4) //4行
	twod[0] = []int{1, 2, 3}
	twod[1] = []int{4, 5}
	twod[2] = []int{6}
	for i := range twod {
		for j := range twod[i] {
			fmt.Printf("%d ", twod[i][j])
		}
		fmt.Println()
	}
	//二维切片的打印相当于是

	items := [...]int{10, 20, 30, 40, 50}
	for _, item := range items {
		item *= 2
	}
	fmt.Println(items) //发现原样，传进去的是值
	//slice resilce切片重组
	//前面有提到
	//此时append容量是足够的，所以a,b仍然是同一个底层数组
	//但是a却还是123，是因为len此时是3，如果将a扩展到4就可以看到了
	//底层是[1 2 3 4 0]
	//我们可以一步步扩展到底层数组的长度

	//因为字符串是纯粹不可变的字节数组，它们也可以被切分成 切片。
	//string -> []rune 字符切片 或者 []byte 后者处理的话比较常用貌似

	var by1, by2 []byte
	by1 = append(by1, "hello"...)
	//等价于一个字符一个字符追加
	by2 = append(by2, ",world!"...)
	by1 = append(by1, by2...)
	fmt.Println(by1)
	//append第二个接收的是元素，直接用切片并不行
	//故用...语法把切片元素展开

	//append的操作并不只是说只能添加，还能实现各种切片操作
	//删除索引为i的元素 arr = append(arr[:i], arr[i+1:]...)
	//删除i到j [:i],[j+1:]...
	//扩展j个元素长度 append(arr,make([]T,j)...)
	//在索引 i 的位置插入元素 x：a = append(a[:i], append([]T{x}, a[i:]...)...)
	//在索引 i 的位置插入长度为 j 的新切片：a = append(a[:i], append(make([]T, j), a[i:]...)...)
	//取出位于切片 a 最末尾的元素 x：x, a = a[len(a)-1], a[:len(a)-1]

	//切片和垃圾回收请见md
}
