package main

import (
	"fmt"
	"sort"
)

func main() {
	//对于切片排序，有sort包
	a := []int{5, 2, 8, 1, 9}
	fmt.Println(a, sort.IntsAreSorted(a))
	//想要在数组或切片中搜索一个元素，该数组或切片必须先被排序
	//（因为标准库的搜索算法使用的是二分法）
	// 然后，您就可以使用函数
	// func SearchInts(a []int, n int) int
	// 进行搜索，并返回对应结果的索引值。
	fmt.Println(sort.SearchInts(a, 2))
	//不可靠的，输出了0 ，实际上排完是a[0]此时是5
	sort.Ints(a)
	fmt.Println(a, sort.IntsAreSorted(a))
	fmt.Println(sort.SearchInts(a, 2))
	fmt.Println(sort.SearchInts(a, 10))
	//没有的话跟普通的搜索算法不大一样
	//其找的是可以安全插入这个值的地方，使得原本也还是升序
	fmt.Println(sort.SearchInts(a, -20))
}
