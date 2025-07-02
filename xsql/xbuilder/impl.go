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
 @Time    : 2025/7/1 -- 17:14
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xbuilder xsql/xbuilder/impl.go
*/

package xbuilder

import (
	"context"
	"github.com/xneogo/Z1ON0101/xsql/factory"

	"github.com/xneogo/Z1ON0101/xsql/sqlutils"
)

type XBuilder struct{}

func (b XBuilder) Build() {

}
func (b XBuilder) BuildSelect(tableName string, where map[string]interface{}, selectedField []string) (query string, args []interface{}, err error) {
	sc, err1 := sqlutils.ParseWhere(where)
	if err1 != nil {
		err = err1
		return
	}
	return sqlutils.BuildSelect(tableName, selectedField, sc)
}
func (b XBuilder) BuildSelectWithContext(ctx context.Context, tableName string, where map[string]interface{}, selectedField []string) (query string, args []interface{}, err error) {
	sc, err1 := sqlutils.ParseWhere(where)
	if err1 != nil {
		err = err1
		return
	}
	return sqlutils.BuildSelectWithContext(ctx, tableName, selectedField, sc)
}
func (b XBuilder) BuildUpdate(tableName string, where map[string]interface{}, update map[string]interface{}) (string, []interface{}, error) {
	clauses, release, err := sqlutils.ParseDMLWhere(where)
	if nil != err {
		return "", nil, err
	}
	defer release()
	return sqlutils.BuildUpdate(tableName, update, clauses)
}
func (b XBuilder) BuildDelete(tableName string, where map[string]interface{}) (string, []interface{}, error) {
	conditions, release, err := sqlutils.GetWhereConditions(where)
	if nil != err {
		return "", nil, err
	}
	defer release()
	return sqlutils.BuildDelete(tableName, conditions...)
}

func (b XBuilder) BuildInsert(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	return sqlutils.BuildInsert(tableName, data, sqlutils.CommonInsert)
}

func (b XBuilder) BuildUpsert(tableName string, data map[string]interface{}) (string, []interface{}, error) {
	return sqlutils.BuildUpsert(tableName, data, sqlutils.CommonInsert)
}

func (b XBuilder) BuildInsertIgnore(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	return sqlutils.BuildInsert(tableName, data, sqlutils.IgnoreInsert)
}

func (b XBuilder) BuildReplaceIgnore(tableName string, data []map[string]interface{}) (string, []interface{}, error) {
	return sqlutils.BuildInsert(tableName, data, sqlutils.ReplaceInsert)
}

func (b XBuilder) AggregateQuery(ctx context.Context, db factory.XDB, tableName string, where map[string]interface{}, aggregate factory.AggregateSymbolBuilder) (factory.ResultResolver, error) {
	cond, vals, err := b.BuildSelect(tableName, where, []string{aggregate.Symbol()})
	if nil != err {
		return resultResolve{0}, err
	}
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err {
		return resultResolve{0}, err
	}
	var result interface{}
	for rows.Next() {
		err = rows.Scan(&result)
	}
	rows.Close()
	return resultResolve{result}, err
}
