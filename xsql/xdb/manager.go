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
 @Time    : 2024/9/30 -- 15:45
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: manager.go
*/

package xmanager

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/qiguanzhu/infra/nerv/magi/xreflect"
	"github.com/qiguanzhu/infra/nerv/xlog"
	"github.com/qiguanzhu/infra/nerv/xtrace"
	"github.com/qiguanzhu/infra/pkg/consts"
	"github.com/qiguanzhu/infra/seele/zconfig"
	"github.com/qiguanzhu/infra/seele/zconfig/zobserver"
	"github.com/qiguanzhu/infra/seele/zsql"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DynamicConfigure ...
type DynamicConfigure struct {
	mysqlConf     *zsql.Cfg
	mu            sync.RWMutex
	closeInsChan  chan zsql.ChangeIns
	reloadInsChan chan zsql.ChangeIns
}

// LoadDynamicConf ...
func (c *DynamicConfigure) LoadDynamicConf(insName string, dynamicConf *zsql.DynamicConf) {
	config := &zsql.DynamicConf{
		Timeout:        consts.DefaultTimeoutSecond * time.Second,
		ReadTimeout:    consts.DefaultReadTimeoutSecond * time.Second,
		WriteTimeout:   consts.DefaultWriteTimeoutSecond * time.Second,
		MaxLifeTimeSec: consts.DefaultMaxLifeTimeSecond * time.Second,
		MaxIdleConns:   consts.DefaultMaxIdleConns,
		MaxOpenConns:   consts.DefaultMaxOpenConns,
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.mysqlConf == nil {
		_ = xreflect.Mirroring(config, dynamicConf)
		return
	}
	if v, ok := c.mysqlConf.ConfMap[insName]; ok {
		if v.MaxIdleConns != 0 {
			config.MaxIdleConns = v.MaxIdleConns
		}
		if v.MaxOpenConns != 0 {
			config.MaxOpenConns = v.MaxOpenConns
		}
		if v.MaxLifeTimeSec != 0 {
			config.MaxLifeTimeSec = time.Duration(v.MaxLifeTimeSec) * time.Second
		}
		if v.TimeoutMsec != 0 {
			config.Timeout = time.Duration(v.TimeoutMsec) * time.Millisecond
		}
		if v.ReadTimeoutMsec != 0 {
			config.ReadTimeout = time.Duration(v.ReadTimeoutMsec) * time.Millisecond
		}
		if v.WriteTimeoutMsec != 0 {
			config.WriteTimeout = time.Duration(v.WriteTimeoutMsec) * time.Millisecond
		}
		if v.Username != "" {
			config.Username = v.Username
		}
		if v.Password != "" {
			config.Password = v.Password
		}
	}
	_ = xreflect.Mirroring(config, dynamicConf)
	return
}

// SetConf ...
func (c *DynamicConfigure) SetConf(cfg *zsql.Cfg) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mysqlConf = cfg
}

func (c *DynamicConfigure) IsInit() bool {
	return c.mysqlConf != nil
}

type InstanceKey struct {
	instanceName string
	dbName       string
}

func CreateInstanceKey(instanceName, dbName string) (InstanceKey, error) {
	if instanceName == "" || dbName == "" {
		return InstanceKey{}, fmt.Errorf("nil params")
	}
	return InstanceKey{
		instanceName: instanceName,
		dbName:       dbName,
	}, nil
}

func (m *InstanceKey) String() string {
	return fmt.Sprintf("%s-%s", m.instanceName, m.dbName)
}

func (m *InstanceKey) GetInstanceName() string {
	return m.instanceName
}

func (m *InstanceKey) GetDbName() string {
	return m.dbName
}

// ProxyAddr proxy地址
type ProxyAddr struct {
	Host string
	Port int
}

type Manager struct {
	// 服务对应的 proxy 地址, 优先级顺序为: ENV > apollo配置 > default
	proxyAddr        ProxyAddr
	instanceMu       sync.RWMutex
	instances        map[InstanceKey]zsql.DBInstanceProxy
	dynamicConfigure zsql.DynamicConfigureProxy[*zsql.Cfg, *zsql.DynamicConf]
}

// NewManager 实例化DB路由对象
func NewManager() zsql.ManagerProxy {
	proxyAddr := getWeirAddrFromEnvOrDefault()

	return &Manager{
		proxyAddr:        proxyAddr,
		instances:        make(map[InstanceKey]zsql.DBInstanceProxy),
		dynamicConfigure: &DynamicConfigure{},
	}
}

