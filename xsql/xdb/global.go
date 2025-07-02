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
 @Time    : 2024/10/9 -- 17:51
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: global.go
*/

package xmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/xneogo/eins/colorlog"
	"github.com/xneogo/extensions/xbreaker"
	"github.com/xneogo/extensions/xtime"
	"github.com/xneogo/matrix/mconfig"
	"github.com/xneogo/matrix/mconfig/mobserver"
	"github.com/xneogo/matrix/msql"
)

// DBX mapper dbrouter DB
type DBX struct {
	*sqlx.DB
}

// GormDB ...
type GormDB struct {
	*gorm.DB
}

var dbManagerSelector *DbManagerSelector

const (
	DB_TYPE_MYSQL = "mysql"
)

const (
	WeirCluster = "weir"
)

func init() {
	var err error
	dbManagerSelector, err = NewDbManagerSelector()
	if err != nil {
		panic(err)
	}
}

// GetDBDefault generalization of GetDB in Manager
func GetDBDefault(ctx context.Context) (msql.XDBWrapper, error) {
	return dbManagerSelector.GetDBDefault(ctx)
}

// InitDBConf init db dynamic conf
func InitDBConf(ctx context.Context, confCenter mconfig.ConfigCenter) error {
	return dbManagerSelector.InitDbManagerConf(ctx, confCenter)
}

// ReloadDBConf reload db conf when change
func ReloadDBConf(ctx context.Context, confCenter mconfig.ConfigCenter, event mobserver.ChangeEvent) error {
	return dbManagerSelector.ReloadDbManagerConf(ctx, confCenter, event)
}

// SqlExecDefault generalization of SqlExec in Manager
// TODO: set cluster or db metrics param
func SqlExecDefault(ctx context.Context, query func(*DBX, []interface{}) error, tables ...string) error {
	fun := "XSQL.SqlExecDefault -->"

	// span, ctx := xtrace.StartSpanFromContext(ctx, "xsql.SqlExecDefault")
	// defer span.Finish()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}
	table := tables[0]

	// setDBSpanTags(span, WeirCluster, table, "")

	// check breaker
	if !xbreaker.Entry(WeirCluster, table) {
		colorlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, WeirCluster, table)
		return errors.New("sql cause breaker, because too many timeout")
	}

	db, err := GetDBDefault(ctx)
	if err != nil {
		return err
	}

	dbx := buildDBX(db.GetDb())
	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}
	st := xtime.NewTimeStat()
	err = query(dbx, tmptables)
	statMetricReqDur(ctx, WeirCluster, table, "sqlExecDefault", st.Millisecond())
	statMetricReqErrTotal(ctx, WeirCluster, table, "sqlExecDefault", err)
	// record breaker
	xbreaker.StatBreaker(WeirCluster, table, err)

	return err
}

func OrmExecDefault(ctx context.Context, query func(*GormDB, []interface{}) error, tables ...string) error {
	fun := "XSQL.OrmExecDefault -->"

	// span, ctx := xtrace.StartSpanFromContext(ctx, "xsql.OrmExecDefault")
	// defer span.Finish()

	if len(tables) <= 0 {
		return fmt.Errorf("tables is empty")
	}
	table := tables[0]

	// setDBSpanTags(span, WeirCluster, table, "")

	// check breaker
	if !xbreaker.Entry(WeirCluster, table) {
		colorlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, WeirCluster, table)
		return errors.New("sql cause breaker, because too many timeout")
	}
	db, err := GetDBDefault(ctx)
	if err != nil {
		return err
	}

	gormDB, err := buildGormDB(db.GetDb())
	if err != nil {
		return fmt.Errorf("%s build grom db, err: %v", fun, err)
	}
	var tmptables []interface{}
	for _, item := range tables {
		tmptables = append(tmptables, item)
	}
	st := xtime.NewTimeStat()
	err = query(gormDB, tmptables)
	statMetricReqDur(ctx, WeirCluster, table, "ormExecDefault", st.Millisecond())
	statMetricReqErrTotal(ctx, WeirCluster, table, "ormExecDefault", err)
	// record breaker
	xbreaker.StatBreaker(WeirCluster, table, err)
	return err
}

// NamedExecWrapper ...
func (db *DBX) NamedExecWrapper(tables []interface{}, query string, arg interface{}) (sql.Result, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.NamedExec(query, arg)
}

// NamedQueryWrapper ...
func (db *DBX) NamedQueryWrapper(tables []interface{}, query string, arg interface{}) (*sqlx.Rows, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.NamedQuery(query, arg)
}

// SelectWrapper ...
func (db *DBX) SelectWrapper(tables []interface{}, dest interface{}, query string, args ...interface{}) error {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Select(dest, query, args...)
}

// ExecWrapper ...
func (db *DBX) ExecWrapper(tables []interface{}, query string, args ...interface{}) (sql.Result, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Exec(query, args...)
}

// QueryRowxWrapper ...
func (db *DBX) QueryRowxWrapper(tables []interface{}, query string, args ...interface{}) *sqlx.Row {
	query = fmt.Sprintf(query, tables...)
	return db.DB.QueryRowx(query, args...)
}

// QueryxWrapper ...
func (db *DBX) QueryxWrapper(tables []interface{}, query string, args ...interface{}) (*sqlx.Rows, error) {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Queryx(query, args...)
}

// GetWrapper ...
func (db *DBX) GetWrapper(tables []interface{}, dest interface{}, query string, args ...interface{}) error {
	query = fmt.Sprintf(query, tables...)
	return db.DB.Get(dest, query, args...)
}

func buildDBX(oriDB *sql.DB) *DBX {
	return &DBX{
		DB: sqlx.NewDb(oriDB, DB_TYPE_MYSQL),
	}
}

func buildGormDB(oriDB *sql.DB) (*GormDB, error) {
	gormDB, err := gorm.Open("mysql", oriDB)
	if err != nil {
		return nil, err
	}
	return &GormDB{
		DB: gormDB,
	}, err
}
