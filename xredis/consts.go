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
 @Time    : 2024/11/1 -- 17:21
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: consts.go
*/

package xredis

import (
	"fmt"
	"time"
)

const (
	SpanLogKeyKey    = "key"
	SpanLogCacheType = "cache"
	SpanLogOp        = "op"
)

const (
	CacheDirtyExpireTime = time.Second * 60
)

var RedisNil = fmt.Sprintf("redis: nil")

type ConfigureType int

const (
	ConfigureTypeSimple ConfigureType = iota
	ConfigureTypeEtcd
	ConfigureTypeApollo
)

func (c ConfigureType) String() string {
	switch c {
	case ConfigureTypeSimple:
		return "simple"
	case ConfigureTypeEtcd:
		return "etcd"
	case ConfigureTypeApollo:
		return "apollo"
	default:
		return "unknown"
	}
}

const (
	DefaultRouteGroup = "default"

	DefaultRedisWrapper = ""
)

type CacheConfig struct {
	Namespace string    `json:"namespace"`
	Prefix    string    `json:"prefix"`
	CacheType CacheType `json:"cache_type"`
}
type CacheType string

const (
	CacheType_Redis   CacheType = "redis"
	CacheType_Gocache CacheType = "gocache"
	CacheType_Empty   CacheType = "empty"
)
