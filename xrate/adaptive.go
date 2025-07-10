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
 @Time    : 2025/5/8 -- 15:30
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xrate xrate/adaptive.go
*/

package xrate

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/xneogo/Z1ON0101/xstat"
	"github.com/xneogo/Z1ON0101/xstat/xmetric/xprometheus"
)

// AdaptiveRateLimiter 自适应限流器实现
type AdaptiveRateLimiter struct {
	// 配置参数
	initialRate        float64
	maxRate            float64
	minRate            float64
	windowSize         int64
	bucketSize         float64
	adjustmentInterval time.Duration

	// 运行时状态
	currentRate float64
	tokens      float64
	lastUpdate  time.Time

	// 统计信息
	requestCount int64
	rejectCount  int64
	totalLatency float64

	// 监控指标
	rateGauge        xstat.Gauge
	requestCounter   xstat.Counter
	rejectCounter    xstat.Counter
	latencyHistogram xstat.Histogram

	// 分布式支持
	redisClient *redis.Client
	mu          sync.Mutex

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
}

// NewAdaptiveRateLimiter 创建新的限流器实例
func NewAdaptiveRateLimiter(opts ...Option) XRate {
	// 默认配置
	limiter := &AdaptiveRateLimiter{
		initialRate:        100.0,
		maxRate:            1000.0,
		minRate:            10.0,
		windowSize:         60,
		bucketSize:         100.0,
		adjustmentInterval: 10 * time.Second,
		currentRate:        100.0,
		tokens:             100.0,
		lastUpdate:         time.Now(),
	}

	// 应用自定义配置
	for _, opt := range opts {
		opt(limiter)
	}

	// 初始化监控指标
	limiter.initMetrics()

	// 创建上下文
	limiter.ctx, limiter.cancel = context.WithCancel(context.Background())

	// 启动调整协程
	go limiter.adjustmentLoop()

	return limiter
}

// Option 限流器配置选项
type Option func(*AdaptiveRateLimiter)

// WithInitialRate 设置初始速率
func WithInitialRate(rate float64) Option {
	return func(r *AdaptiveRateLimiter) {
		r.initialRate = rate
		r.currentRate = rate
	}
}

// WithMaxRate 设置最大速率
func WithMaxRate(rate float64) Option {
	return func(r *AdaptiveRateLimiter) {
		r.maxRate = rate
	}
}

// WithMinRate 设置最小速率
func WithMinRate(rate float64) Option {
	return func(r *AdaptiveRateLimiter) {
		r.minRate = rate
	}
}

// WithWindowSize 设置时间窗口大小
func WithWindowSize(size int64) Option {
	return func(r *AdaptiveRateLimiter) {
		r.windowSize = size
	}
}

// WithBucketSize 设置令牌桶大小
func WithBucketSize(size float64) Option {
	return func(r *AdaptiveRateLimiter) {
		r.bucketSize = size
		r.tokens = size
	}
}

// WithAdjustmentInterval 设置调整间隔
func WithAdjustmentInterval(interval time.Duration) Option {
	return func(r *AdaptiveRateLimiter) {
		r.adjustmentInterval = interval
	}
}

// WithRedisClient 设置 Redis 客户端
func WithRedisClient(client *redis.Client) Option {
	return func(r *AdaptiveRateLimiter) {
		r.redisClient = client
	}
}

// initMetrics 初始化监控指标
func (r *AdaptiveRateLimiter) initMetrics() {
	r.rateGauge = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Name: "rate_limiter_current_rate",
		Help: "Current rate limit",
	})

	r.requestCounter = xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Name: "rate_limiter_requests_total",
		Help: "Total requests",
	})

	r.rejectCounter = xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Name: "rate_limiter_rejects_total",
		Help: "Total rejected requests",
	})

	r.latencyHistogram = xprometheus.NewHistogram(&xprometheus.HistogramVecOpts{
		Name:    "rate_limiter_latency_seconds",
		Help:    "Request latency",
		Buckets: prometheus.DefBuckets,
	})
}

