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
 @Time    : 2024/10/12 -- 11:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: metrics.go
*/

package sys

import (
	"github.com/xneogo/Z1ON0101/xstat/xmetric/xprometheus"
)

const namespace = "runtime_resource"

var (
	// cpu usage
	_metricCPUUsage = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "cpu_usage",
		Name:       "current",
		Help:       "cup usage percentage",
		LabelNames: []string{"group", "service"},
	})
	// load avg 1 min
	_metricLoadAvg1min = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "loadavg",
		Name:       "current",
		Help:       "load avg 1 min",
		LabelNames: []string{"group", "service"},
	})
	// memory of process
	_metricMemory = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "memory",
		Name:       "current",
		Help:       "memory of process from runtime memory.sys",
		LabelNames: []string{"group", "service"},
	})
	// goroutine number
	_metricGoroutine = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "goroutine",
		Name:       "current",
		Help:       "goroutine number",
		LabelNames: []string{"group", "service"},
	})
	// heap objects number
	_metricHeapObjects = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "heap_objects",
		Name:       "current",
		Help:       "heap objects",
		LabelNames: []string{"group", "service"},
	})
	// last gc pause
	_metricLastGCPause = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "gc_pause",
		Name:       "last",
		Help:       "last gc pause",
		LabelNames: []string{"group", "service"},
	})
	// heap alloc
	_metricHeapAlloc = xprometheus.NewGauge(&xprometheus.GaugeVecOpts{
		Namespace:  namespace,
		Subsystem:  "heap_alloc",
		Name:       "current",
		Help:       " HeapAlloc is bytes of allocated heap objects.",
		LabelNames: []string{"group", "service"},
	})
)
