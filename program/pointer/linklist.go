package main

import "fmt"

type intlinklist struct {
	nextaddress *intlinklist
	value       int
}

/*
func initintlinklist(a *intlinklist, n int) intlinklist {
	if n == 0 {
		a.nextaddress = nil
		return a
	} else {
		//我想创建一个变量来着。。。。但是不会，所以初始化这一块
		//之后学到了再说吧
		//算是半成品了
		initintlinklist(a.nextaddress, n-1)
		//递归也不是很好
		return a
	}
}
*/
func main() {
	//学了一下指针，顺便写个链表玩玩
	/*
	   var x, y, z intlinklist
	   	x.nextaddress = &y
	   	y.nextaddress = &z
	   	z.nextaddress = nil
	   	x.value = 1
	   	y.value = 2
	   	z.value = 3
	   	var a *intlinklist = &x
	   	for {
	   		fmt.Println(a.value)
	   		if a.nextaddress == nil {
	   			break
	   		} else {
	   			a = a.nextaddress
	   		}
	   	}
	   	//go的结构体a.value自动就等价于c中的a->value或者说(*a).value
	   	//相当于简化了
	*/
	var start intlinklist
	//	start = initintlinklist(start.nextaddress, 10)
	//我想做到的效果可能是，返回个初始点的就可以
	//但是剩余的是靠函数体自己创建
	var a *intlinklist = &start
	for {
		fmt.Println(a.value)
		if a.nextaddress == nil {
			break
		} else {
			a = a.nextaddress
		}
	}
}
