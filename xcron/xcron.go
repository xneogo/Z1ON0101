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
 @Time    : 2025/4/15 -- 15:29
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcron xcron/xcron.go
*/

package xcron

import (
	"context"
	"reflect"

	"github.com/robfig/cron/v3"
)

type XCron struct {
	C *cron.Cron
}

// New new a XCron
func New() *XCron {
	return &XCron{
		C: cron.New(),
	}
}

func (x XCron) Register(rate string, f func()) error {
	ctx := context.Background()
	fnName := func() string {
		fnVal := reflect.ValueOf(f)
		if fnVal.Type().Kind() != reflect.Func {
			return "not a func"
		}
		return fnVal.Type().Name()
	}()
	id, err := x.C.AddFunc(rate, f)
	if err != nil {
		xlog.Errorf(ctx, "Xcron cronJob.AddFunc %s has an error. e: %+v", fnName, err)
	}
	xlog.Infof(ctx, "Xcron cronJob.AddFunc %s  id: %d successfully added", fnName, id)
}

func (x XCron) Run() {
	x.C.Run()
}