func (m *Manager) InitConf(ctx context.Context, confCenter zconfig.ConfigCenter) error {
	f := "Manager.InitConf ->"

	if confCenter == nil {
		return fmt.Errorf("init xsql conf err: configcenter nil")
	}
	mysqlConf := new(zsql.Cfg)
	err := confCenter.UnmarshalWithNamespace(ctx, consts.MysqlConfNamespace, mysqlConf)
	if err != nil {
		return err
	}
	m.dynamicConfigure.SetConf(mysqlConf)

	if modified := m.modifyProxyAddrIfNeeded(mysqlConf); modified {
		xlog.Infof(ctx, "%s proxy addr modified from dynamic conf, addr: %v", f, m.proxyAddr)
	}
	return nil
}

// GetDB return xmanager.DB without transaction
func (m *Manager) GetDB(ctx context.Context, insName, dbName string) (zsql.XDBWrapper, error) {
	if !m.dynamicConfigure.IsInit() {
		return nil, fmt.Errorf("dynamic configer not init")
	}
	ins, err := m.GetInstance(insName, dbName)
	if err != nil {
		return nil, err
	}
	db := new(DB)
	db.cluster = insName
	// db.table = dbName
	db.db = ins.GetDB()
	return db, nil
}

func (m *Manager) ReloadConf(ctx context.Context, config zconfig.ConfigCenter, event zobserver.ChangeEvent) error {
	f := "Manager.ReloadConf -->"
	if event.Namespace != consts.MysqlConfNamespace {
		return nil
	}
	if config == nil {
		return fmt.Errorf("reload xsql conf err: configcenter nil")
	}
	c := new(zsql.Cfg)
	err := config.UnmarshalWithNamespace(ctx, consts.MysqlConfNamespace, c)
	if err != nil {
		return err
	}
	m.dynamicConfigure.SetConf(c)

	if modified := m.modifyProxyAddrIfNeeded(c); modified {
		xlog.Infof(ctx, "%s proxy addr is modified from apollo, async close all db instances", f)
		m.asyncCloseAllDbInstance(ctx)
		return nil
	}

	closeInsMap, reloadInsMap, err := getCloseAndReloadInsMap(event)
	if err != nil {
		return err
	}
	m.closeDbInstance(ctx, closeInsMap)
	m.reloadDbInstance(ctx, reloadInsMap)
	return nil
}

func (m *Manager) GetInstance(insName, dbName string) (zsql.DBInstanceProxy, error) {
	fun := "ManagerV2 -->"
	m.instanceMu.RLock()
	insKey, err := CreateInstanceKey(insName, dbName)
	if err != nil {
		return nil, err
	}
	if in, ok := m.instances[insKey]; ok {
		m.instanceMu.RUnlock()
		return in, nil
	}
	m.instanceMu.RUnlock()
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()
	if in, ok := m.instances[insKey]; ok {
		return in, nil
	}
	newIns, err := NewDBInstance(m.proxyAddr, &insKey, m.dynamicConfigure)
	if err != nil {
		return nil, errors.Wrapf(err, "%s buildInstance error!", fun)
	}
	m.instances[insKey] = newIns
	return newIns, nil
}

func (m *Manager) modifyProxyAddrIfNeeded(mysqlConf zsql.SqlConfigProxy) (modified bool) {
	if !isWeirProxyHostEnvSet() {
		if mysqlConf.IsProxyHostSet() && m.proxyAddr.Host != mysqlConf.GetProxyHost() {
			m.setProxyHost(mysqlConf.GetProxyHost())
			modified = true
		}
	}

	if !isWeirProxyPortEnvSet() {
		if mysqlConf.IsProxyPortSet() && m.proxyAddr.Port != mysqlConf.GetProxyPort() {
			m.setProxyPort(mysqlConf.GetProxyPort())
			modified = true
		}
	}

	return modified
}

func (m *Manager) setProxyHost(host string) {
	m.proxyAddr.Host = host
}

func (m *Manager) setProxyPort(port int) {
	m.proxyAddr.Port = port
}

