package main

import (
	"fmt"
)

// 责任链
type ApprovalRequest struct {
	Name   string
	Amount float64
	// 5000 10000 100000
}
type Handler interface {
	HandleRequest(req ApprovalRequest)
	SetNextHandler(next Handler) Handler //这里返回Handler方便链式调用
}

type BaseHandler struct {
	next Handler
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (b *BaseHandler) SetNextHandler(next Handler) Handler {
	b.next = next
	return next
}

type LimitHandler struct {
	*BaseHandler
	limit     float64
	onSuccess []func(req ApprovalRequest) //可选的自定义处理逻辑，按顺序执行
	name      string                      //用于打印
}

func NewLimitHandler(name string, limit float64, onSuccess ...func(req ApprovalRequest)) *LimitHandler {
	var fns []func(req ApprovalRequest)
	for _, f := range onSuccess {
		if f != nil {
			fns = append(fns, f)
		}
	}
	return &LimitHandler{
		BaseHandler: NewBaseHandler(),
		name:        name,
		limit:       limit,
		onSuccess:   fns,
	}
}

func (l *LimitHandler) HandleRequest(req ApprovalRequest) {
	if req.Amount <= l.limit {
		if len(l.onSuccess) > 0 {
			for _, f := range l.onSuccess {
				f(req)
			}
		} else {
			fmt.Printf("%s 审批通过: %s 申请金额 %.2f\n", l.name, req.Name, req.Amount)
		}
	} else if l.next != nil {
		fmt.Printf("%s 无法审批 %.2f, 转交上级\n", l.name, req.Amount)
		l.next.HandleRequest(req)
	} else {
		fmt.Println("审批流程结束，无人处理")
	}
}

func main() {
	manager := NewLimitHandler("组长", 5000, nil)
	//如果没判断OnSuccess的元素是不是nil
	//这里就不要传nil,nil是nil函数，会panic
	deptHead := NewLimitHandler("组长", 10000)
	ceo := NewLimitHandler("CEO", 100000, func(req ApprovalRequest) {
		fmt.Printf("%s 获得预算为: %.2f", req.Name, req.Amount)
	})

	manager.SetNextHandler(deptHead).SetNextHandler(ceo)

	req1 := ApprovalRequest{Name: "张三", Amount: 3000}   // 组长处理
	req2 := ApprovalRequest{Name: "李四", Amount: 8000}   // 经理处理
	req3 := ApprovalRequest{Name: "王五", Amount: 50000}  // CEO处理
	req4 := ApprovalRequest{Name: "赵六", Amount: 200000} // 无人能处理

	fmt.Println("=== 请求1 ===")
	manager.HandleRequest(req1)

	fmt.Println("\n=== 请求2 ===")
	manager.HandleRequest(req2)

	fmt.Println("\n=== 请求3 ===")
	manager.HandleRequest(req3)

	fmt.Println("\n=== 请求4 ===")
	manager.HandleRequest(req4)
}
