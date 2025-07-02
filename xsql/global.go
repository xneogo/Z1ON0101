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
 @Time    : 2025/4/15 -- 12:05
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xsql xsql/global.go
*/

package xsql

import (
	"context"
	"github.com/xneogo/Z1ON0101/xsql/sqlutils"

	"github.com/pkg/errors"
	"github.com/xneogo/Z1ON0101/xsql/factory"
	"github.com/xneogo/Z1ON0101/xsql/xbuilder"
	xmanager "github.com/xneogo/Z1ON0101/xsql/xdb"
	"github.com/xneogo/Z1ON0101/xsql/xscanner"
)

var Constructor constructor

type constructor struct {
	_scanner xscanner.XScanner
	_builder xbuilder.XBuilder
}

func init() {
	Constructor = constructor{
		_scanner: xscanner.XScanner{},
		_builder: xbuilder.XBuilder{},
	}
}

func (c constructor) GetBuilder() factory.Builder {
	return c._builder
}

func (c constructor) GetScanner() factory.Scanner {
	return c._scanner
}

func (c constructor) ComplexSelect(tableF func() []string, target any, query string, args ...interface{}) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := xmanager.SqlExecDefault(ctx, func(dbx *xmanager.DBX, tables []interface{}) error {
			return dbx.SelectWrapper(tables, &target, query, args)
		}, tableF()...)
		return err
	}
}

func (c constructor) ComplexExec(tableF func() []string, query string, args ...interface{}) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		err := xmanager.SqlExecDefault(ctx, func(dbx *xmanager.DBX, tables []interface{}) error {
			_, er := dbx.ExecWrapper(tables, query, args)
			return er
		}, tableF()...)
		return err
	}
}

func NotFound(err error) bool {
	return errors.Is(err, sqlutils.ErrScannerEmptyResult)
}
