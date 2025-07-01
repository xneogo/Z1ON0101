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
 @Time    : 2024/10/9 -- 18:00
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: manager_selector.go
*/

package xmanager

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/qiguanzhu/infra/nerv/xlog"
	"github.com/qiguanzhu/infra/pkg/consts"
	"github.com/qiguanzhu/infra/seele/zconfig"
	"github.com/qiguanzhu/infra/seele/zconfig/zobserver"
	"github.com/qiguanzhu/infra/seele/zsql"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

const configTypeKey = "configType"

type configType int

const (
	ConfigFromApollo configType = iota
	ConfigFromEnv
	ConfigFromYaml
)

type dbManagerSelectorConfig struct {
	// 用于兼容新不同的配置读取方式
	ConfigType        configType `properties:"configType"`
	NewManagerInsName string     `properties:"newManagerInsName"`
	NewManagerDbName  string     `properties:"newManagerDbName"`
	SqlDsn            string     `properties:"sqlDsn"`
}

type DbManagerSelectorConfiger struct {
	mu  sync.RWMutex
	cfg *dbManagerSelectorConfig
}

func NewDbManagerSelectorConfiger() *DbManagerSelectorConfiger {
	return &DbManagerSelectorConfiger{}
}

func (m *DbManagerSelectorConfiger) IsInitialized() bool {
	return m.cfg != nil
}

func (m *DbManagerSelectorConfiger) GetManagerInfo() (key1 string, key2 string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg.NewManagerInsName, m.cfg.NewManagerDbName
}

func (m *DbManagerSelectorConfiger) loadConfigFromApollo(ctx context.Context, confCenter zconfig.ConfigCenter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cfg := new(dbManagerSelectorConfig)
	err := confCenter.UnmarshalWithNamespace(ctx, consts.MysqlConfNamespace, cfg)
	if err != nil {
		return err
	}
	m.cfg = cfg
	return nil
}

func (m *DbManagerSelectorConfiger) loadConfigFromEnv(ctx context.Context, _ zconfig.ConfigCenter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cfg := new(dbManagerSelectorConfig)
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
	m.cfg = cfg
	return nil
}

func (m *DbManagerSelectorConfiger) loadConfigFromYaml(ctx context.Context, _ zconfig.ConfigCenter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cfg := new(dbManagerSelectorConfig)
	f, err := os.Open("./configs/config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(cfg); err != nil {
		panic(err)
	}
	m.cfg = cfg
	return nil
}

type DbManagerSelector struct {
	newDbManager zsql.ManagerProxy
	configer     *DbManagerSelectorConfiger
}

func NewDbManagerSelector() (*DbManagerSelector, error) {
	dbMs := &DbManagerSelector{
		newDbManager: NewManager(),
		configer:     NewDbManagerSelectorConfiger(),
	}
	return dbMs, nil
}

// InitDbManagerConf init db dynamic conf
func (m *DbManagerSelector) InitDbManagerConf(ctx context.Context, configCenter zconfig.ConfigCenter) error {
	if configCenter == nil {
		return fmt.Errorf("init xsql conf err: configcenter nil")
	}
	err := m.configer.loadConfigFromApollo(ctx, configCenter)
	if err != nil {
		return err
	}
	err = m.newDbManager.InitConf(ctx, configCenter)
	if err != nil {
		return err
	}
	return nil
}

// ReloadDbManagerConf reload db conf when change
func (m *DbManagerSelector) ReloadDbManagerConf(ctx context.Context, configCenter zconfig.ConfigCenter, event zobserver.ChangeEvent) error {
	fun := "DbManagerSelector -->"
	for k, v := range event.Changes {
		if k == configTypeKey {
			switch v.NewValue {
			case fmt.Sprintf("%d", ConfigFromApollo):
				err := m.configer.loadConfigFromApollo(ctx, configCenter)
				if err != nil {
					return err
				}
			case fmt.Sprintf("%d", ConfigFromEnv):
				err := m.configer.loadConfigFromEnv(ctx, configCenter)
				if err != nil {
					return err
				}
			case fmt.Sprintf("%d", ConfigFromYaml):
				err := m.configer.loadConfigFromYaml(ctx, configCenter)
				if err != nil {
					return err
				}
			}

			break
		}
	}
	err := m.newDbManager.ReloadConf(ctx, configCenter, event)
	if err != nil {
		return err
	}
	xlog.Infof(ctx, "%s reload config success!", fun)
	return nil
}

// GetDBDefault in newManager key1: insName, key2: dbname
func (m *DbManagerSelector) GetDBDefault(ctx context.Context) (db zsql.XDBWrapper, err error) {
	if m.configer == nil || !m.configer.IsInitialized() {
		return nil, fmt.Errorf("the DB config not initialized")
	}

	key1, key2 := m.configer.GetManagerInfo()
	return m.newDbManager.GetDB(ctx, key1, key2)
}
