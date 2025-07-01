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
 @Time    : 2024/10/28 -- 11:43
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: stat_logger.go
*/

package xlog

import (
	"context"
)

// InitStatLog init stat log
func InitStatLog(logRoot, fileName string, maxRetainDays int, compress bool) (err error) {
	statLogger, err = New(logRoot, fileName, InfoLevel, JSONFormatType, false, maxRetainDays, compress)
	if err != nil {
		return err
	}
	return nil
}

// StatInfow 统计日志，日志级别默认为Info
func StatInfow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	statLogger.Infow(ctx, msg, keysAndValues...)
}

// SetStatLogService set service in stat logger
func SetStatLogService(service string) {
	statLogger.SetService(service)
}

// SetStatLogSkip set skip of app logger
func SetStatLogSkip(skip int) {
	statLogger.skip = skip
}

// StatLogSync stat log sync
func StatLogSync() {
	statLogger.Sync()
}

// StatLogWithMap 统计日志，通过传入map的方式
func StatLogWithMap(ctx context.Context, msg string, kvmap map[string]interface{}) {
	statLogger.Infom(ctx, msg, kvmap)
}
