package main

import (
	"fmt"
)

func main() {
	s := make([]byte, 10)
	fmt.Println(len(s), cap(s)) // 10 10
	for i := range s {
		s[i] = byte(i + 1) //第几个元素
	}
	s = s[2:4]                  //索引起始为0
	fmt.Println(len(s), cap(s)) // 2 3 因为2:4不包括4，经典end不是最后一个
	//fmt.Println(s[3]) 这样会panic，因为len是2
	//但是这里3其实是从2切到最后，所以cap是3
	//总结一下就是
	//cap切片=底层数组长度 - 切片起始索引 因为包括起始索引那个元素，实际上是+1-1
	//len=后值-前值 s[a:b] 就是b-a
	s = s[2:8] //这里相当于切底层数组的 4:10 再大会panic
	//但是 索引是基于上面s[2:4]
	fmt.Println(len(s), cap(s)) // 6 6
	for i := range s {
		fmt.Println(s[i])
	} //而后面的数是基于底层数组的
	s = s[2:2]                  // basic[6:6] cap是4 len是0,因为重合
	fmt.Println(len(s), cap(s)) // 0 4
	//切片没法再访问前面的数了
	fmt.Println(len(s), cap(s))
	str := "hello，世界"
	for i := range str {
		fmt.Println(i, str[i])
	}
	fmt.Println(len(str))
	//这里的话应该想说的是ascii码的事情，中文符号和字都不是
	//其实是UTF8编码的事情，这一块好像还挺复杂
	//只占用一个字节，所以读取次数和长度并不一样，并且跳着蹦着的
	//01234是前五个英文字母
	//5到7是中文逗号
	//8到10 11到13分别是世界
}
