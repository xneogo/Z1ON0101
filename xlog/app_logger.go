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
 @Time    : 2024/10/28 -- 11:42
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: app_logger.go
*/

package xlog

import (
	"context"
)

// InitAppLog init app logger
func InitAppLog(logRoot, fileName string, level Level) (err error) {
	appLogger, err = New(logRoot, fileName, level, JSONFormatType, false, 0, false)
	if err != nil {
		return err
	}
	appLogger.skip = 4
	return nil
}

func InitAppLogHeaders(logRoot, fileName string, level Level, extraHeaders map[string]interface{}) (err error) {
	return initAppLogCore(logRoot, fileName, level, JSONFormatType, extraHeaders, 0, false)
}

// initAppLogCore init app logger with extra headers
func initAppLogCore(logRoot, fileName string, level Level, format FormatType, extraHeaders map[string]interface{}, maxRetainDays int, compress bool) (err error) {
	appLogger, err = NewPassThrough(logRoot, fileName, level, format, false, extraHeaders, maxRetainDays, compress)
	if err != nil {
		return err
	}
	appLogger.skip = 4
	return nil
}

// SetAppLogService set service in stat logger
func SetAppLogService(service string) {
	appLogger.SetService(service)
}

// SetAppLogLevel set level of app logger
func SetAppLogLevel(level string) {
	Infof(context.TODO(), "set app logger level from %s to %s", convertLevelToStr(AppLogLevel()), level)
	appLogger.SetLevel(ConvertLevel(level))
}

// AppLogLevel  return level of app logger
func AppLogLevel() Level {
	return appLogger.Level()
}

// SetAppLogSkip set skip of app logger
func SetAppLogSkip(skip int) {
	appLogger.skip = skip
}

// AppLogSync app logger sync
func AppLogSync() {
	appLogger.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(ctx context.Context, args ...interface{}) {
	appLogger.Debug(ctx, args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(ctx context.Context, template string, args ...interface{}) {
	appLogger.Debugf(ctx, template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg
func Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Debugw(ctx, msg, keysAndValues...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(ctx context.Context, args ...interface{}) {
	appLogger.Info(ctx, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(ctx context.Context, template string, args ...interface{}) {
	appLogger.Infof(ctx, template, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With
func Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Infow(ctx, msg, keysAndValues...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(ctx context.Context, args ...interface{}) {
	appLogger.Warn(ctx, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(ctx context.Context, template string, args ...interface{}) {
	appLogger.Warnf(ctx, template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With
func Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Warnw(ctx, msg, keysAndValues...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(ctx context.Context, args ...interface{}) {
	appLogger.Error(ctx, args...)
}

// Errorf uses fmt.Sprintf to log a templated message
func Errorf(ctx context.Context, template string, args ...interface{}) {
	appLogger.Errorf(ctx, template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Errorw(ctx, msg, keysAndValues...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(ctx context.Context, args ...interface{}) {
	appLogger.Fatal(ctx, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(ctx context.Context, template string, args ...interface{}) {
	appLogger.Fatalf(ctx, template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With
func Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Fatalw(ctx, msg, keysAndValues...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(ctx context.Context, args ...interface{}) {
	appLogger.Panic(ctx, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func Panicf(ctx context.Context, template string, args ...interface{}) {
	appLogger.Panicf(ctx, template, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With
func Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	appLogger.Panicw(ctx, msg, keysAndValues...)
}
