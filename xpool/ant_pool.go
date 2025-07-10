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
 @Time    : 2024/11/4 -- 17:43
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: ant_pool.go
*/

package xpool

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"github.com/xneogo/Z1ON0101/xlog"
)

// Pool alias of ants.Pool
type Pool struct {
	ap *ants.Pool
}

// NewWorkerPoolWithOptions 创建Worker池, 支持传入ants Option.
// 注意, logger不允许变更, 必须使用xlog, 因此WithLogger()会被覆盖
func NewWorkerPoolWithOptions(size int, options ...ants.Option) (*Pool, error) {
	options = append(options, ants.WithLogger(&WorkerLogger{}))
	ap, err := ants.NewPool(size, options...)
	if err != nil {
		return nil, err
	}
	return &Pool{ap: ap}, nil
}

// NewWorkerPool constructor of Pool
func NewWorkerPool(size int) (*Pool, error) {
	return NewWorkerPoolWithOptions(size)
}

// Tune 调整Worker池容量
func (p *Pool) Tune(size int) {
	p.ap.Tune(size)
}

// Release close pool and release resources
func (p *Pool) Release() {
	p.ap.Release()
}

// Submit submit a task
func (p *Pool) Submit(task func()) error {
	return p.ap.Submit(task)
}

// Running return goroutines of runnning
func (p *Pool) Running() int {
	return p.ap.Running()
}

// WorkerLogger log handler
type WorkerLogger struct {
}

// Printf implements ants.Logger
func (WorkerLogger) Printf(format string, args ...interface{}) {
	xlog.Infof(context.Background(), format, args...)
}
