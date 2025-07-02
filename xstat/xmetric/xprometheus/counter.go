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
 @Time    : 2024/10/12 -- 11:42
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: counter.go
*/

package xprometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xneogo/Z1ON0101/xstat"
)

// CounterVecOpts is an alias of VectorOpts.
type CounterVecOpts xstat.VectorOpts

// counterVec counter vec.
type promCounterVec struct {
	cv  *prometheus.CounterVec
	lvs xstat.LabelValues
}

// NewCounter constructs and register a Prometheus CounterVec,
// and return a usable Counter object.
func NewCounter(cfg *CounterVecOpts) xstat.Counter {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.LabelNames)
	prometheus.MustRegister(vec)
	return &promCounterVec{
		cv: vec,
	}
}

// Inc increments the counter by 1.
func (c *promCounterVec) Inc() {
	if err := c.lvs.Check(); err != nil {
		fmt.Printf("counter label value invalid:%s\n", err.Error())
		return
	}
	c.cv.With(makeLabels(c.lvs...)).Inc()
}

// Add adds the given value to the counter. It panics if the value is < 0
func (c *promCounterVec) Add(delta float64) {
	if err := c.lvs.Check(); err != nil {
		fmt.Printf("counter label value invalid:%s\n", err.Error())
		return
	}
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}

// With implements Counter.
func (c *promCounterVec) With(labelValues ...string) xstat.Counter {
	return &promCounterVec{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}
