package main

import (
	"fmt"
)

type intlinklist struct {
	nextaddress *intlinklist
	value       int
}

func main() {
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
}