// adjustmentLoop 动态调整循环
func (r *AdaptiveRateLimiter) adjustmentLoop() {
	ticker := time.NewTicker(r.adjustmentInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.adjustRate()
		case <-r.ctx.Done():
			return
		}
	}
}

// adjustRate 调整限流速率
func (r *AdaptiveRateLimiter) adjustRate() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 计算系统负载指标
	loadFactor := r.calculateLoadFactor()

	// 计算响应时间指标
	latencyFactor := r.calculateLatencyFactor()

	// 计算错误率指标
	errorFactor := r.calculateErrorFactor()

	// 综合调整因子
	adjustmentFactor := (loadFactor + latencyFactor + errorFactor) / 3

	// 计算新的限流速率
	newRate := r.currentRate * adjustmentFactor

	// 确保在允许范围内
	newRate = math.Max(math.Min(newRate, r.maxRate), r.minRate)

	// 平滑调整
	r.currentRate = r.currentRate*0.7 + newRate*0.3

	// 更新监控指标
	r.rateGauge.Set(r.currentRate)
}

// calculateLoadFactor 计算系统负载因子
func (r *AdaptiveRateLimiter) calculateLoadFactor() float64 {
	if r.requestCount == 0 {
		return 1.0
	}

	// 基于请求量和响应时间计算负载
	avgLatency := r.totalLatency / float64(r.requestCount)
	load := (float64(r.requestCount) / float64(r.windowSize)) * avgLatency

	// 归一化到 [0.5, 1.5] 范围
	return 1.0 + (load-r.currentRate)/(2*r.currentRate)
}

// calculateLatencyFactor 计算响应时间因子
func (r *AdaptiveRateLimiter) calculateLatencyFactor() float64 {
	if r.requestCount == 0 {
		return 1.0
	}

	avgLatency := r.totalLatency / float64(r.requestCount)
	targetLatency := 0.1 // 目标响应时间（秒）

	// 归一化到 [0.5, 1.5] 范围
	return 1.0 + (targetLatency-avgLatency)/(2*targetLatency)
}

// calculateErrorFactor 计算错误率因子
func (r *AdaptiveRateLimiter) calculateErrorFactor() float64 {
	if r.requestCount == 0 {
		return 1.0
	}

	errorRate := float64(r.rejectCount) / float64(r.requestCount)
	targetErrorRate := 0.01 // 目标错误率

	// 归一化到 [0.5, 1.5] 范围
	return 1.0 + (targetErrorRate-errorRate)/(2*targetErrorRate)
}

// Allow 检查是否允许请求通过
func (r *AdaptiveRateLimiter) Allow() bool {
	startTime := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	// 更新令牌桶
	now := time.Now()
	timePassed := now.Sub(r.lastUpdate).Seconds()
	r.lastUpdate = now

	// 添加新令牌
	newTokens := timePassed * r.currentRate
	r.tokens = math.Min(r.bucketSize, r.tokens+newTokens)

	// 检查是否有足够的令牌
	if r.tokens >= 1.0 {
		r.tokens -= 1.0
		r.requestCount++
		r.totalLatency += time.Since(startTime).Seconds()
		r.requestCounter.Inc()
		r.latencyHistogram.Observe(time.Since(startTime).Seconds())
		return true
	} else {
		r.rejectCount++
		r.rejectCounter.Inc()
		return false
	}
}

// Stats 获取限流器统计信息
func (r *AdaptiveRateLimiter) Stats() Stats {
	r.mu.Lock()
	defer r.mu.Unlock()

	return Stats{
		CurrentRate:  r.currentRate,
		RequestCount: r.requestCount,
		RejectCount:  r.rejectCount,
		RejectRate:   float64(r.rejectCount) / math.Max(1, float64(r.requestCount)),
		AvgLatency:   r.totalLatency / math.Max(1, float64(r.requestCount)),
		Tokens:       r.tokens,
	}
}

// Close 关闭限流器
func (r *AdaptiveRateLimiter) Close() {
	r.cancel()
}
