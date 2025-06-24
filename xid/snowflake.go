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
 @Time    : 2025/4/15 -- 13:55
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xid xid/snowflake.go
*/

package xid

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type SnowFlakeMode int32

const (
	NORMAL_MOD SnowFlakeMode = 0 // 标准雪花算法模式
)

// snowflake 默认 bit 分配方案
const (
	DEFAULT_TIMESTAMP_BIT_NUM = 41
	DEFAULT_WORK_ID_BIT_NUM   = 2
	DEFAULT_SEQUENCE_BIT_NUM  = 10
	// DEFAULT_WORK_ID_BIT_NUM   = 10
	// DEFAULT_SEQUENCE_BIT_NUM  = 12
	SLOW_SEQUENCE_BIT_NUM = 1

	START_TIME int64 = 1672502400000 // 2023-01-01 00:00:00

	DEFAULT_TIMEOUT = 50
)

type SnowFlakeIdGen struct {
	mutex sync.Mutex

	key      string
	lastTime int64

	workId   int64 // 当前 servId
	sequence int64 // 当前自增序列

	MaxTimestamp int64
	MaxSequence  int64

	// 配置加载的 各段bit位长度
	WorkIdBitNum    uint16
	TimeStampBitNum uint16
	SequenceBitNum  uint16

	timeout int64 // 毫秒级超时

	concurrency      int64
	concurrencyLimit int64
}

func CreateSnowFlakeIdGen(ctx context.Context, cfg *SnowFlakeModeConfig, workId int64) (idGen *SnowFlakeIdGen, err error) {
	idGen = &SnowFlakeIdGen{
		key:             cfg.GetKey(),
		workId:          workId,
		sequence:        0,
		MaxTimestamp:    -1 ^ (-1 << cfg.TimeStampBitNum),
		MaxSequence:     -1 ^ (-1 << cfg.SequenceBitNum),
		WorkIdBitNum:    DEFAULT_WORK_ID_BIT_NUM,
		TimeStampBitNum: cfg.TimeStampBitNum,
		SequenceBitNum:  cfg.SequenceBitNum,
		timeout:         cfg.Timeout,
		concurrency:     0,
	}
	idGen.concurrencyLimit = (idGen.MaxSequence + 1) * idGen.timeout
	idGen.lastTime = idGen.CurrentElapsedTime()
	return idGen, nil
}

func (sfg *SnowFlakeIdGen) GetId(ctx context.Context) (int64, error) {
	if !sfg.getToken() {
		return 0, fmt.Errorf("excessive concurrency")
	}
	defer sfg.releaseToken()

	sfg.mutex.Lock()
	defer sfg.mutex.Unlock()

	ts := sfg.CurrentElapsedTime()
	var id int64
	for {
		if ts < sfg.lastTime { // 发生时钟回拨，等待
			time.Sleep(time.Duration(ts - sfg.lastTime))
			ts = sfg.CurrentElapsedTime()
			continue
		} else if ts == sfg.lastTime { // 同一时间段内
			sfg.sequence = (sfg.sequence + 1) & sfg.MaxSequence
			if sfg.sequence == 0 { // 耗尽，等待到下一时间段
				ts = sfg.tilNextMillis()
			}
		} else {
			sfg.sequence = 0 // 已到下一时间段
		}
		id = sfg.ToId(ts, sfg.sequence) // 生成id
		sfg.lastTime = ts               // 标记时间
		break
	}
	return id, nil
}

// ToSnowFlakeTime 按照时间 bit 位格式化时间
func (sfg *SnowFlakeIdGen) ToSnowFlakeTime(t int64) int64 {
	return (t/1e6 - START_TIME) & sfg.MaxTimestamp
}

func (sfg *SnowFlakeIdGen) CurrentElapsedTime() int64 {
	return sfg.ToSnowFlakeTime(time.Now().UnixNano())
}

func (sfg *SnowFlakeIdGen) ToId(lastTime int64, sequence int64) int64 {
	return sfg.toId(lastTime, sequence)
}

func (sfg *SnowFlakeIdGen) toId(lastTime int64, sequence int64) int64 {
	return lastTime<<(sfg.SequenceBitNum+sfg.WorkIdBitNum) |
		sfg.workId<<(sfg.SequenceBitNum) |
		sequence
}

func (sfg *SnowFlakeIdGen) tilNextMillis() int64 {

	un := time.Now().UnixNano()
	dw := 1000000 - un%1000000
	time.Sleep(time.Duration(dw))

	i := sfg.CurrentElapsedTime()
	for i <= sfg.lastTime {
		i = sfg.CurrentElapsedTime()
	}
	return i
}

func (sfg *SnowFlakeIdGen) getToken() bool {

	if atomic.LoadInt64(&sfg.concurrency) >= sfg.concurrencyLimit {
		return false
	}
	atomic.AddInt64(&sfg.concurrency, 1)
	return true
}

func (sfg *SnowFlakeIdGen) releaseToken() {
	atomic.AddInt64(&sfg.concurrency, -1)
}

// UnmarshalId 解析标准snowflake算法生成的id
func (sfg *SnowFlakeIdGen) UnmarshalId(ctx context.Context, id int64) (timeStampMs, workId, sequence int64) {
	var MaxTimestamp int64
	var MaxSequence int64
	var MaxWorkId int64
	MaxTimestamp = -1 ^ (-1 << sfg.TimeStampBitNum)
	MaxSequence = -1 ^ (-1 << sfg.SequenceBitNum)
	MaxWorkId = -1 ^ (-1 << sfg.WorkIdBitNum)
	return ((id >> (sfg.WorkIdBitNum + sfg.SequenceBitNum)) & MaxTimestamp) + START_TIME,
		(id >> sfg.SequenceBitNum) & MaxWorkId, id & MaxSequence
}
