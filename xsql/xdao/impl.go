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
 @Time    : 2025/7/1 -- 17:19
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xdao xsql/xdao/impl.go
*/

package xdao

import (
	"context"
	"errors"
	"fmt"
	"github.com/xneogo/zion/xsql/factory"
	"github.com/xneogo/zion/xsql/sqlutils"
	"github.com/xneogo/zion/xsql/xbuilder"
)

type XDao[DObj any] interface {
	factory.DaoModel[DObj]
}

type DefaultDao[DObj any] struct {
	tableName func() string
	omits     func() []string
	_scanner  factory.Scanner
	_builder  factory.Builder
	bind      factory.BindFunc
}

// TableName anyone wrapper *DefaultDao should write their own
func (dao *DefaultDao[DObj]) TableName() string { return dao.tableName() }

// Omits anyone wrapper *DefaultDao should write their own
func (dao *DefaultDao[DObj]) Omits() []string { return dao.omits() }

func (dao *DefaultDao[DObj]) Init(cons factory.SqlConstructor, tableName func() string, omits func() []string, b factory.BindFunc) {
	dao._builder = cons.GetBuilder()
	dao._scanner = cons.GetScanner()
	dao.omits = omits
	dao.tableName = tableName
	dao.bind = b
}

func (dao *DefaultDao[DObj]) GetScanner() factory.Scanner {
	return dao._scanner
}

func (dao *DefaultDao[DObj]) GetBuilder() factory.Builder {
	return dao._builder
}

func (dao *DefaultDao[DObj]) SelectOne(ctx context.Context, db factory.XDB, where map[string]interface{}) (res DObj, err error) {
	if nil == db {
		return res, errors.New("manager.XDB object couldn't be nil")
	}
	tar := sqlutils.CopyWhere(where)
	if _, ok := tar["_limit"]; !ok {
		tar["_limit"] = []uint{0, 1}
	}
	cond, vals, err := dao._builder.BuildSelectWithContext(ctx, dao.TableName(), tar, dao.Omits())
	if nil != err {
		return res, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err || nil == rows {
		return res, err
	}
	defer rows.Close()
	err = dao._scanner.Scan(rows, &res, dao.bind)
	fmt.Println("res", res)
	return res, err
}

func (dao *DefaultDao[DObj]) SelectMulti(ctx context.Context, db factory.XDB, where map[string]interface{}) (res []DObj, err error) {
	if nil == db {
		return res, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildSelectWithContext(ctx, dao.TableName(), where, dao.Omits())
	if nil != err {
		return nil, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err || nil == rows {
		return nil, err
	}
	defer rows.Close()
	err = dao._scanner.Scan(rows, &res, dao.bind)
	return res, err
}

func (dao *DefaultDao[DObj]) Insert(ctx context.Context, db factory.XDB, data []map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildInsert(dao.TableName(), data)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err || nil == result {
		return 0, err
	}
	return result.LastInsertId()
}

func (dao *DefaultDao[DObj]) Upsert(ctx context.Context, db factory.XDB, data map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildUpsert(dao.TableName(), data)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err || nil == result {
		return 0, err
	}
	return result.LastInsertId()
}

func (dao *DefaultDao[DObj]) Update(ctx context.Context, db factory.XDB, where, data map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildUpdate(dao.TableName(), where, data)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DefaultDao[DObj]) Delete(ctx context.Context, db factory.XDB, where map[string]interface{}) (int64, error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildDelete(dao.TableName(), where)
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	result, err := db.ExecContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DefaultDao[DObj]) CountOf(ctx context.Context, db factory.XDB, where map[string]interface{}) (count int, err error) {
	if nil == db {
		return 0, errors.New("manager.XDB object couldn't be nil")
	}
	cond, vals, err := dao._builder.BuildSelect(dao.TableName(), where, []string{xbuilder.AggregateCount("*").Symbol()})
	if nil != err {
		return 0, err
	}
	xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
	rows, err := db.QueryContext(ctx, cond, vals...)
	if nil != err {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	return
}

// ComplexQuery
// you can use this default logic or
// you can build your own query logic with or without tableName or columns
// depends on your ToSql func
func ComplexQuery[ans any](tableName string, columns ...string) factory.ComplexQueryMod[ans] {
	return func(
		ctx context.Context,
		db factory.XDB,
		scanner factory.Scanner,
		f factory.ToSql,
		bind factory.BindFunc,
	) (res []ans, err error) {
		if nil == db {
			return nil, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName, columns...)
		if nil != err {
			return nil, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		xlog.Infof(ctx, "build cond: %s vals: %v", cond, vals)
		rows, err := db.QueryContext(ctx, cond, vals...)
		if nil != err || nil == rows {
			return nil, err
		}
		defer rows.Close()
		err = scanner.Scan(rows, &res, bind)
		return res, err
	}
}

func ComplexExec(tableName string) factory.ComplexExecMod {
	return func(
		ctx context.Context,
		db factory.XDB,
		f factory.ToSql,
	) (int64, error) {
		if nil == db {
			return 0, errors.New("manager.XDB object couldn't be nil")
		}
		cond, vals, err := f(tableName)
		if nil != err {
			return 0, err
		}
		xlog.Debugf(ctx, "build cond: %s vals: %v", cond, vals)
		result, err := db.ExecContext(ctx, cond, vals...)
		if nil != err {
			return 0, err
		}
		return result.RowsAffected()
	}
}
