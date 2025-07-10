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
 @Time    : 2024/11/5 -- 18:10
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: metrics.go
*/

package xredis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/xneogo/Z1ON0101/xstat/xmetric/xprometheus"
	"github.com/xneogo/Z1ON0101/xtrace"
	"github.com/xneogo/extensions/xtime"
)

const (
	namespace = "zion"
	subsystem = "redis_requests"
)

var (
	buckets = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500}

	// caller_endpoint 调用方方法
	// cluster 集群名
	_metricRequestDuration = xprometheus.NewHistogram(&xprometheus.HistogramVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "duration_ms",
		Help:       "redisext requests duration(ms).",
		LabelNames: []string{"namespace", "command", "cluster", xprometheus.LabelCallerEndpoint},
		Buckets:    buckets,
	})

	_metricReqErr = xprometheus.NewCounter(&xprometheus.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       "err_total",
		Help:       "redisext requests error total",
		LabelNames: []string{"namespace", "command", "cluster", xprometheus.LabelCallerEndpoint},
	})
)

func statReqDuration(ctx context.Context, namespace, command string, cluster string, durationMS int64) {
	_metricRequestDuration.With("namespace", namespace, "command", command, "cluster", cluster,
		xprometheus.LabelCallerEndpoint, xtrace.ExtractSpanOperationName(ctx)).Observe(float64(durationMS))
}

func statReqErr(ctx context.Context, namespace, command string, cluster string, err error) {
	if err != nil && !errors.Is(err, redis.Nil) {
		_metricReqErr.With("namespace", namespace, "command", command, "cluster", cluster,
			xprometheus.LabelCallerEndpoint, xtrace.ExtractSpanOperationName(ctx)).Inc()
	}
	return
}

type RedisMetric struct {
	Command   string
	Namespace string
	Cluster   string
}

func WrapMetric(ctx context.Context, redisMetric RedisMetric, executor func() error) error {
	span, ctx := xtrace.StartSpanFromContext(ctx, redisMetric.Command)
	st := xtime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(ctx, redisMetric.Namespace, redisMetric.Command, redisMetric.Cluster, st.Millisecond())
	}()
	err := executor()
	statReqErr(ctx, redisMetric.Namespace, redisMetric.Command, redisMetric.Cluster, err)
	return err
}
