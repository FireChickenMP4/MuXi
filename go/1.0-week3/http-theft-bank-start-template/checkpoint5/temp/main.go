package main

import (
	"fmt"
)

func permute(nums []int) [][]int {
	//所谓全排列，其实就是说我定住任意一个数字，然后去全排列剩下所有的其他数字
	//所以可以用递归来搞，但是感觉dfs递归深度过大可能会崩掉
	//先把递归的逻辑写好得了
	if len(nums) == 1 {
		return [][]int{nums}
	}
	back := [][]int{}
	//作为返回的
	for i := range nums {
		//选择某个数为现在的位置
		num := []int{}
		for j, val := range nums {
			if j != i {
				num = append(num, val)
			}
		}
		// fmt.Println(num)
		//没有nums[i]的切片num
		temp := permute(num)
		for _, val := range temp {
			now := []int{}
			now = append(now, nums[i])
			now = append(now, val...)
			//现在now是一个[]int,追加到back后面
			back = append(back, now)
		}
	}
	return back
}

func main() {
	var n int
	fmt.Scanf("%d\n", &n)
	testSlice := make([]int, n)
	// 标准输入n个不重复的数字
	for i := range testSlice {
		fmt.Scanf("%d", &testSlice[i])
	}
	// fmt.Println(testSlice)
	res := permute(testSlice)
	fmt.Println(res)
	//如果要优化的话，就不能硬生生递归去找了，要去查找下一个，插入啊，换序啊，都可以
}
