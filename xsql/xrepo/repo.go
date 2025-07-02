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

import "github.com/xneogo/matrix/msql"

// XRepo only for convenient, a wrapper of zsql interfaces
// to make a difference between fastsql and xsql
type XRepo[EntityObj any] interface {
	msql.RepoModel[EntityObj]
}

type XQueryRequest[EntityObj any] interface {
	msql.QueryRequest[EntityObj]
}

type XInsertRequest[EntityObj any] interface {
	msql.InsertRequest[EntityObj]
}

type XUpsertRequest[EntityObj any] interface {
	msql.UpsertRequest[EntityObj]
}

type XUpdateRequest interface {
	msql.UpdateRequest
}

type XDeleteRequest interface {
	msql.DeleteRequest
}

type XComplexRequest[EntityObj any] interface {
	msql.ComplexRequest[EntityObj]
}
