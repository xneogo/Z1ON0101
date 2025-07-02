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
 @Time    : 2024/11/4 -- 10:31
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xdb_tx.go
*/

package xmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/Z1ON0101/xtrace"
	"github.com/xneogo/extensions/xbreaker"
	"github.com/xneogo/extensions/xtime"
)

// Tx wrapper of sql.Tx
type Tx struct {
	tx      *sql.Tx
	cluster string
}

// ExecContext exec insert/update/delete and so on.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	fun := "xsql.Tx.ExecContext"
	table := tx.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(tx.cluster, table) {
		xlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, tx.cluster, table)
		return nil, errors.New("sql cause breaker, because too many timeout")
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, tx.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res, err := tx.tx.ExecContext(ctx, query, args...)
	statMetricReqDur(ctx, tx.cluster, table, "exec", st.Millisecond())
	// stat breaker
	xbreaker.StatBreaker(tx.cluster, table, err)
	statMetricReqErrTotal(ctx, tx.cluster, table, "exec", err)
	return res, err
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	fun := "xsql.Tx.QueryContext"
	table := tx.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(tx.cluster, table) {
		xlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, tx.cluster, table)
		return nil, errors.New("sql cause breaker, because too many timeout")
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, tx.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res, err := tx.tx.QueryContext(ctx, query, args...)
	statMetricReqDur(ctx, tx.cluster, table, "query", st.Millisecond())
	// stat breaker
	xbreaker.StatBreaker(tx.cluster, table, err)
	statMetricReqErrTotal(ctx, tx.cluster, table, "query", err)
	return res, err
}

// QueryRowContext executes a query that is expected to return at most one row.
func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	fun := "xsql.Tx.QueryRowContext"
	table := tx.fetchTableName(query)
	// check breaker
	if !xbreaker.Entry(tx.cluster, table) {
		xlog.Errorf(ctx, "%s trigger tidb breaker, because too many timeout sqls, cluster: %s, table: %s", fun, tx.cluster, table)
		return nil
	}
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, fun)
	defer span.Finish()
	query = injectSQLTraceIDLineComment(ctx, query)
	setDBSpanTags(span, tx.cluster, table, fmt.Sprintf("%s %v", query, args))

	st := xtime.NewTimeStat()
	res := tx.tx.QueryRowContext(ctx, query, args...)
	statMetricReqDur(ctx, tx.cluster, table, "query row", st.Millisecond())
	return res
}

// Commit wrapper of sql.Tx commit
func (tx *Tx) Commit(ctx context.Context) error {
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, "xsql.Tx.Commit")
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.cluster, "")

	st := xtime.NewTimeStat()
	err := tx.tx.Commit()
	statMetricReqDur(ctx, tx.cluster, tx.cluster, "commit", st.Millisecond())
	statMetricReqErrTotal(ctx, tx.cluster, tx.cluster, "commit", err)
	return err
}

// Rollback wrapper of sql.Tx rollback
func (tx *Tx) Rollback(ctx context.Context) error {
	// trace
	span, ctx := xtrace.StartSpanFromContext(ctx, "xsql.Tx.Commit")
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.cluster, "")

	st := xtime.NewTimeStat()
	err := tx.tx.Rollback()
	statMetricReqDur(ctx, tx.cluster, tx.cluster, "rollback", st.Millisecond())
	statMetricReqErrTotal(ctx, tx.cluster, tx.cluster, "rollback", err)
	return err
}

func (tx *Tx) fetchTableName(query string) (table string) {
	table = extractSQLTableName(query)
	if table != "" {
		return
	}

	if tx != nil {
		return tx.cluster
	}

	return
}
