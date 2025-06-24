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
 @Time    : 2025/4/15 -- 15:26
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcron xcron/global.go
*/

package xcron

import (
	"context"
	"errors"
)

func NewRunCronOnlyOnce(leaderSelection func() bool) (CronRunnerProxy, error) {
	if leaderSelection == nil {
		return nil, errors.New("[xcron] leaderSelection is nil")
	}
	if !leaderSelection() {
		return nil, errors.New("[xcron] current instance elect for leader failed")
	}
	return &CronRunner{
		jobs: make(map[string]CronUnit),
		run: func(_ctx context.Context) func(remark string, unit CronUnit) {
			return func(remark string, unit CronUnit) {
				scheduler := newCronScheduler(_ctx)
				xlog.Infof(_ctx, "%s loading to scheduler", remark)
				_ = scheduler.Register(unit)
			}
		},
	}, nil
}

func NewRunCronEveryInstance() (CronRunnerProxy, error) {
	return &CronRunner{
		jobs: make(map[string]CronUnit),
		run: func(_ctx context.Context) func(remark string, unit CronUnit) {
			return func(remark string, unit CronUnit) {
				scheduler := newCronScheduler(_ctx)
				xlog.Infof(_ctx, "%s loading to scheduler", remark)
				_ = scheduler.Register(unit)
			}
		},
	}, nil
}
