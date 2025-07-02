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
 @Time    : 2025/7/1 -- 18:18
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xsvc xsvc/flow.go
*/

package xsvc

import (
	"github.com/xneogo/Z1ON0101/xredis"
	"github.com/xneogo/Z1ON0101/xsql/factory"
)

type XFlow[EntityObj any] interface {
	Repo(factory.RepoModel[EntityObj]) *XFlowImpl[EntityObj]
	GetRepo() factory.RepoModel[EntityObj]

	KeyVal(*xredis.RedisClient) *XFlowImpl[EntityObj]
	GetKeyVal() *xredis.RedisClient
}

type XFlowImpl[EntityObj any] struct {
	Rp factory.RepoModel[EntityObj]
	Kv *xredis.RedisClient
	// Mongo
}

func (f *XFlowImpl[EntityObj]) Repo(rp factory.RepoModel[EntityObj]) *XFlowImpl[EntityObj] {
	f.Rp = rp
	return f
}

func (f *XFlowImpl[EntityObj]) GetRepo() factory.RepoModel[EntityObj] {
	return f.Rp
}

func (f *XFlowImpl[EntityObj]) KeyVal(kv *xredis.RedisClient) *XFlowImpl[EntityObj] {
	f.Kv = kv
	return f
}

func (f *XFlowImpl[EntityObj]) GetKeyVal() *xredis.RedisClient {
	return f.Kv
}

func NewXFlow[EntityObj any]() XFlow[EntityObj] {
	return &XFlowImpl[EntityObj]{}
}
