## 简介

xlog是基于[zap](https://github.com/uber-go/zap)封装的内部库，zap是uber开源的异步、高效、结构化的日志库。我们在zap的基础上提供了公司内部格式需要的encoder、易用性和安全性。

## 直接使用

目前xlog内包含应用日志、统计日志两个内置类型，可以直接使用。

1. 日志打印风格

目前支持三种打印风格，具体如下

* fmt.Sprint方式

方法包括Debug(ctx context.Context, args ...interface{})、Info(ctx context.Context, args ...interface{})等方法，示例:

`xlog.Debug(ctx, "test debug")`

* fmt.Sprintf方式

方法包括Debugf(ctx context.Context, template string, args ...interface{})、Infof(ctx context.Context, template string, args ...interface{})等方法，示例:

`xlog.Debugf(ctx, "test user: %d", uid)`

* key-value方式

方法包括Debugw(ctx context.Context, msg string, keysAndValues ...interface{})、Infow(ctx context.Context, msg string, keysAndValues ...interface{})等方法，示例:

`xlog.Debugw(ctx, "test debug kv", "key1", "value1", "key2", 10, "key3", 2.5)`

## 自定义使用
### 自定义说明
当需要主动调用xlog创建日志句柄时，需要关注该部分使用说明，伴鱼内部微服务都直接使用即可，不需要自己创建句柄。
1. 初始化日志句柄

```go
// logDir: 日志目录.
// fileName: 日志文件名称，如命名为: test.log
// level: 日志级别，可以直接引用xlog.DebugLevel/xlog.InfoLevel/xlog.WarnLevel/xlog.ErrorLevel，级别依次升高，低于level级别的日志不会打印.
// formatType: 格式类型，目前支持三种，xlog.JSONFormatType、xlog.AccessLogFormatType和xlog.PlainTextFormatType，其中accesslog使用AccessLogFormatType类型，应用日志建议使用JSONFormatType类型，PlainTextFormatType为普通文本类型。
// debug: 标识是否为调试状态，调试状态下日志会打印到console.
func New(logDir, fileName string, level Level, 
formatType FormatType, debug bool) (logger *XLogger, err error) {
    ...
    ...
}
```
### 代码示例

* 可运行示例

[xlog/examples/main.go](xlog/examples/main.go)


* 代码示例

```golang
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
	xl, err := xlog.New("./", "test.log", xlog.InfoLevel, xlog.JSONFormatType, false, 7, false)
	if err != nil {
		log.Fatalf("xlog new failed, %v", err)
	}
	defer xl.Sync()

	xlConsole, _ := xlog.NewConsole(xlog.DebugLevel)
	defer xlConsole.Sync()
	for {
		LogFile(xl)
		//LogConsole(xlConsole)
		time.Sleep(time.Second * 1)
	}
}
```