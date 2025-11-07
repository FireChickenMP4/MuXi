package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4} //初始化也可以用make
	//数组的话[...]int{xxx} ...是编译器自行推断大小
	str := [...]string{3: "awa", 5: "qwq"} //index:val可以赋值对应索引的
	fmt.Println(len(str), cap(str))        //大小开正好的
	for index, val := range str {
		fmt.Println("index:", index, "value:", val)
	}
	a = append(a, 5) //append追加，空切片也可以
	b := make([]int, len(a))
	//cap(a)是上限容量
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
}
