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
 @Time    : 2024/10/12 -- 11:10
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: metrics.go
*/

package xmanager

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	xprom "github.com/xneogo/Z1ON0101/xstat/xmetric/xprometheus"
	"github.com/xneogo/Z1ON0101/xtrace"
)

const namespace = "xsql"

var (
	msBuckets     = []float64{1, 3, 5, 10, 25, 50, 100, 200, 300, 500, 1000, 3000, 5000, 10000, 15000}
	_metricReqDur = xprom.NewHistogram(&xprom.HistogramVecOpts{
		Namespace:  namespace,
		Subsystem:  "requests",
		Name:       "duration_ms",
		Help:       "mysql client requests duration(ms).",
		Buckets:    msBuckets,
		LabelNames: []string{"cluster", "table", "command", xprom.LabelCallerService, xprom.LabelCallerEndpoint},
	})

	_metricReqErrTotal = xprom.NewCounter(&xprom.CounterVecOpts{
		Namespace:  namespace,
		Subsystem:  "requests",
		Name:       "err_total",
		Help:       "mysql client err requests total.",
		LabelNames: []string{"cluster", "table", "command", xprom.LabelCallerService, xprom.LabelCallerEndpoint},
	})

	_metricConnTotal = xprom.NewGauge(&xprom.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "connections",
		Name:       "total",
		Help:       "mysql client connections total count.",
		LabelNames: []string{"dbname"},
	})

	_metricConnInUse = xprom.NewGauge(&xprom.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "connections",
		Name:       "in_use",
		Help:       "mysql client connections in use.",
		LabelNames: []string{"dbname"},
	})

	_metricConnIdle = xprom.NewGauge(&xprom.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "connections",
		Name:       "idle",
		Help:       "mysql client connections idle.",
		LabelNames: []string{"dbname"},
	})
)

func statMetricReqDur(ctx context.Context, cluster, table, command string, durationMS int64) {
	_metricReqDur.With("cluster", cluster, "table", table, "command", command,
		xprom.LabelCallerService, xtrace.ServiceName(), xprom.LabelCallerEndpoint, xtrace.ExtractSpanOperationName(ctx)).
		Observe(float64(durationMS))
}

func statMetricReqErrTotal(ctx context.Context, cluster, table, command string, err error) {
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		_metricReqErrTotal.With("cluster", cluster, "table", table, "command", command,
			xprom.LabelCallerService, xtrace.ServiceName(), xprom.LabelCallerEndpoint, xtrace.ExtractSpanOperationName(ctx)).Inc()
	}
}
