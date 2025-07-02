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
 @Time    : 2025/7/1 -- 17:23
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xrepo xsql/xrepo/impl.go
*/

package xrepo

import "github.com/xneogo/Z1ON0101/xsql/factory"

// XRepo only for convenient, a wrapper of zsql interfaces
// to make a difference between fastsql and xsql
type XRepo[EntityObj any] interface {
	factory.RepoModel[EntityObj]
}

type XQueryRequest[EntityObj any] interface {
	factory.QueryRequest[EntityObj]
}

type XInsertRequest[EntityObj any] interface {
	factory.InsertRequest[EntityObj]
}

type XUpsertRequest[EntityObj any] interface {
	factory.UpsertRequest[EntityObj]
}

type XUpdateRequest interface {
	factory.UpdateRequest
}

type XDeleteRequest interface {
	factory.DeleteRequest
}

type XComplexRequest[EntityObj any] interface {
	factory.ComplexRequest[EntityObj]
}
