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
 @Time    : 2024/10/28 -- 11:41
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: logger.go
*/

package xlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"gitee.com/go-mid/infra/xcontext"

	"github.com/shawnfeng/lumberjack.v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gitee.com/go-mid/infra/xtrace/typealias"
)

// FormatType log format type
type FormatType string

// log format type, it's encoder in zap
const (
	JSONFormatType      FormatType = "json"           // 业务日志建议用格式化的json
	AccessLogFormatType            = "accesslog.text" // 目前主要用于accesslog
	PlainTextFormatType            = "plain.text"     // 用于测试，建议使用以上两种方式
)

const (
	defaultInitSkip = 3
)

const tracingErrorLogTagName = "error"

// A Level is a logging priority. Higher levels are more important.
type Level = zapcore.Level

// log level
const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

// XLogger wrapper of zap.Logger
type XLogger struct {
	zl       *zap.Logger
	zlSugar  *zap.SugaredLogger
	level    zap.AtomicLevel
	service  string
	hostName string

	extraHeaders map[string]interface{}

	skip int // runtime由xlog进行计算，没有使用zap的caller
}

// Level return level
func (xl *XLogger) Level() Level {
	return xl.level.Level()
}

// SetLevel set level  runtime
func (xl *XLogger) SetLevel(level Level) {
	xl.level.SetLevel(level)
}

// SetService set service name
func (xl *XLogger) SetService(service string) {
	xl.service = service
}

// SetSkip set skip
func (xl *XLogger) SetSkip(skip int) {
	xl.skip = skip
}

// Sync force sync data to log faile
func (xl *XLogger) Sync() {
	xl.zl.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func (xl *XLogger) Debug(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	xl.Debugw(ctx, msg)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (xl *XLogger) Debugf(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > DebugLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Debug(msg)
	atomic.AddInt64(&cnDebug, 1)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg
func (xl *XLogger) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, false, "DEBUG")
	body := xl.buildBody(msg, keysAndValues...)
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Debugw(msg, kvs...)
	if xl.Level() > DebugLevel {
		return
	}
	atomic.AddInt64(&cnDebug, 1)
}

// Info uses fmt.Sprint to construct and log a message.
func (xl *XLogger) Info(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	xl.Infow(ctx, msg)
}

// Infof uses fmt.Sprintf to log a templated message.
func (xl *XLogger) Infof(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > InfoLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Info(msg)
	atomic.AddInt64(&cnInfo, 1)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With
func (xl *XLogger) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, false, "INFO")
	body := xl.buildBody(msg, keysAndValues...)
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Infow(msg, kvs...)
	if xl.Level() > InfoLevel {
		return
	}
	atomic.AddInt64(&cnInfo, 1)
}

// Infom logs a message with some additional context. The variadic map
// key-string value-interface pairs are treated as they are in With
func (xl *XLogger) Infom(ctx context.Context, msg string, keyAndValueMap map[string]interface{}) {
	head := xl.buildHead(ctx, false, "INFO")
	keyAndValueMap["msg"] = msg
	kvs := append([]interface{}{}, "head", head, "body", keyAndValueMap)
	xl.zlSugar.Infow(msg, kvs...)
	atomic.AddInt64(&cnInfo, 1)
}

// Warn uses fmt.Sprint to construct and log a message.
func (xl *XLogger) Warn(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	head := xl.buildHead(ctx, true, "WARN")
	body := xl.buildBody(msg)
	xl.warnw(msg, head, body)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (xl *XLogger) Warnf(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > WarnLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Warn(msg)
	atomic.AddInt64(&cnWarn, 1)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With
func (xl *XLogger) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, true, "WARN")
	body := xl.buildBody(msg, keysAndValues...)
	xl.warnw(msg, head, body)
}

// 为了runtim层次相同而封装
func (xl *XLogger) warnw(msg string, head, body map[string]interface{}) {
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Warnw(msg, kvs...)
	if xl.Level() > WarnLevel {
		return
	}
	atomic.AddInt64(&cnWarn, 1)
}

// Error uses fmt.Sprint to construct and log a message.
func (xl *XLogger) Error(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	head := xl.buildHead(ctx, true, "ERROR")
	body := xl.buildBody(msg)
	xl.errorw(msg, head, body)
}

