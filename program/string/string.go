package main

import (
	"bytes"
	"fmt"
)

var (
	testData = []string{"Hello", "World!", "qwq"}
	index    = 0
)

func getNextString() (string, bool) {
	if index < len(testData) {
		s := testData[index]
		index++
		return s, true
	}
	return "", false
}

func pop(buf bytes.Buffer) ([]byte, []byte) {
	return buf.Bytes()[:len(buf.Bytes())-1], buf.Bytes()[len(buf.Bytes())-1 : len(buf.Bytes())]
}

func main() {
	//这里试一下用[]byte处理字符串
	//主要是go的字符串是常量
	//差不多是常量（，要么就赋值个新的，应该说是不可变的
	//一般方法就是转成[]byte或者[]rune,前者更普遍，后者。。。我反正这么干的
	//[]byte和string之间的转换比较灵活，所以值得学习
	//以及该程序学习一下string的处理方法，库是strings我记得

	s := "abc"
	by := []byte(s)
	fmt.Println(string(by))
	//这样就转好了

	//然后说，我之前一直有点误解
	//string的类型其实算是byte的只读序列
	//然后[]byte字节切片相当于可变字符串
	//[]rune也是，但是不同的是存的不是字节
	//而是Unicode字符
	//for i:=range str遍历的其实是Unicode字符
	//而不是字节遍历
	//所以有中文字符的时候是跳着的索引

	s1 := s[1:] //这样可以字符串直接截取,仍然是字符串
	fmt.Printf("%T %s\n", s1, s1)
	//字符串在内存中的结构其实是指向底层字节序列的指针加上长度

	//追加字符串，效率高一些的方法
	var buf bytes.Buffer //长度可变的bytes buffer，提供了Read和Write方法
	//因为读写长度位置的bytes最好用buffer
	//也可以这么定义
	//var r *bytes.Buffer = new(bytes.Buffer)
	//以及 func NewBuffer(buf []byte) *Buffer
	// 创建一个 Buffer 对象并且用 buf 初始化好
	//NewBuffer 最好用在从 buf 读取的时候使用。

	//接下来串联字符串
	//buffer.WriteString(s)可以将字符串追加到后面
	//在学接口的时候总是见到这些，也算熟人了
	for {
		if s, ok := getNextString(); ok {
			//method getNextString() not shown here
			buf.WriteString(s)
		} else {
			break
		}
	}
	fmt.Println(buf.String())
	//如果是直接是buf输出，会是
	//{[72 101 108 108 111 87 111 114 108 100 33 113 119 113] 0 0}
	//最后两个不是[]byte数据，而是buffer结构体的一部分
	//是off和lastRead
	//读偏移和上次操作类型

	var by1, by2 []byte
	by1 = append(by1, "hello"...)
	//等价于一个字符一个字符追加
	by2 = append(by2, ",world!"...)
	by1 = append(by1, by2...)
	fmt.Println(string(by1))
	//append第二个接收的是元素，直接用切片并不行
	//故用...语法把切片元素展开

	//然后写了bytes.Buffer个返回前n个[]byte和剩余的[]byte
	by3, by4 := pop(buf)
	fmt.Println(string(by3), string(by4))
}
