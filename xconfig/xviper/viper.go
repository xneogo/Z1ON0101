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
 @Time    : 2024/11/10 -- 17:15
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: viper.go
*/

package xviper

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

type ConfigSourceType string

const (
	FILE   ConfigSourceType = "file"
	ETCD   ConfigSourceType = "etcd"
	ETCD3  ConfigSourceType = "etcd3"
	APPLLO ConfigSourceType = "appllo"
	REDIS  ConfigSourceType = "redis"
)

type ConfigType string

const (
	YAML ConfigType = "yaml"
	JSON ConfigType = "json"
)

type ViperConfigCenterConf struct {
	Namespace string
	ViperConfigCenterManagerConf
}
type ViperConfigCenter struct {
	Conf ViperConfigCenterConf

	viperClient *viper.Viper
	observers   []*mobserver.ConfigObserver
	mu          sync.Mutex
}

func NewViperConfigCenter(conf ViperConfigCenterConf) (*ViperConfigCenter, error) {
	v := viper.New()
	confCenter := &ViperConfigCenter{
		Conf:        conf,
		viperClient: v,
		observers:   make([]*mobserver.ConfigObserver, 0),
		mu:          sync.Mutex{},
	}
	v.AutomaticEnv()
	v.SetConfigName(conf.Namespace)
	v.SetConfigType(string(conf.ConfigType))
	switch conf.SourceType {
	case FILE:
		var ext = ""
		switch conf.ConfigType {
		case YAML:
			ext = ".yaml"
		case JSON:
			ext = ".json"
		}
		v.AddConfigPath(conf.Path) // or other path where your yaml files are located
		v.SetConfigFile(fmt.Sprintf("%s/%s%s", conf.Path, conf.Namespace, ext))
		v.ReadInConfig()
		// Register a callback function to be called when the configuration changes.
		// viper的watch机制只支持本地文件
		v.OnConfigChange(confCenter.notify)
		v.WatchConfig()
	case ETCD:
		var location = fmt.Sprintf("%s/%s", conf.Path, conf.Namespace)
		v.AddRemoteProvider(string(conf.SourceType), strings.Join(conf.Addrs, ";"), location)
		v.SetConfigType(string(conf.ConfigType)) // assuming that etcd stores json data
		v.ReadRemoteConfig()
		fmt.Println(viper.RemoteConfig)
		// viper没有实现变更通知机制，这里需要再创建一个watch，一直监听变更，有变更了，通知出去。
		// 监听配置变更
		go func() {
			for {
				v.WatchRemoteConfigOnChannel()
				confCenter.notify(fsnotify.Event{
					Name: "change",
					Op:   0,
				})
			}

		}()
	default:
		return nil, fmt.Errorf("not support source type: %s", conf.SourceType)
	}
	return confCenter, nil
}

func (m *ViperConfigCenter) RegisterObserver(ctx context.Context, observer *mobserver.ConfigObserver) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 注册时 启动监听
	observer.StartWatch(ctx)
	m.observers = append(m.observers, observer)
	return func() {
	}
}

// todo viper针对文件方式没有每个key的变化，每个key的变化都会通知，并没有oldvalue
func (m *ViperConfigCenter) notify(e fsnotify.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// viper是真个文件维度的变更通知，因此需要手动对比不同key的变化
	keys := m.viperClient.AllKeys()
	var valMap = make(map[string]string)
	for _, key := range keys {
		valMap[key] = m.viperClient.GetString(key)
	}
	var changes = make(map[string]*mobserver.Change)
	for key, oval := range valMap {
		var change = &mobserver.Change{
			NewValue:   oval,
			ChangeType: mobserver.MODIFY,
		}
		changes[key] = change
	}
	for _, observer := range m.observers {
		observer.HandleChangeEvent(&mobserver.ChangeEvent{
			Namespace: m.Conf.Namespace,
			Changes:   changes,
		})
	}
	return
}

func (m *ViperConfigCenter) Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.Stop")
	defer span.Finish()
	return nil
}

func (m *ViperConfigCenter) GetString(ctx context.Context, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetString")
	defer span.Finish()
	if ok := m.viperClient.IsSet(key); !ok {
		return "", ok
	}
	return m.viperClient.GetString(key), true
}

func (m *ViperConfigCenter) GetBool(ctx context.Context, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetBool")
	defer span.Finish()
	if ok := m.viperClient.IsSet(key); !ok {
		return false, ok
	}
	return m.viperClient.GetBool(key), true
}

func (m *ViperConfigCenter) GetInt(ctx context.Context, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetInt")
	defer span.Finish()

	if ok := m.viperClient.IsSet(key); !ok {
		return 0, ok
	}
	return m.viperClient.GetInt(key), true
}

func (m *ViperConfigCenter) GetFloat64(ctx context.Context, key string) (float64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetFloat64")
	defer span.Finish()

	if ok := m.viperClient.IsSet(key); !ok {
		return 0.0, ok
	}
	return m.viperClient.GetFloat64(key), true
}

func (m *ViperConfigCenter) GetInt64(ctx context.Context, key string) (int64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetInt64")
	defer span.Finish()

	if ok := m.viperClient.IsSet(key); !ok {
		return 0, ok
	}
	return m.viperClient.GetInt64(key), true
}

func (m *ViperConfigCenter) GetInt32(ctx context.Context, key string) (int32, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetInt32")
	defer span.Finish()

	if ok := m.viperClient.IsSet(key); !ok {
		return 0, ok
	}
	return m.viperClient.GetInt32(key), true
}
func (m *ViperConfigCenter) GetAllKeys(ctx context.Context) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenter.GetAllKeys")
	defer span.Finish()
	return m.viperClient.AllKeys()
}

func (m *ViperConfigCenter) Unmarshal(ctx context.Context, v interface{}) error {
	return m.viperClient.Unmarshal(v, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
}

func (m *ViperConfigCenter) UnmarshalKey(ctx context.Context, key string, v interface{}) error {
	return m.viperClient.UnmarshalKey(key, v, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
}
