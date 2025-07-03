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
 @Time    : 2024/11/10 -- 17:20
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: yaml.go
*/

package xyaml

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/xneogo/eins/oneenv"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

const (
	defaultCluster              = "default"
	defaultCacheDir             = "/tmp/sconfcenter"
	defaultCfgDir               = "/opt/conf"
	defaultNamespaceApplication = "application"
)

type YamlConfigCenter struct {
	servLoc           string
	cfgdir            string
	mu                sync.Mutex
	observers         []*mobserver.ConfigObserver
	recalledObservers map[*mobserver.ConfigObserver]interface{}
	configFiles       map[string]string
	fileDatas         map[string][]byte
	parsers           map[string]*viper.Viper
}

type YamlConfigCenterConf struct {
	ServLoc        string   `json:"serv_loc"`
	Cfgdir         string   `json:"cfgdir"`
	NameSpaceNames []string `json:"namespaceNames,omitempty"`
}

func NewYamlConfigCenter(ctx context.Context, conf *YamlConfigCenterConf) (*YamlConfigCenter, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "apolloDriver.New")
	defer span.Finish()
	cfgDir := conf.Cfgdir
	if cfgDir == "" {
		cfgDir = oneenv.GetEnvWithDefault("CFGDIR", defaultCfgDir)
	}
	center := &YamlConfigCenter{
		cfgdir:            fmt.Sprintf("%s/%s", cfgDir, conf.ServLoc),
		observers:         make([]*mobserver.ConfigObserver, 0),
		servLoc:           conf.ServLoc,
		configFiles:       make(map[string]string),
		fileDatas:         make(map[string][]byte),
		parsers:           make(map[string]*viper.Viper),
		recalledObservers: make(map[*mobserver.ConfigObserver]interface{}),
	}
	namespaceMap := make(map[string]interface{})
	for _, ns := range conf.NameSpaceNames {
		namespaceMap[strings.ToLower(ns)] = struct {
		}{}
	}
	filepath.Walk(center.cfgdir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && namespaceMap[strings.TrimSuffix(strings.ToLower(info.Name()), ".yaml")] != nil {
			namespace := fileNameToNamespace(info.Name())
			center.configFiles[namespace] = path
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			center.fileDatas[namespace] = data
			parser := viper.New()
			parser.SetConfigFile(path)
			parser.SetConfigType("yml")
			parser.ReadInConfig()
			center.parsers[namespace] = parser
		}
		return nil
	})
	center.StartWatchUpdate()

	return center, nil
}

func (m *YamlConfigCenter) StartWatchUpdate() {
	go func() {
		m.watchAllFiles()
	}()
}
func fileNameToNamespace(filename string) string {
	namespace := strings.ReplaceAll(strings.ToLower(filename), ".yaml", "")
	return namespace
}
func (m *YamlConfigCenter) watchAllFiles() {
	for k, parser := range m.parsers {
		parser.OnConfigChange(func(in fsnotify.Event) {
			var chgType mobserver.ChangeType
			switch in.Op {
			case fsnotify.Create:
				chgType = mobserver.ADD
			case fsnotify.Write:
				chgType = mobserver.MODIFY
			case fsnotify.Remove:
				chgType = mobserver.DELETE
			default:
				return
			}
			newData, _ := os.ReadFile(in.Name)
			changEvent := &mobserver.ChangeEvent{
				Namespace: k,
				Changes: map[string]*mobserver.Change{
					"_": &mobserver.Change{
						OldValue:   string(m.fileDatas[k]),
						NewValue:   string(newData),
						ChangeType: chgType,
					},
				},
			}
			if changEvent != nil {
				for _, ob := range m.observers {
					go func() {
						ob.HandleChangeEvent(changEvent)
					}()
				}
			}
			return
		})
		parser.WatchConfig()
	}

}

func (m *YamlConfigCenter) RegisterObserver(ctx context.Context, namespace string, observer *mobserver.ConfigObserver) func() {
	// 注册时 启动监听
	observer.StartWatch(ctx)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, observer)
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.recalledObservers[observer] = struct {
		}{}
	}
}

func (m *YamlConfigCenter) Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.Stop")
	defer span.Finish()
	return nil
}

func (m *YamlConfigCenter) SubscribeNamespaces(ctx context.Context, namespaceNames []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.SubscribeNamespaces")
	defer span.Finish()
	return nil
}

func (m *YamlConfigCenter) GetString(ctx context.Context, namespace, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetString")
	defer span.Finish()
	m.mu.Lock()
	defer m.mu.Unlock()
	if parser, ok := m.parsers[namespace]; ok {
		return parser.GetString(key), true
	}
	return "", false
}

func (m *YamlConfigCenter) GetBool(ctx context.Context, namespace, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetBool")
	defer span.Finish()

	m.mu.Lock()
	defer m.mu.Unlock()

	if parser, ok := m.parsers[namespace]; ok {
		return parser.GetBool(key), true
	}
	return false, false
}

func (m *YamlConfigCenter) GetInt(ctx context.Context, namespace, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetInt")
	defer span.Finish()
	i, ok := m.GetInt64(ctx, namespace, key)
	return int(i), ok
}

func (m *YamlConfigCenter) GetFloat64(ctx context.Context, namespace, key string) (float64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetFloat64")
	defer span.Finish()

	m.mu.Lock()
	defer m.mu.Unlock()
	if parser, ok := m.parsers[namespace]; ok {
		return parser.GetFloat64(key), true
	}
	return 0, false
}

func (m *YamlConfigCenter) GetInt64(ctx context.Context, namespace, key string) (int64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetInt64")
	defer span.Finish()
	m.mu.Lock()
	defer m.mu.Unlock()

	if parser, ok := m.parsers[namespace]; ok {
		return parser.GetInt64(key), true
	}
	return 0, false
}

func (m *YamlConfigCenter) GetInt32(ctx context.Context, namespace, key string) (int32, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetInt32")
	defer span.Finish()
	i, ok := m.GetInt64(ctx, namespace, key)
	return int32(i), ok
}

func (m *YamlConfigCenter) GetAllKeys(ctx context.Context, namespace string) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.GetAllKeys")
	defer span.Finish()
	m.mu.Lock()
	defer m.mu.Unlock()
	var keys []string
	if parser, ok := m.parsers[namespace]; ok {
		return parser.AllKeys()
	}
	return keys
}

func (m *YamlConfigCenter) Unmarshal(ctx context.Context, namespace string, v interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.Unmarshal")
	defer span.Finish()
	m.mu.Lock()
	defer m.mu.Unlock()
	if parser, ok := m.parsers[namespace]; ok {
		return parser.Unmarshal(v)
	}
	return fmt.Errorf("no namespace: %s", namespace)
}

func (m *YamlConfigCenter) UnmarshalKey(ctx context.Context, namespace string, key string, v interface{}) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "YamlConfigCenter.Unmarshal")
	defer span.Finish()
	m.mu.Lock()
	defer m.mu.Unlock()
	if parser, ok := m.parsers[namespace]; ok {
		return parser.UnmarshalKey(key, v)
	}
	return fmt.Errorf("no namespace: %s", namespace)
}
