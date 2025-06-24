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
 @Time    : 2025/4/15 -- 13:56
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xid xid/snow-config.go
*/

package xid

import (
	"bytes"
)

const ErrorIdResult int64 = -1

type SnowFlakeModeConfig struct {
	SName           string        `properties:"sName"`
	Namespace       string        `properties:"namespace"`
	Timeout         int64         `properties:"timeout"`
	Mode            SnowFlakeMode `properties:"mode"`
	TimeStampBitNum uint16        `properties:"timestampBitNum"`
	SequenceBitNum  uint16        `properties:"sequenceBitNum"`
}

func NewSnowFlakeModeConfig(sName, namespace string) SnowFlakeModeConfig {
	var cfg = SnowFlakeModeConfig{
		SName:     sName,
		Namespace: namespace,
	}
	_ = cfg.Init()
	return cfg
}

type SegmentModeConfig struct {
	ServName   string `properties:"servName"`
	Namespace  string `properties:"namespace"`
	Timeout    int64  `properties:"timeout"`
	DefaultQps int64  `properties:"defaultQps"`
}

func GetGeneratorKey(servName, namespace string) string {
	var buf bytes.Buffer
	if servName != "" {
		buf.WriteString(servName)
	}
	if namespace != "" {
		buf.WriteString(".")
		buf.WriteString(namespace)
	}
	return buf.String()
}

func (s *SnowFlakeModeConfig) GetKey() string {
	return GetGeneratorKey(s.SName, s.Namespace)
}

func (s *SnowFlakeModeConfig) Init() error {
	if s.Mode == NORMAL_MOD {
		s.SequenceBitNum = DEFAULT_SEQUENCE_BIT_NUM
		s.TimeStampBitNum = DEFAULT_TIMESTAMP_BIT_NUM
		s.Timeout = DEFAULT_TIMEOUT
	}
	return nil
}

func (c *SegmentModeConfig) GetKey() string {
	return GetGeneratorKey(c.ServName, c.Namespace)
}
