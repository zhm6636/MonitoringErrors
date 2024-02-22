package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	variables map[string]interface{} // 存储监控的变量
	mu        sync.Mutex
)

func main() {
	variables = make(map[string]interface{})

	// 使用context创建一个可以取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitorVariable(ctx, cancel, "例子变量")

	// 模拟在一段时间后给变量赋值
	time.Sleep(2 * time.Second)
	setVariable("例子变量", "示例值")

	// 模拟在一段时间后发生错误，取消监控goroutine
	time.Sleep(2 * time.Second)
	cancel()

	// 主goroutine等待一段时间，然后关闭程序
	time.Sleep(5 * time.Second)
}

func monitorVariable(ctx context.Context, cancel context.CancelFunc, variableName string) {
	for {
		select {
		case <-ctx.Done():
			// 上下文被取消，结束监控goroutine
			fmt.Println("监控 goroutine 已取消")
			return
		default:
			mu.Lock()
			value := variables[variableName]
			mu.Unlock()

			if value != nil {
				fmt.Printf("变量 %s 有值: %v\n", variableName, value)
				// 在这里可以根据变量值进行相关逻辑处理
				// 如果发现错误，取消上下文
				if isError(value) {
					fmt.Println("检测到错误，取消监视器 goroutine")
					cancel()
					return
				}
			}

			// 模拟定时检查变量的间隔
			time.Sleep(1 * time.Second)
		}
	}
}

func setVariable(variableName string, value interface{}) {
	mu.Lock()
	variables[variableName] = value
	mu.Unlock()
}

func isError(value interface{}) bool {
	// 在这里添加检测错误的逻辑
	// 例如，当某个特定值出现时认为是错误
	return value == "error"
}
