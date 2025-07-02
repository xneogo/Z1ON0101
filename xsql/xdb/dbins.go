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
 @Time    : 2024/10/25 -- 17:54
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: dbins.go
*/

package xmanager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/msql"
)

// DBInstance ...
type DBInstance struct {
	insName          string
	group            string
	dbType           string
	dbName           string
	dbAddr           string
	timeOut          time.Duration
	userName         string
	passWord         string
	db               *sql.DB
	dynamicConfigure msql.DynamicConfigureProxy[*msql.Cfg, *msql.DynamicConf]
}

// GetType ...
func (m *DBInstance) GetType() string {
	return m.dbType
}

// Close ...
func (m *DBInstance) Close() error {
	err := m.db.Close()
	if err != nil {
		return fmt.Errorf("xsql Instance Close err: %v", err)
	}

	return nil
}

// Reload ...
func (m *DBInstance) Reload() error {
	dynamicConf := new(msql.DynamicConf)
	m.dynamicConfigure.LoadDynamicConf(m.insName, dynamicConf)
	m.db.SetMaxIdleConns(dynamicConf.MaxIdleConns)
	m.db.SetMaxOpenConns(dynamicConf.MaxOpenConns)
	m.db.SetConnMaxLifetime(dynamicConf.MaxLifeTimeSec)

	return nil
}

func (m *DBInstance) GetDbName() string {
	return m.dbName
}
func (m *DBInstance) GetDB() *sql.DB {
	return m.db
}

func concatDSN(settings []msql.Setting) string {
	s := ""
	for _, f := range settings {
		s = f(s)
	}
	return strings.TrimRight(s, "&")
}

func realDSN(info *msql.Option) string {
	format := "%s:%s@tcp(%s:%d)/%s?%s"
	return strings.TrimRight(fmt.Sprintf(format,
		info.GetUser(), info.GetPassword(), info.GetHost(), info.GetPort(), info.GetDbName(), concatDSN(info.GetSettings())), "?")
}

func open(o *msql.Option) (*sql.DB, error) {
	return sql.Open(o.GetDriver(), realDSN(o))
}

func openFromDsn(o *msql.Option, dsn string) (*sql.DB, error) {
	return sql.Open(o.GetDriver(), dsn)
}

// NewDBInstance 实例化DB实例
func NewDBInstance(addr ProxyAddr, insKey *InstanceKey, dynamicConfigure msql.DynamicConfigureProxy[*msql.Cfg, *msql.DynamicConf]) (*DBInstance, error) {
	dynamicConf := new(msql.DynamicConf)
	dynamicConfigure.LoadDynamicConf(insKey.GetInstanceName(), dynamicConf)
	var err error
	if dynamicConf == nil {
		return nil, errors.New("dynamic conf is nil")
	}
	db, err := New(insKey.GetDbName(), dynamicConf.Username, dynamicConf.Password, addr.Host).Set(GetSettingFunctionList(dynamicConf)...).Port(addr.Port).Open(true, open)
	if err != nil {
		xlog.Errorf(context.TODO(), "new db instance error:%+v", err.Error())
		return nil, err
	}
	db.SetMaxIdleConns(dynamicConf.MaxIdleConns)
	db.SetMaxOpenConns(dynamicConf.MaxOpenConns)
	db.SetConnMaxLifetime(dynamicConf.MaxLifeTimeSec)
	instance := &DBInstance{
		insName:          insKey.instanceName,
		group:            "",
		dbType:           msql.DefaultDbType,
		dbName:           insKey.dbName,
		dbAddr:           addr.Host,
		userName:         dynamicConf.Username,
		passWord:         dynamicConf.Password,
		timeOut:          dynamicConf.Timeout,
		dynamicConfigure: dynamicConfigure,
		db:               db,
	}
	return instance, err
}
