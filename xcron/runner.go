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
 @Time    : 2025/4/15 -- 15:27
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcron xcron/runner.go
*/

package xcron

import "context"

type CronRunner struct {
	jobs map[string]CronUnit
	run  func(ctx context.Context) func(remark string, unit CronUnit)
}

func (r *CronRunner) Run(ctx context.Context) {
	for remark, job := range r.jobs {
		r.run(ctx)(remark, job)
	}
}

func (r *CronRunner) Register(remark string, unit CronUnit) CronRunnerProxy {
	if _, ok := r.jobs[remark]; !ok {
		r.jobs[remark] = unit
	}
	return r
}

func (r *CronRunner) ForceRegister(remark string, unit CronUnit) CronRunnerProxy {
	r.jobs[remark] = unit
	return r
}
