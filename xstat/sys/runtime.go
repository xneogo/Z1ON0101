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
 @Time    : 2024/10/12 -- 15:02
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: runtime.go
*/

package sys

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xneogo/Z1ON0101/xstat/sys/cpu"
	"github.com/xneogo/extensions/xfile"
)

var gRSMonitor *ResourceMonitor

// ResourceMonitor calc resources of instance in one service
type ResourceMonitor struct {
	Group               string // 组名
	Service             string // 服务名
	Instance            string // 实例名称,可能来源于k8s分配
	Hostname            string // 主机名称
	CPUDetector         cpu.CPU
	CPUUsage            uint64
	LoadAvg1Min         float64
	LastInspectUnixNano int64
}

// Init init calc runtime resource
func Init(group, service, instance string) {
	var cpuDetector cpu.CPU
	var err error
	cpuDetector, err = cpu.NewCgroupCPU()
	if err != nil {
		cpuDetector, err = cpu.NewPsutilCPU(0)
		fmt.Printf("cgroup cpu init failed(%v),switch to psutil cpu\n", err)
		if err != nil {
			fmt.Printf("psutil init failed(%v)\n", err)
			return
		}
	}
	gRSMonitor = NewResourceMonitor(group, service, instance, cpuDetector)
	go func() {
		for {
			work()
		}
	}()
}

func work() {
	defer func() {
		recover()
	}()
	time.Sleep(10 * time.Second)
	_ = gRSMonitor.calcCPUUsagePercentage()
	_ = gRSMonitor.calcLoadAvg()
	gRSMonitor.record()
}

// NewResourceMonitor  new monitor
func NewResourceMonitor(group, service, instance string, cpuDetector cpu.CPU) *ResourceMonitor {
	rs := &ResourceMonitor{Group: group, Service: service, Instance: instance, CPUDetector: cpuDetector}
	hostname, _ := os.Hostname()
	rs.Hostname = strings.Replace(hostname, ".", "_", -1)
	rs.LastInspectUnixNano = time.Now().UnixNano()
	return rs
}

func (p *ResourceMonitor) getCurrentCPUUsagePercentage() float64 {
	return float64(p.CPUUsage)
}

func (p *ResourceMonitor) getLoadAvg1Min() float64 {
	return p.LoadAvg1Min
}

func (p *ResourceMonitor) calcLoadAvg() error {
	loadAvg, err := xfile.ReadAll("/proc/loadavg")
	if err != nil {
		return err
	}
	L := strings.Fields(strings.TrimSpace(string(loadAvg)))
	if p.LoadAvg1Min, err = strconv.ParseFloat(L[0], 64); err != nil {
		return err
	}
	return nil
}

func (p *ResourceMonitor) calcCPUUsagePercentage() error {
	u, err := p.CPUDetector.Usage()
	if err == nil && u != 0 {
		p.CPUUsage = u
	}
	return err
}

// GetCurrentCPUUsagePercentage ...
func GetCurrentCPUUsagePercentage() float64 {
	if gRSMonitor == nil {
		return 0.0
	}
	return gRSMonitor.getCurrentCPUUsagePercentage()
}

// GetMemStat ...
func GetMemStat() runtime.MemStats {
	if gRSMonitor == nil {
	}
	return gRSMonitor.getMemStat()
}

// GetNumGoroutine ...
func GetNumGoroutine() int {
	if gRSMonitor == nil {
	}
	return gRSMonitor.getNumOfGoroutine()
}

func (p *ResourceMonitor) getNumOfGoroutine() int {
	return runtime.NumGoroutine()
}

func (p *ResourceMonitor) getMemStat() runtime.MemStats {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	return stat
}

func (p *ResourceMonitor) record() {
	memStat := p.getMemStat()

	_metricCPUUsage.With("group", p.Group, "service", p.Service).Set(p.getCurrentCPUUsagePercentage())
	_metricLoadAvg1min.With("group", p.Group, "service", p.Service).Set(p.getLoadAvg1Min() * 100)
	_metricMemory.With("group", p.Group, "service", p.Service).Set(float64(memStat.Sys))
	_metricGoroutine.With("group", p.Group, "service", p.Service).Set(float64(p.getNumOfGoroutine()))
	_metricHeapObjects.With("group", p.Group, "service", p.Service).Set(float64(memStat.HeapObjects))
	_metricLastGCPause.With("group", p.Group, "service", p.Service).Set(float64(memStat.PauseNs[(memStat.NumGC-1)%uint32(len(memStat.PauseNs))]))
	_metricHeapAlloc.With("group", p.Group, "service", p.Service).Set(float64(memStat.HeapAlloc))
}
