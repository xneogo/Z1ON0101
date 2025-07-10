/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2025/5/8 -- 16:27
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: examples xrate/examples/main.go
*/

package main

import (
	"fmt"
	"github.com/xneogo/Z1ON0101/xrate"
	"math/rand"
	"time"
)

func simulateRequest(limiter xrate.XRate, requestID int) {
	if limiter.Allow() {
		// 模拟请求处理时间
		time.Sleep(time.Duration(rand.Float64()*90+100) * time.Millisecond)
		fmt.Printf("请求 %d 处理成功\n", requestID)
	} else {
		fmt.Printf("请求 %d 被限流\n", requestID)
	}
}

func main() {
	// 创建限流器实例
	limiter := xrate.NewAdaptiveRateLimiter(
		xrate.WithInitialRate(5),                    // 初始每秒10个请求
		xrate.WithMaxRate(10),                       // 最大每秒100个请求
		xrate.WithMinRate(5),                        // 最小每秒5个请求
		xrate.WithWindowSize(60),                    // 60秒时间窗口
		xrate.WithBucketSize(10),                    // 令牌桶大小
		xrate.WithAdjustmentInterval(5*time.Second), // 每5秒调整一次
	)
	defer limiter.Close()

	// 模拟突发流量
	fmt.Println("模拟突发流量...")
	for i := 0; i < 50; i++ {
		simulateRequest(limiter, i)
		time.Sleep(10 * time.Millisecond) // 快速发送请求
	}

	// 等待一段时间
	fmt.Println("\n等待5秒...")
	time.Sleep(5 * time.Second)

	// 打印统计信息
	stats := limiter.Stats()
	fmt.Println("\n限流器统计信息:")
	fmt.Printf("当前限流速率: %.2f 请求/秒\n", stats.CurrentRate)
	fmt.Printf("总请求数: %d\n", stats.RequestCount)
	fmt.Printf("限流请求数: %d\n", stats.RejectCount)
	fmt.Printf("限流率: %.2f%%\n", stats.RejectRate*100)
	fmt.Printf("平均延迟: %.2fms\n", stats.AvgLatency*1000)
	fmt.Printf("当前令牌数: %.2f\n", stats.Tokens)
}
