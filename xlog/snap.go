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
 @Time    : 2024/10/28 -- 11:46
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: snap.go
*/

package xlog

import (
	"sync"
	"sync/atomic"
	"time"
)

const logCacheMaxSize = 10

// log counters of different level
var (
	// Deprecated
	cnTrace int64

	cnDebug int64
	cnInfo  int64
	cnWarn  int64
	cnError int64
	cnFatal int64
	cnPanic int64
)

var (
	// 每次日志统计的起始时刻
	cnStamp int64

	slogMutex sync.Mutex

	// Deprecated, 使用logInfos
	logs []string

	// 新版本日志详情信息
	logInfos []*LogInfo
)

// LogStatInfo 日志统计信息
type LogStatInfo struct {
	// 统计的起始时刻
	StatStartTime int64
	// 各日志级别的日志数量
	Counters map[Level]int64
	LogInfos []*LogInfo
}

type LogInfo struct {
	Head map[string]interface{} `json:"head"`
	Body map[string]interface{} `json:"body"`
}

func init() {
	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

func newLogInfo(head, body map[string]interface{}) *LogInfo {
	return &LogInfo{
		Head: head,
		Body: body,
	}
}

// DumpLogStatInfo 传入log header和body, 方便报警处理
func DumpLogStatInfo() *LogStatInfo {
	slogMutex.Lock()
	defer slogMutex.Unlock()

	return &LogStatInfo{
		StatStartTime: dumpStatStartTime(),
		Counters:      dumpLogCounters(),
		LogInfos:      dumpLogInfos(),
	}
}

// AddLogInfo
// head, body用于新版日志报警上报
func AddLogInfo(head, body map[string]interface{}) {
	slogMutex.Lock()
	defer slogMutex.Unlock()
	l := newLogInfo(head, body)
	logInfos = append(logInfos, l)
	if len(logInfos) > logCacheMaxSize {
		logInfos = logInfos[len(logInfos)-logCacheMaxSize:]
	}
}

// 不再导出trace level的日志数量
func dumpLogCounters() map[Level]int64 {
	return map[Level]int64{
		DebugLevel: atomic.SwapInt64(&cnDebug, 0),
		InfoLevel:  atomic.SwapInt64(&cnInfo, 0),
		WarnLevel:  atomic.SwapInt64(&cnWarn, 0),
		ErrorLevel: atomic.SwapInt64(&cnError, 0),
		FatalLevel: atomic.SwapInt64(&cnFatal, 0),
		PanicLevel: atomic.SwapInt64(&cnPanic, 0),
	}
}

func dumpLogInfos() []*LogInfo {
	l := logInfos
	logInfos = make([]*LogInfo, 0, 128)
	return l
}

func dumpStatStartTime() int64 {
	return atomic.SwapInt64(&cnStamp, time.Now().Unix())
}