// Errorf uses fmt.Sprintf to log a templated message
func (xl *XLogger) Errorf(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > ErrorLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Error(msg)
	atomic.AddInt64(&cnError, 1)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (xl *XLogger) Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, true, "ERROR")
	body := xl.buildBody(msg, keysAndValues...)
	xl.errorw(msg, head, body)
}

func (xl *XLogger) errorw(msg string, head, body map[string]interface{}) {
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Errorw(msg, kvs...)
	atomic.AddInt64(&cnError, 1)
	AddLogInfo(head, body)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (xl *XLogger) Fatal(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	head := xl.buildHead(ctx, true, "FATAL")
	body := xl.buildBody(msg)
	xl.fatalw(msg, head, body)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (xl *XLogger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > FatalLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Fatal(msg)
	atomic.AddInt64(&cnFatal, 1)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With
func (xl *XLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, true, "FATAL")
	body := xl.buildBody(msg, keysAndValues...)
	xl.fatalw(msg, head, body)
	atomic.AddInt64(&cnFatal, 1)
	AddLogInfo(head, body)
}

func (xl *XLogger) fatalw(msg string, head, body map[string]interface{}) {
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Fatalw(msg, kvs...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (xl *XLogger) Panic(ctx context.Context, args ...interface{}) {
	msg := fmt.Sprint(args...)
	head := xl.buildHead(ctx, true, "PANIC")
	body := xl.buildBody(msg)
	xl.panicw(msg, head, body)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (xl *XLogger) Panicf(ctx context.Context, template string, args ...interface{}) {
	if xl.Level() > PanicLevel {
		return
	}
	msg := xl.templateWrapTraceMsg(ctx, template, args...)
	xl.zl.Panic(msg)
	atomic.AddInt64(&cnPanic, 1)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With
func (xl *XLogger) Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	head := xl.buildHead(ctx, true, "PANIC")
	body := xl.buildBody(msg, keysAndValues...)
	xl.panicw(msg, head, body)
}

func (xl *XLogger) panicw(msg string, head, body map[string]interface{}) {
	kvs := append([]interface{}{}, "head", head, "body", body)
	xl.zlSugar.Panicw(msg, kvs...)
	atomic.AddInt64(&cnPanic, 1)
	AddLogInfo(head, body)
}

func (xl *XLogger) templateWrapTraceMsg(ctx context.Context, template string, args ...interface{}) string {
	template, args = xl.templateWrapTrace(ctx, template, args...)
	return fmt.Sprintf(template, args...)
}
func (xl *XLogger) templateWrapTrace(ctx context.Context, template string, args ...interface{}) (string, []interface{}) {
	template = "%s\t" + template

	var traceID string
	if v, ok := xl.getTraceID(ctx); ok {
		traceID = v
	}
	var argg = []interface{}{
		traceID,
	}
	if args != nil {
		argg = append(argg, args...)
	}
	return template, argg
}

func NewPassThrough(logRoot, fileName string, level Level, formatType FormatType, debug bool, extraHeaders map[string]interface{}, maxRetainDays int, compress bool) (logger *XLogger, err error) {
	logger, err = New(logRoot, fileName, level, formatType, debug, maxRetainDays, compress)
	if err != nil {
		return nil, err
	}
	logger.extraHeaders = extraHeaders
	return logger, nil
}

// New create and init logger, rotate by hour
// fileName: 日志文件名称，如: test.log
func New(logRoot, fileName string, level Level, formatType FormatType, debug bool, maxRetainDays int, compress bool) (logger *XLogger, err error) {
	logger = &XLogger{level: zap.NewAtomicLevel(), skip: defaultInitSkip, extraHeaders: make(map[string]interface{})}
	logger.hostName, err = os.Hostname()
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)
	logFile := ""
	if logRoot != "" && fileName != "" {
		logFile = logRoot + "/" + fileName
	}
	var w io.Writer
	if logFile != "" {
		ljWriter := lumberjack.NewLogger(logFile, 1024, maxRetainDays, 0, true, compress)

		go func() {
			for {
				now := time.Now().Unix()
				duration := 3600 - now%3600
				select {
				case <-time.After(time.Second * time.Duration(duration)):
					ljWriter.Rotate()
				}
			}
		}()
		w = ljWriter
	} else {
		w = os.Stdout
	}

	zc := zapcore.AddSync(w)

	var core zapcore.Core
	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "",
		NameKey:        "",
		CallerKey:      "",
		MessageKey:     "",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch formatType {
	case AccessLogFormatType:
		encoderConfig := TextEncoderConfig{
			&zapEncoderConfig,
			"\t",
		}
		core = zapcore.NewCore(
			NewAccessLogEncoder(encoderConfig),
			zc,
			logger.level,
		)
	case JSONFormatType:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zapEncoderConfig),
			zc,
			logger.level,
		)
	case PlainTextFormatType:
		zapEncoderConfig.TimeKey = "ts"
		zapEncoderConfig.LevelKey = "level"
		zapEncoderConfig.NameKey = "logger"
		zapEncoderConfig.CallerKey = "caller"
		zapEncoderConfig.MessageKey = "msg"
		zapEncoderConfig.StacktraceKey = "stacktrace"
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapEncoderConfig),
			zc,
			logger.level,
		)
	}

	// for debug, printed to console
	if debug {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapEncoderConfig),
			zapcore.Lock(os.Stdout),
			logger.level,
		)
		core = zapcore.NewTee(core, consoleCore)
	}

	// AccessLogFormatType不需要caller，其他格式的日志默认添加caller
	if formatType != AccessLogFormatType {
		// add caller
		caller := zap.AddCaller()
		// add skip
		skip := zap.AddCallerSkip(defaultInitSkip)
		logger.zl = zap.New(core, caller, skip)
	} else {
		logger.zl = zap.New(core)

	}
	logger.zlSugar = logger.zl.Sugar()
	return logger, nil
}

