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
 @Time    : 2024/10/12 -- 15:14
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: psutilCPU.go
*/

package cpu

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type PsutilCPU struct {
	interval time.Duration
}

func NewPsutilCPU(interval time.Duration) (cpu *PsutilCPU, err error) {
	cpu = &PsutilCPU{interval: interval}
	_, err = cpu.Usage()
	if err != nil {
		return
	}
	return
}

func (ps *PsutilCPU) Usage() (u uint64, err error) {
	var percents []float64
	percents, err = cpu.Percent(ps.interval, false)
	if err == nil {
		u = uint64(percents[0] * 10)
	}
	return
}

func (ps *PsutilCPU) Info() (info Info) {
	stats, err := cpu.Info()
	if err != nil {
		return
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		return
	}
	info = Info{
		Frequency: uint64(stats[0].Mhz),
		Quota:     float64(cores),
	}
	return
}
