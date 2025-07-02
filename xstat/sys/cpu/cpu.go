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
 @Description: cpu.go
*/

package cpu

import (
	"sync/atomic"
	"time"
)

const (
	interval = time.Millisecond * 500
)

var (
	stats CPU
	usage uint64
)

// CPU is cpu stat usage.
type CPU interface {
	Usage() (u uint64, e error)
	Info() Info
}

// Stat cpu stat.
type Stat struct {
	Usage uint64 // cpu use ratio.
}

// Info cpu info.
type Info struct {
	Frequency uint64
	Quota     float64
}

// ReadStat read cpu stat.
func ReadStat(stat *Stat) {
	stat.Usage = atomic.LoadUint64(&usage)
}

// GetInfo get cpu info.
func GetInfo() Info {
	return stats.Info()
}
