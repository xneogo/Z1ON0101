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
 @Time    : 2025/7/3 -- 10:35
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: zion zion/replace.go
*/

package zion

import (
	"bytes"
	"context"
	"text/template"

	"github.com/xneogo/Z1ON0101/xlog"
)

func Replace(ctx context.Context, format string, params map[string]string) string {
	fun := "zion.Replace"
	b := &bytes.Buffer{}
	defer func() {
		if r := recover(); r != nil {
			xlog.Errorf(ctx, "%s fail err: template.Panic. format:%s, params:%+v", fun, format, params)
		}
	}()
	_ = template.Must(template.New("").Parse(format)).Execute(b, params)
	return b.String()
}
