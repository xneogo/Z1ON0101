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
 @Time    : 2025/4/15 -- 15:28
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcron xcron/scheduler.go
*/

package xcron

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xneogo/saferun"
)

// CronScheduler
// 本地 cron 定时调度器，可以注册多个 unit，按固定/非固定时间执行
// 配合 自行实现的选主函数 实现全局单例定时任务
// 具体使用示例见 /prototype/cron/xx.go
type CronScheduler struct {
	sync.WaitGroup
	ctx     context.Context
	isClose int64
}

func newCronScheduler(ctx context.Context) CronSchedulerModel {
	return &CronScheduler{
		ctx: ctx,
	}
}

func (w *CronScheduler) Register(c CronUnit) (err error) {
	// 先初始化
	if err = c.Init(w.ctx); err != nil {
		// w.logger.Errorf(w.ctx, "init crond %s fail %s", c.Name(), err)
		return
	}

	checkStopInv := time.Tick(time.Second)

	w.Add(1)
	go func(c CronUnit) {
		defer func() {
			if err := saferun.DumpStack(recover()); err != nil {
				// xlog.Errorf(w.ctx, "crond panic %s: %s", c.Name(), err)
			}
			w.Done()
		}()

		for {
			func(c CronUnit) {
				defer func() {
					err := saferun.DumpStack(recover())
					if err != nil {
						// xlog.Errorf(w.ctx, "crond panic %s: %v", c.Name(), err)
						time.Sleep(1 * time.Second)
					}
				}()

				c.Do()
			}(c)

			ticker := c.GetTicker()

		WaitNext:
			for {
				select {
				case <-ticker:
					break WaitNext
				case <-checkStopInv:
					if w.IsClose() {
						// xlog.Infof(w.ctx, "crond %s is stopping", c.Name())
						return
					}
				}
			}
		}

	}(c)
	return
}

func (w *CronScheduler) Close() {
	atomic.StoreInt64(&w.isClose, 1)
}

func (w *CronScheduler) IsClose() bool {
	return atomic.LoadInt64(&w.isClose) == 1
}
