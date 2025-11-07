package main

import (
	"fmt"
)

func main() {
	var v interface{} = "qwqawa"
	if str, ok := v.(string); ok {
		fmt.Println("Is string:", str)
	}
	if val, ok := v.([]rune); ok {
		for i := range val {
			fmt.Print(val[i])
		}
		fmt.Print('\n')
	}
	//any类型 相当于空接口
	//然后我们输出的时候可用类型断言来检测其动态类型
}
