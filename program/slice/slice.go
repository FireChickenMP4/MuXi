package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4} //初始化也可以用make
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
}
