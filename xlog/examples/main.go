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
 @Time    : 2024/10/28 -- 11:48
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: main.go
*/

package main

import (
	"context"
	"github.com/qiguanzhu/infra/nerv/xlog"
	"log"
	"time"
)

// LogFile write log to file+console(debug)
func LogFile(xl *xlog.XLogger) {
	xl.Info(context.TODO(), "info message")
	xl.Infof(context.TODO(), "info message %s", "format")
	xl.Infow(context.TODO(), "info message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.1)

	xl.Warn(context.TODO(), "warn message")
	xl.Warnf(context.TODO(), "warn message %s", "format")
	xl.Warnw(context.TODO(), "warn message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.2)

	xl.Error(context.TODO(), "error message")
	xl.Errorf(context.TODO(), "error message: %s", "format")
	xl.Errorw(context.TODO(), "error message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.3)
}

// LogConsole log to console
func LogConsole(xl *xlog.XLogger) {
	xl.Info(context.TODO(), "info message")
	xl.Infof(context.TODO(), "info message %s", "format")
	xl.Infow(context.TODO(), "info message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.0)

	xl.Warn(context.TODO(), "warn message")
	xl.Warnf(context.TODO(), "warn message %s", "format")
	xl.Warnw(context.TODO(), "warn message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.0)

	xl.Error(context.TODO(), "error message")
	xl.Errorf(context.TODO(), "error message: %s", "format")
	xl.Errorw(context.TODO(), "error message key-value",
		"bkey", false,
		"key1", "value1",
		"intk", 10,
		"fkey", 32.0)
}

func main() {
	// 默认按小时切割, path, filename, maxage, level, format type
	xl, err := xlog.New("./nerv/xlog/examples", "test.log", xlog.InfoLevel, xlog.JSONFormatType, false, 0, false)
	if err != nil {
		log.Fatalf("xlog new failed, %v", err)
	}
	defer xl.Sync()

	xlConsole, _ := xlog.NewConsole(xlog.DebugLevel)
	defer xlConsole.Sync()
	for {
		LogFile(xl)
		// LogConsole(xlConsole)
		time.Sleep(time.Second * 1)
	}
}
