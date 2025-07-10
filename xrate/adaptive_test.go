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
 @Time    : 2025/5/8 -- 16:22
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xrate xrate/adaptive_test.go
*/

package xrate

import (
	"testing"
	"time"
)

func TestInitialization(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(100),
		WithMaxRate(1000),
		WithMinRate(10),
		WithWindowSize(60),
		WithBucketSize(100),
	)

	if limiter.Stats().CurrentRate != 100 {
		t.Errorf("期望初始速率为 100，实际为 %f", limiter.Stats().CurrentRate)
	}
	if limiter.Stats().Tokens != 100 {
		t.Errorf("期望初始令牌数为 100，实际为 %f", limiter.Stats().Tokens)
	}
	if limiter.Stats().RequestCount != 0 {
		t.Errorf("期望初始请求数为 0，实际为 %d", limiter.Stats().RequestCount)
	}
	if limiter.Stats().RejectCount != 0 {
		t.Errorf("期望初始拒绝数为 0，实际为 %d", limiter.Stats().RejectCount)
	}
}

func TestBasicRateLimiting(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(10),
		WithBucketSize(10),
	)

	// 前10个请求应该通过
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Errorf("第 %d 个请求被错误限流", i+1)
		}
	}

	// 第11个请求应该被限流
	if limiter.Allow() {
		t.Error("第11个请求应该被限流")
	}
}

func TestTokenRefill(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(10),
		WithBucketSize(10),
	)

	// 消耗所有令牌
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Errorf("第 %d 个请求被错误限流", i+1)
		}
	}

	// 等待1秒，应该补充10个令牌
	time.Sleep(1100 * time.Millisecond)

	// 应该可以通过10个请求
	for i := 0; i < 10; i++ {
		if !limiter.Allow() {
			t.Errorf("补充令牌后第 %d 个请求被错误限流", i+1)
		}
	}
}

func TestRateAdjustment(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(100),
		WithMaxRate(1000),
		WithMinRate(10),
		WithAdjustmentInterval(time.Second),
	)

	// 模拟高负载
	for i := 0; i < 1000; i++ {
		limiter.Allow()
		time.Sleep(time.Millisecond) // 模拟请求处理时间
	}

	// 等待调整
	time.Sleep(1100 * time.Millisecond)

	// 速率应该降低
	if limiter.Stats().CurrentRate >= 100 {
		t.Errorf("期望速率降低，实际为 %f", limiter.Stats().CurrentRate)
	}
}

func TestStats(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(10),
		WithBucketSize(10),
	)

	// 发送一些请求
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}

	// 被限流一些请求
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}

	stats := limiter.Stats()
	if stats.RequestCount != 5 {
		t.Errorf("期望请求数为 5，实际为 %d", stats.RequestCount)
	}
	if stats.RejectCount != 5 {
		t.Errorf("期望拒绝数为 5，实际为 %d", stats.RejectCount)
	}
	if stats.RejectRate != 1.0 {
		t.Errorf("期望拒绝率为 1.0，实际为 %f", stats.RejectRate)
	}
	if stats.AvgLatency < 0 || stats.AvgLatency > 0.1 {
		t.Errorf("期望平均延迟在 0-100ms 之间，实际为 %f", stats.AvgLatency)
	}
}

func TestConcurrentRequests(t *testing.T) {
	limiter := NewAdaptiveRateLimiter(
		WithInitialRate(100),
		WithBucketSize(100),
	)

	results := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			results <- limiter.Allow()
		}()
	}

	// 收集结果
	successCount := 0
	for i := 0; i < 100; i++ {
		if <-results {
			successCount++
		}
	}

	// 验证结果
	if successCount != 100 {
		t.Errorf("期望所有请求都通过，实际通过 %d 个", successCount)
	}
}
