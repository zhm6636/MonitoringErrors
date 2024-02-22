package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu           sync.Mutex
	cond         *sync.Cond
	monitorMap   map[string]interface{}
	notification chan string
)

func main() {
	cond = sync.NewCond(&mu)
	monitorMap = make(map[string]interface{})
	notification = make(chan string)

	go monitorVariable("变量1")
	go monitorVariable("变量2")

	// 模拟一些操作，更新变量的值
	go updateVariable("变量1", "值1")
	go updateVariable("变量2", 42)

	// 主goroutine等待一段时间，然后关闭程序
	time.Sleep(10 * time.Second)
}

func monitorVariable(variableName string) {
	for {
		mu.Lock()
		for _, value := range monitorMap {
			if value != nil {
				// 发送信号，表示有值
				notification <- variableName
				break
			}
		}
		cond.Wait() // 等待下一次监控
		mu.Unlock()

		// 处理相关逻辑
		select {
		case variable := <-notification:
			fmt.Printf("变量 %s 具有值\n", variable)
			// 处理相关逻辑
		}
	}
}

func updateVariable(variableName string, value interface{}) {
	time.Sleep(2 * time.Second)
	mu.Lock()
	monitorMap[variableName] = value
	cond.Signal() // 通知等待的监控goroutine
	mu.Unlock()
}
