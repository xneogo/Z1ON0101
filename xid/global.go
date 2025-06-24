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
 @Time    : 2025/4/15 -- 13:57
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xid xid/global.go
*/

package xid

import "context"

var defaultIdConf = NewSnowFlakeModeConfig("_default_", "_default_")
var defaultIdGen, _ = CreateSnowFlakeIdGen(context.Background(), &defaultIdConf, 0)

func GetID(ctx context.Context) (int64, error) {
	return defaultIdGen.GetId(ctx)
}

// SlowIdGenerator
// 慢id生成器，适合id产生不是非常快的场景,基于毫秒时间戳，每毫秒最多产生2个id，过快会自动阻塞，直到毫秒递增
// id表示可以再52bit完成，用double表示不会丢失精度，javascript等弱类型语音可以直接使用
type SlowIdGenerator interface {
	GenId(tp string) (int64, error)
	GetIdStamp(sid int64) int64
	GetIdWithStamp(stamp int64) int64
}

// SnowFlakeIdGenerator 雪花id生成器
type SnowFlakeIdGenerator interface {
	// GenId 雪花id生成
	GenId() (int64, error)
	// GetIdStamp 获取snowflakeid生成时间戳，单位ms
	GetIdStamp(sid int64) int64
	// GetIdWithStamp 按给定的时间点构造一个起始snowflakeid，一般用于区域判断
	GetIdWithStamp(stamp int64) int64
}
