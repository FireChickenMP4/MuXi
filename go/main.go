package main

import "fmt"

func main() {
	x := []int{1, 2, 3}
	y := x[1:] //2之后 2 3
	y = append(y, 4)
	for _, val := range x {
		fmt.Println(val) //输出1 2 3 而没有4
	}
	fmt.Print("\n")
	for _, val := range y {
		fmt.Println(val) //输出2 3 4 而没有1
	}
	fmt.Print("\n")
	//因为append在cap不够会自动扩容，导致y指向新的底层数组
	//这过程中只是数据的拷贝
	a := make([]int, 3, 5)
	a[0] = 1
	a[1] = 2
	a[2] = 3
	b := a[1:3]
	b = append(b, 4)
	fmt.Println(a)
	fmt.Print("\n")
	fmt.Println(b)
	fmt.Print("\n")
	//此时append容量是足够的，所以a,b仍然是同一个底层数组
	//但是a却还是123，是因为len此时是3，如果将a扩展到4就可以看到了
	aex := a[:4]
	fmt.Println(aex)
	/*[1 2 3]

	[2 3 4]

	[1 2 3 4]
	*/
}
