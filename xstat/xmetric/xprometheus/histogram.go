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
 @Time    : 2024/10/12 -- 11:44
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: histogram.go
*/

package xprometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xneogo/Z1ON0101/xstat"
)

// HistogramVecOpts is histogram vector opts.
type HistogramVecOpts struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	LabelNames []string
	Buckets    []float64
}

// Histogram prom histogram collection.
type promHistogramVec struct {
	hv  *prometheus.HistogramVec
	lvs xstat.LabelValues
}

// NewHistogram constructs and registers a Prometheus HistogramVec,
// and returns a usable Histogram object.
func NewHistogram(cfg *HistogramVecOpts) xstat.Histogram {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
			Buckets:   cfg.Buckets,
		}, cfg.LabelNames)
	prometheus.MustRegister(vec)
	return &promHistogramVec{
		hv: vec,
	}
}

// With append k-v pairs to histogram lvs
func (h *promHistogramVec) With(labelValues ...string) xstat.Histogram {
	return &promHistogramVec{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}

// Timing adds a single observation to the histogram.
func (h *promHistogramVec) Observe(v float64) {
	if err := h.lvs.Check(); err != nil {
		fmt.Printf("histogram label value invalid:%s\n", err.Error())
		return
	}
	h.hv.With(makeLabels(h.lvs...)).Observe(v)
}
