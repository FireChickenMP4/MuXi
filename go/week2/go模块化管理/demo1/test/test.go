package test

import (
	"demo1/model"
	"fmt"
)

func ShowScores() {
	Stu := model.Student{
		Name:   "小明",
		Age:    18,
		Scores: []int{90, 80, 70},
	}

	fmt.Printf("姓名：%s，年龄：%d，成绩：%v\n", Stu.Name, Stu.Age, Stu.Scores)
}