func getCloseAndReloadInsMap(event zobserver.ChangeEvent) (closeInsMap, reloadInsMap map[string]struct{}, err error) {
	insMap := make(map[string]struct{})
	closeInsMap = make(map[string]struct{})
	reloadInsMap = make(map[string]struct{})
	for k, v := range event.Changes {
		if v != nil {
			parts := strings.Split(k, consts.KeySep)
			if len(parts) < 3 {
				continue
			}
			insName := parts[1]
			if _, ok := insMap[insName]; ok {
				continue
			}
			insMap[insName] = struct{}{}
			if !IsReloadConn(k) {
				closeInsMap[insName] = struct{}{}
			} else {
				reloadInsMap[insName] = struct{}{}
			}
		}
	}
	return
}

func (m *Manager) closeDbInstance(ctx context.Context, insNameMap map[string]struct{}) {
	fun := "Manager.closeDbInstance -->"
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()
	for key, ins := range m.instances {
		if _, ok := insNameMap[key.GetInstanceName()]; ok {
			delete(m.instances, key)
			go func(insKey InstanceKey, dbIns zsql.DBInstanceProxy) {
				if err := dbIns.Close(); err == nil {
					xlog.Infof(ctx, "%s succeed to close db, instance: %v", fun, insKey)
				} else {
					xlog.Errorf(ctx, "%s fail close db, instance: %v, error: %v", fun, insKey, err)
				}
			}(key, ins)
		}
	}
}

func (m *Manager) asyncCloseAllDbInstance(ctx context.Context) {
	fun := "Manager.asyncCloseAllDbInstance -->"
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()

	for key, ins := range m.instances {
		delete(m.instances, key)
		go func(insKey InstanceKey, dbIns zsql.DBInstanceProxy) {
			if err := dbIns.Close(); err == nil {
				xlog.Infof(ctx, "%s succeed to close db, instance: %v", fun, insKey)
			} else {
				xlog.Errorf(ctx, "%s fail to close db, instance: %v, error: %v", fun, insKey, err)
			}
		}(key, ins)
	}
}

// TODO://配置修改不是原子操作。读取到的配置可能有中间状态
func (m *Manager) reloadDbInstance(ctx context.Context, insNameMap map[string]struct{}) {
	fun := "Manager.reloadDbInstance -->"
	m.instanceMu.Lock()
	defer m.instanceMu.Unlock()
	for key, ins := range m.instances {
		if _, ok := insNameMap[key.GetInstanceName()]; ok {
			go func(insKey InstanceKey, dbIns zsql.DBInstanceProxy) {
				if err := dbIns.Reload(); err == nil {
					xlog.Infof(ctx, "%s succeed to reload db, instance: %v, dbName: %s", fun, insKey, dbIns.GetDbName())
				} else {
					xlog.Errorf(ctx, "%s fail to reload db, instance: %v, dbName: %s, error: %v", fun, insKey, dbIns.GetDbName(), err)
				}
			}(key, ins)
		}
	}
}

func getWeirAddrFromEnvOrDefault() ProxyAddr {
	var host string
	var port int

	envHost := os.Getenv(consts.WeirProxyHostEnv)
	if envHost != "" {
		host = envHost
	} else {
		host = consts.WeirProxyHost
	}

	envPortStr := os.Getenv(consts.WeirProxyPortEnv)
	if envPort, err := strconv.Atoi(envPortStr); err == nil {
		port = envPort
	} else {
		port = consts.WeirProxyPort
	}

	return ProxyAddr{Host: host, Port: port}
}

func isWeirProxyHostEnvSet() bool {
	return os.Getenv(consts.WeirProxyHostEnv) != ""
}

func isWeirProxyPortEnvSet() bool {
	return os.Getenv(consts.WeirProxyPortEnv) != ""
}

func setDBSpanTags(span opentracing.Span, cluster, table, stmt string) {
	span.SetTag(xtrace.TagComponent, traceComponent)
	span.SetTag(xtrace.TagDBType, xtrace.DBTypeSQL)
	span.SetTag(xtrace.TagSpanKind, xtrace.SpanKindClient)
	setTagIfNonEmpty(span, xtrace.TagPalfishDBCluster, cluster)
	setTagIfNonEmpty(span, xtrace.TagPalfishDBTable, table)
	setTagIfNonEmpty(span, xtrace.TagDBSQLTable, table)
	setTagIfNonEmpty(span, xtrace.TagDBStatement, stmt)
}

func setTagIfNonEmpty(span opentracing.Span, key, val string) {
	if val != "" {
		span.SetTag(key, val)
	}
}
