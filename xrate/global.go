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
 @Time    : 2025/5/8 -- 15:27
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xrate xrate/global.go
*/

package xrate

// Stats 获取限流器统计信息
type Stats struct {
	CurrentRate  float64
	RequestCount int64
	RejectCount  int64
	RejectRate   float64
	AvgLatency   float64
	Tokens       float64
}

type XRate interface {
	Allow() bool
	Close()
	Stats() Stats
}
