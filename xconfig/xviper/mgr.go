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
 @Time    : 2024/11/10 -- 17:16
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: mgr.go
*/

package xviper

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

type ViperConfigCenterManagerConf struct {
	Path       string           `json:"path"`
	Addrs      []string         `json:"addrs"`
	SourceType ConfigSourceType `json:"source_type"`
	ConfigType ConfigType       `json:"config_type"`
}
type ViperConfigCenterManager struct {
	conf         ViperConfigCenterManagerConf
	viperClients map[string]*ViperConfigCenter
	mu           sync.Mutex
}

func NewViperConfigCenterManager(conf ViperConfigCenterManagerConf) *ViperConfigCenterManager {
	return &ViperConfigCenterManager{
		conf:         conf,
		viperClients: make(map[string]*ViperConfigCenter),
		mu:           sync.Mutex{},
	}
}

func (m *ViperConfigCenterManager) RegisterObserver(ctx context.Context, namespace string, observer *mobserver.ConfigObserver) func() {
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return nil
	}
	return client.RegisterObserver(ctx, observer)
}
func (m *ViperConfigCenterManager) getViper(ctx context.Context, namespace string) (*ViperConfigCenter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var viperCenter = m.viperClients[namespace]
	if viperCenter == nil {
		if viperCenter0, err := NewViperConfigCenter(ViperConfigCenterConf{
			Namespace:                    namespace,
			ViperConfigCenterManagerConf: m.conf,
		}); err == nil {
			viperCenter = viperCenter0
			m.viperClients[namespace] = viperCenter
		} else {
			return nil, err
		}
	}
	return viperCenter, nil
}

func (m *ViperConfigCenterManager) Stop(ctx context.Context, namespace string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.Stop")
	defer span.Finish()
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return err
	}
	client.Stop(ctx)
	return nil
}

func (m *ViperConfigCenterManager) GetString(ctx context.Context, namespace string, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetString")
	defer span.Finish()
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return "", false
	}
	return client.GetString(ctx, key)
}

func (m *ViperConfigCenterManager) GetBool(ctx context.Context, namespace string, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetBool")
	defer span.Finish()
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return false, false
	}
	return client.GetBool(ctx, key)
}

func (m *ViperConfigCenterManager) GetInt(ctx context.Context, namespace string, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetInt")
	defer span.Finish()

	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return 0, false
	}
	return client.GetInt(ctx, key)
}

func (m *ViperConfigCenterManager) GetFloat64(ctx context.Context, namespace string, key string) (float64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetFloat64")
	defer span.Finish()

	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return 0, false
	}
	return client.GetFloat64(ctx, key)
}

func (m *ViperConfigCenterManager) GetInt64(ctx context.Context, namespace string, key string) (int64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetInt64")
	defer span.Finish()

	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return 0, false
	}
	return client.GetInt64(ctx, key)
}

func (m *ViperConfigCenterManager) GetInt32(ctx context.Context, namespace string, key string) (int32, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetInt32")
	defer span.Finish()

	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return 0, false
	}
	return client.GetInt32(ctx, key)
}
func (m *ViperConfigCenterManager) GetAllKeys(ctx context.Context, namespace string) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "ViperConfigCenterManager.GetAllKeys")
	defer span.Finish()

	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return nil
	}
	return client.GetAllKeys(ctx)
}

func (m *ViperConfigCenterManager) Unmarshal(ctx context.Context, namespace string, v interface{}) error {
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return err
	}
	return client.Unmarshal(ctx, v)
}

func (m *ViperConfigCenterManager) UnmarshalKey(ctx context.Context, namespace string, key string, v interface{}) error {
	client, err := m.getViper(ctx, namespace)
	if err != nil {
		xlog.Errorf(ctx, "get viper err: ns: %s, err: %v", namespace, err)
		return err
	}
	return client.UnmarshalKey(ctx, key, v)
}