// NewConsole some cli tools will use it
func NewConsole(level Level) (logger *XLogger, err error) {
	var core zapcore.Core
	logger = &XLogger{level: zap.NewAtomicLevel()}
	logger.hostName, err = os.Hostname()
	if err != nil {
		return nil, err
	}
	logger.level.SetLevel(level)
	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "",
		NameKey:        "",
		CallerKey:      "",
		MessageKey:     "",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core = zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapEncoderConfig),
		zapcore.Lock(os.Stdout),
		logger.level,
	)
	// add caller
	caller := zap.AddCaller()
	// add skip
	skip := zap.AddCallerSkip(defaultInitSkip)
	logger.zl = zap.New(core, caller, skip)
	logger.zlSugar = logger.zl.Sugar()
	return
}

func getRuntimeInfo(skip int) (function, filename string, lineno int) {
	function = "???"
	pc, filename, lineno, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}

	return
}

func (xl *XLogger) getTraceID(ctx context.Context) (string, bool) {
	span := typealias.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(typealias.SpanContext); ok {
			return fmt.Sprint(sc.TraceID()), true
		}
	}
	return "", false
}

func (xl *XLogger) buildHead(ctx context.Context, needCaller bool, level string) (head map[string]interface{}) {
	head = map[string]interface{}{
		"ts":       time.Now().Format("2006-01-02T15:04:05.000Z0700"),
		"level":    level,
		"host":     xl.hostName,
		"ctx_lane": xcontext.GetControlRouteGroupWithMasterName(ctx, ""),
	}

	for k, v := range xl.extraHeaders {
		head[k] = v
	}

	// service name
	if xl.service != "" {
		head["service"] = xl.service
	}
	// trace id
	span := typealias.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(typealias.SpanContext); ok {
			head["trace_id"] = fmt.Sprint(sc.TraceID())
		}

		// set tracing error log tag
		if ConvertLevel(level) >= ErrorLevel {
			span.SetTag(tracingErrorLogTagName, true)
		}
	}
	// uid
	uid := ctx.Value("uid")
	if uid != nil {
		head["uid"] = uid
	}
	// caller
	if needCaller {
		_, fileName, lineno := getRuntimeInfo(xl.skip)
		head["caller"] = fmt.Sprintf("%s:%d", fileName, lineno)
	}
	return head
}

func (xl *XLogger) buildBody(msg string, keysAndValues ...interface{}) map[string]interface{} {
	body := make(map[string]interface{})
	body["msg"] = msg
	for i := 0; i < len(keysAndValues); i += 2 {
		// ignore non-equal keys and values
		if i == len(keysAndValues)-1 {
			break
		}

		k, v := keysAndValues[i], keysAndValues[i+1]
		// fail-fast if key is not of string type
		if ks, ok := k.(string); !ok {
			break
		} else {
			body[ks] = v
		}
	}
	return body
}

// ConvertLevel convert string level to Level
func ConvertLevel(level string) Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "panic":
		return PanicLevel
	default:
		return InfoLevel
	}
}

func convertLevelToStr(level Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	default:
		return "info"
	}
}
