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
 @Time    : 2024/9/30 -- 16:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: db.go
*/

package xmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Bishoptylaor/go-toolkit/xtime"
	"github.com/qiguanzhu/infra/gehirn/colorlog"
	"github.com/qiguanzhu/infra/nerv/magi/xbreaker"
	"github.com/qiguanzhu/infra/nerv/xtrace"
	"github.com/qiguanzhu/infra/pkg/consts"
	"github.com/qiguanzhu/infra/seele/zsql"
	"github.com/xwb1989/sqlparser"
	"strings"
)

const (
	traceComponent = "xsql"
)

var bCheckTableName = true

// DB 实现了XDB接口，同时可以通过GetTx获取一个Tx句柄并进行提交
type DB struct {
	db      *sql.DB
	cluster string
}

// ExecContext exec insert/update/delete and so on.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fun := "xsql.DB.ExecContext"
	table := db.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(db.cluster, table) {
		colorlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, db.cluster, table)
		return nil, errors.New("sql cause breaker, because too many timeout")
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, db.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res, err := db.db.ExecContext(ctx, query, args...)
	statMetricReqDur(ctx, db.cluster, table, "exec", st.Millisecond())
	// stat breaker
	xbreaker.StatBreaker(db.cluster, table, err)
	statMetricReqErrTotal(ctx, db.cluster, table, "exec", err)
	return res, err
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fun := "xsql.DB.QueryContext"
	table := db.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(db.cluster, table) {
		colorlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, db.cluster, table)
		return nil, errors.New("sql cause breaker, because too many timeout")
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, db.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res, err := db.db.QueryContext(ctx, query, args...)
	statMetricReqDur(ctx, db.cluster, table, "query", st.Millisecond())
	// stat breaker
	xbreaker.StatBreaker(db.cluster, table, err)
	statMetricReqErrTotal(ctx, db.cluster, table, "query", err)
	return res, err
}

// QueryRowContext executes a query that is expected to return at most one row.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	fun := "xsql.DB.QueryRowContext"
	table := db.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(db.cluster, table) {
		colorlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, db.cluster, table)
		return nil
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, db.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res := db.db.QueryRowContext(ctx, query, args...)
	statMetricReqDur(ctx, db.cluster, table, "query row", st.Millisecond())
	return res
}

func (db *DB) GetDb() *sql.DB {
	return db.db
}

func (db *DB) TxExecute(ctx context.Context, exec func(db zsql.XDB) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		colorlog.Errorf(ctx, "start transaction error :  err:%v", err)
		return err
	}
	err = exec(tx)
	if err != nil {
		er := tx.Rollback(ctx)
		if er != nil {
			colorlog.Errorf(ctx, "rollback transaction error :  err:%v", er)
		}
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		colorlog.Errorf(ctx, "commit transaction error :  err:%v", err)
		return err
	}
	return nil
}

// SetSqlDB mock时使用
func (db *DB) SetSqlDB(outDb *sql.DB) {
	db.db = outDb
}

// BeginTx return Tx, wrapper of sql.Tx
func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	var err error
	tx := &Tx{cluster: db.cluster}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, "xsql.DB.Begin")
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.cluster, "")

	st := xtime.NewTimeStat()
	tx.tx, err = db.db.Begin()
	statMetricReqDur(ctx, tx.cluster, tx.cluster, "begin", st.Millisecond())
	statMetricReqErrTotal(ctx, tx.cluster, tx.cluster, "begin", err)
	return tx, err
}

func (db *DB) fetchTableName(query string) (table string) {
	table = extractSQLTableName(query)
	if table != "" {
		return
	}

	if db != nil {
		return db.cluster
	}

	return
}

// New returns an Option
func New(dbName, user, password, host string) *zsql.Option {
	op := &zsql.Option{
		DbName:   dbName,
		User:     user,
		Password: password,
		Host:     host,
	}
	op.Port(consts.DefaultPort).Driver(consts.DefaultDriver)
	return op.Set()
}

func bCloseConn(key string) bool {
	if strings.Contains(key, consts.TimeoutMsecKey) || strings.Contains(key, consts.ReadTimeoutMsecKey) || strings.Contains(key, consts.WriteTimeoutMsecKey) || strings.Contains(key, consts.MaxLifeTimeSecKey) {
		return true
	}

	return false
}

func IsReloadConn(key string) bool {
	if strings.Contains(key, consts.MaxIdleConnsKey) || strings.Contains(key, consts.MaxOpenConnsKey) || strings.Contains(key, consts.MaxLifeTimeSecKey) {
		return true
	}
	return false
}

func injectSQLTraceIDLineComment(ctx context.Context, query string) string {
	var traceID string
	traceID = xtrace.ExtractTraceID(ctx)
	if traceID == "" {
		return query
	}

	return fmt.Sprintf("/*%s*/ %s", traceID, query)
}

func extractSQLTableName(query string) (table string) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}

	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		if tableIdent, ok := node.(sqlparser.TableIdent); ok {
			table = tableIdent.String()
			if table != "" {
				return false, fmt.Errorf("has found")
			}
		}

		return true, nil
	}, stmt)

	return
}
