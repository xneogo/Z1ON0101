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
 @Time    : 2024/10/12 -- 15:47
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: apollo.go
*/

package xapollo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eva-nigouki/apogo"
	"github.com/opentracing/opentracing-go"
	"github.com/xneogo/Z1ON0101/xconfig"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/mconfig"
	"github.com/xneogo/matrix/mconfig/mobserver"
	"github.com/xneogo/matrix/mmarshal"
)

const (
	envApolloCluster            = "APOLLO_CLUSTER"
	envApolloHostPort           = "APOLLO_HOST_PORT"
	defaultCluster              = "default"
	defaultHostPort             = "apollo-meta.xxxxx.com:30002"
	defaultCacheDir             = "/tmp/sconfcenter"
	defaultNamespaceApplication = "application"
)

type apolloDriver struct{}

type apolloConfigCenter struct {
	conf *apogo.Conf
	ag   *apogo.Apogo
}

func init() {
	xconfig.Register(ConfigTypeApollo, &apolloDriver{})
}

// SetCluster set apollo center cluster, default value "default"
func SetCluster(cluster string) mconfig.Option {
	return func(center mconfig.ConfigCenter) {
		if apCenter, ok := center.(*apolloConfigCenter); ok {
			apCenter.SetCluster(cluster)
		} else {
			fmt.Printf("SetCluster assert apollo center type err\n")
		}
	}
}

// SetCacheDir set apollo cache dir, default value "/tmp/sconfcenter"
func SetCacheDir(cacheDir string) mconfig.Option {
	return func(center mconfig.ConfigCenter) {
		if apCenter, ok := center.(*apolloConfigCenter); ok {
			apCenter.SetCacheDir(cacheDir)
		} else {
			fmt.Printf("SetCacheDir assert apollo center type err\n")
		}
	}
}

// SetIPHost set apollo remote host, default value "apollo-meta.xxxxx.com:30002"
func SetIPHost(ipHost string) mconfig.Option {
	return func(center mconfig.ConfigCenter) {
		if apCenter, ok := center.(*apolloConfigCenter); ok {
			apCenter.SetIPHost(ipHost)
		} else {
			fmt.Printf("SetIPHost assert apollo center type err\n")
		}
	}
}

// New return apollo config center
func (driver *apolloDriver) New(ctx context.Context, serviceName string, namespaceNames []string, opts ...mconfig.Option) (mconfig.ConfigCenter, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "apolloDriver.New")
	defer span.Finish()
	fun := "apolloDriver.New-->"
	apogo.SetLogger(xlog.GetInfoLogger())
	center := &apolloConfigCenter{
		conf: confFromEnv(),
	}
	for _, opt := range opts {
		opt(center)
	}

	center.conf.AppID = normalizeServiceName(serviceName)
	if len(namespaceNames) > 0 {
		center.conf.NameSpaceNames = namespaceNames
	} else {
		center.conf.NameSpaceNames = []string{defaultNamespaceApplication}
	}
	for i, namespace := range center.conf.NameSpaceNames {
		center.conf.NameSpaceNames[i] = normalizeServiceName(namespace)
	}

	center.ag = apogo.NewApogo(center.conf)

	xlog.Infof(ctx, "%s start agollo with conf:%v", fun, center.conf)

	if err := center.ag.Start(); err != nil {
		xlog.Errorf(ctx, "%s agollo starts err:%v", fun, err)
		return nil, err
	} else {
		xlog.Infof(ctx, "%s agollo starts succeed:%v", fun, err)
	}
	center.ag.StartWatchUpdate()

	return center, nil
}

func (ap *apolloConfigCenter) HandleChangeEvent(event *apogo.ChangeEvent) {
	// TODO implement me
	panic("implement me")
}

func (ap *apolloConfigCenter) RegisterObserver(ctx context.Context, observer *mobserver.ConfigObserver) func() {
	// 注册时 启动监听
	observer.StartWatch(ctx)
	return ap.ag.RegisterObserver(&aObserver{
		observer,
	})
}

func (ap *apolloConfigCenter) Stop(ctx context.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.Stop")
	defer span.Finish()
	return ap.ag.Stop()
}

func (ap *apolloConfigCenter) SubscribeNamespaces(ctx context.Context, namespaceNames []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.SubscribeNamespaces")
	defer span.Finish()
	return ap.ag.SubscribeToNamespaces(namespaceNames...)
}

func (ap *apolloConfigCenter) GetString(ctx context.Context, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetString")
	defer span.Finish()
	return ap.ag.GetString(key)
}

func (ap *apolloConfigCenter) GetStringWithNamespace(ctx context.Context, namespace, key string) (string, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetStringWithNamespace")
	defer span.Finish()

	return ap.ag.GetStringWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetBool(ctx context.Context, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetBool")
	defer span.Finish()

	return ap.ag.GetBool(key)
}

func (ap *apolloConfigCenter) GetBoolWithNamespace(ctx context.Context, namespace, key string) (bool, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetBoolWithNamespace")
	defer span.Finish()

	return ap.ag.GetBoolWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetInt(ctx context.Context, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt")
	defer span.Finish()

	return ap.ag.GetInt(key)
}

func (ap *apolloConfigCenter) GetIntWithNamespace(ctx context.Context, namespace, key string) (int, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetIntWithNamespace")
	defer span.Finish()

	return ap.ag.GetIntWithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetFloat64(ctx context.Context, key string) (float64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetFloat64")
	defer span.Finish()
	return ap.ag.GetFloat64(key)
}

func (ap *apolloConfigCenter) GetFloat64WithNamespace(ctx context.Context, namespace, key string) (float64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetFloat64WithNamespace")
	defer span.Finish()

	return ap.ag.GetFloat64WithNamespace(namespace, key)
}

func (ap *apolloConfigCenter) GetInt64(ctx context.Context, key string) (int64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt64")
	defer span.Finish()
	val, ok := ap.ag.GetInt(key)
	return int64(val), ok
}

func (ap *apolloConfigCenter) GetInt64WithNamespace(ctx context.Context, namespace, key string) (int64, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt64WithNamespace")
	defer span.Finish()
	val, ok := ap.ag.GetIntWithNamespace(namespace, key)
	return int64(val), ok
}

func (ap *apolloConfigCenter) GetInt32(ctx context.Context, key string) (int32, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt32")
	defer span.Finish()
	val, ok := ap.ag.GetInt(key)
	return int32(val), ok
}

func (ap *apolloConfigCenter) GetInt32WithNamespace(ctx context.Context, namespace, key string) (int32, bool) {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt32WithNamespace")
	defer span.Finish()
	val, ok := ap.ag.GetIntWithNamespace(namespace, key)
	return int32(val), ok
}

func (ap *apolloConfigCenter) GetIntSlice(ctx context.Context, keyPrefix string) ([]int, bool) {
	fun := "apolloConfigCenter.GetIntSlice -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetIntSlice")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getIntSlice(ctx, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetIntSliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int, bool) {
	fun := "apolloConfigCenter.GetIntSliceWithNamespace -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetIntSliceWithNamespace")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getIntSliceWithNamespace(ctx, namespace, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetInt64Slice(ctx context.Context, keyPrefix string) ([]int64, bool) {
	fun := "apolloConfigCenter.GetInt64Slice -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt64Slice")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getInt64Slice(ctx, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetInt64SliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int64, bool) {
	fun := "apolloConfigCenter.GetInt64SliceWithNamespace -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt64SliceWithNamespace")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getInt64SliceWithNamespace(ctx, namespace, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetInt32Slice(ctx context.Context, keyPrefix string) ([]int32, bool) {
	fun := "apolloConfigCenter.GetInt32Slice -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt32Slice")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getInt32Slice(ctx, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetInt32SliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]int32, bool) {
	fun := "apolloConfigCenter.GetInt32SliceWithNamespace -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetInt32SliceWithNamespace")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getInt32SliceWithNamespace(ctx, namespace, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetStringSlice(ctx context.Context, keyPrefix string) ([]string, bool) {
	fun := "apolloConfigCenter.GetStringSlice -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetStringSlice")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getStringSlice(ctx, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) GetStringSliceWithNamespace(ctx context.Context, namespace, keyPrefix string) ([]string, bool) {
	fun := "apolloConfigCenter.GetStringSliceWithNamespace -->"
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetStringSliceWithNamespace")
	defer span.Finish()

	ks := ap.GetAllKeys(ctx)
	targetMap, err := ap.getSliceIdxMap(ctx, ks, keyPrefix)
	if err != nil {
		xlog.Errorf(ctx, "%s agollo atoi err:%v", fun, err)
		return nil, false
	}
	return ap.getStringSliceWithNamespace(ctx, namespace, keyPrefix, targetMap)
}

func (ap *apolloConfigCenter) getSliceIdxMap(ctx context.Context, ks []string, keyPrefix string) (map[int]struct{}, error) {
	targetMap := make(map[int]struct{})
	for _, k := range ks {
		if strings.HasPrefix(k, keyPrefix) {
			karray := strings.Split(k, "[")
			if len(karray) != 2 {
				continue
			}
			if karray[0] == keyPrefix {
				kidxStr := strings.Split(karray[1], "]")
				kidx, err := strconv.Atoi(kidxStr[0])
				if err != nil {
					return targetMap, err
				}
				targetMap[kidx] = struct{}{}
			}
		}
	}
	return targetMap, nil
}

func (ap *apolloConfigCenter) GetAllKeys(ctx context.Context) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetAllKeys")
	defer span.Finish()

	return ap.ag.GetAllKeys("application")
}

func (ap *apolloConfigCenter) GetAllKeysWithNamespace(ctx context.Context, namespace string) []string {
	span, _ := opentracing.StartSpanFromContext(ctx, "apolloConfigCenter.GetAllKeysWithNamespace")
	defer span.Finish()

	return ap.ag.GetAllKeys(namespace)
}

func (ap *apolloConfigCenter) Unmarshal(ctx context.Context, v interface{}) error {
	return ap.UnmarshalWithNamespace(ctx, defaultNamespaceApplication, v)
}

func (ap *apolloConfigCenter) UnmarshalWithNamespace(ctx context.Context, namespace string, v interface{}) error {
	var kv = map[string]string{}

	ks := ap.GetAllKeysWithNamespace(ctx, namespace)
	for _, k := range ks {
		if v, ok := ap.GetStringWithNamespace(ctx, namespace, k); ok {
			kv[k] = v
		}
	}

	return mmarshal.UnmarshalKV(kv, v)
}

func (ap *apolloConfigCenter) UnmarshalKey(ctx context.Context, key string, v interface{}) error {
	return ap.UnmarshalKeyWithNamespace(ctx, defaultNamespaceApplication, key, v)
}

func (ap *apolloConfigCenter) UnmarshalKeyWithNamespace(ctx context.Context, namespace string, key string, v interface{}) error {
	var kv = map[string]string{}

	ks := ap.GetAllKeysWithNamespace(ctx, namespace)
	for _, k := range ks {
		if v, ok := ap.GetStringWithNamespace(ctx, namespace, k); ok {
			kv[k] = v
		}
	}

	bs, err := mmarshal.Marshal(&kv)
	if err != nil {
		return err
	}

	return mmarshal.UnmarshalKey(key, bs, v)
}

func (ap *apolloConfigCenter) getIntSlice(ctx context.Context, keyPrefix string, kMap map[int]struct{}) ([]int, bool) {
	fun := "apolloConfigCenter.getIntSlice -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = val
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getIntSliceWithNamespace(ctx context.Context, namespace, keyPrefix string, kMap map[int]struct{}) ([]int, bool) {
	fun := "apolloConfigCenter.getIntSliceWithNamespace -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = val
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getInt64Slice(ctx context.Context, keyPrefix string, kMap map[int]struct{}) ([]int64, bool) {
	fun := "apolloConfigCenter.getInt64Slice -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int64, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = int64(val)
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getInt64SliceWithNamespace(ctx context.Context, namespace, keyPrefix string, kMap map[int]struct{}) ([]int64, bool) {
	fun := "apolloConfigCenter.getInt64SliceWithNamespace -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int64, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = int64(val)
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getInt32Slice(ctx context.Context, keyPrefix string, kMap map[int]struct{}) ([]int32, bool) {
	fun := "apolloConfigCenter.getInt32Slice -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int32, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = int32(val)
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getInt32SliceWithNamespace(ctx context.Context, namespace, keyPrefix string, kMap map[int]struct{}) ([]int32, bool) {
	fun := "apolloConfigCenter.getInt32SliceWithNamespace -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]int32, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetInt(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = int32(val)
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getStringSlice(ctx context.Context, keyPrefix string, kMap map[int]struct{}) ([]string, bool) {
	fun := "apolloConfigCenter.getStringSlice -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]string, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetString(keyPrefix + "[" + strconv.Itoa(k) + "]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = val
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) getStringSliceWithNamespace(ctx context.Context, namespace, keyPrefix string, kMap map[int]struct{}) ([]string, bool) {
	fun := "apolloConfigCenter.getStringSliceWithNamespace -->"

	arrayLen := len(kMap)
	if arrayLen == 0 {
		return nil, false
	}

	kSlice := make([]string, arrayLen)

	for k, _ := range kMap {
		val, ok := ap.ag.GetStringWithNamespace(namespace, keyPrefix+"["+strconv.Itoa(k)+"]")
		if !ok {
			xlog.Errorf(ctx, "%s get error key: %s", fun, keyPrefix+"["+strconv.Itoa(k)+"]")
			return nil, false
		}
		kSlice[k] = val
	}

	return kSlice, true
}

func (ap *apolloConfigCenter) SetCluster(cluster string) {
	ap.conf.Cluster = cluster
}
func (ap *apolloConfigCenter) SetCacheDir(cacheDir string) {
	ap.conf.CacheDir = cacheDir
}
func (ap *apolloConfigCenter) SetIPHost(ipHost string) {
	ap.conf.IP = ipHost
}

func getEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func confFromEnv() *apogo.Conf {
	cluster := getEnvWithDefault(envApolloCluster, defaultCluster)
	hostport := getEnvWithDefault(envApolloHostPort, defaultHostPort)

	return &apogo.Conf{
		Cluster:  cluster,
		CacheDir: defaultCacheDir,
		IP:       hostport,
	}
}

// NOTE: apollo 不支持在项目名称中使用 '/'，因此规定用 '.' 代替 '/'
//       base/authapi => base.authapi
func normalizeServiceName(serviceName string) string {
	return strings.Replace(serviceName, "/", ".", -1)
}
