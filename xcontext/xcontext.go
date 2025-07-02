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
 @Time    : 2024/10/10 -- 12:03
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: xcontext.go
*/

package xcontext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/xneogo/matrix/mtransport/gen-go/util/thriftutil"
)

// 由于请求的上下文信息的 thrift 定义在 util 项目中，本模块主要为了避免循环依赖

// ContextCaller store caller info
type ContextCaller struct {
	Method string
}

const (
	// ContextKeyTraceID ...
	ContextKeyTraceID = "traceID"
	// ContextKeyHead ...
	ContextKeyHead = "Head"
	// ContextKeyHeadUID ...
	ContextKeyHeadUID = "uid"
	// ContextKeyHeadSource ...
	ContextKeyHeadSource = "source"
	// ContextKeyHeadIP ...
	ContextKeyHeadIP = "ip"
	// ContextKeyHeadRegion ...
	ContextKeyHeadRegion = "region"
	// ContextKeyHeadDt ...
	ContextKeyHeadDt = "dt"
	// ContextKeyHeadUnionID ...
	ContextKeyHeadUnionID = "unionid"
	// ContextKeyHeadDID
	ContextKeyHeadDID = "h_did"
	// ContextKeyHeadZone
	ContextKeyHeadZone = "zone"
	// ContextKeyHeadZoneName
	ContextKeyHeadZoneName = "zone_name"

	// ContextKeyControl ...
	ContextKeyControl = "Control"

	ContextKeyCaller = "Caller"

	ContextKeyProperties = "properties"

	ContextPropertiesKeyHLC        = "qgz-h_lc"
	ContextPropertiesKeyZone       = "qgz-zone"
	ContextPropertiesKeyZoneName   = "qgz-zone_name"
	ContextPropertiesKeyHiiiHeader = "qgz-hiii-header"
)

// DefaultGroup ...
const DefaultGroup = ""

// ErrInvalidContext ...
var ErrInvalidContext = errors.New("invalid context")

// ContextHeader ...
type ContextHeader interface {
	ToKV() map[string]interface{}
}

type ContextHeaderSetter interface {
	SetKV(key string, value interface{})
}

// ContextControlRouter ...
type ContextControlRouter interface {
	GetControlRouteGroup() (string, bool)
	SetControlRouteGroup(string) error
}

// ContextControlCaller ...
type ContextControlCaller interface {
	GetControlCallerServerName() (string, bool)
	SetControlCallerServerName(string) error
	GetControlCallerServerId() (string, bool)
	SetControlCallerServerId(string) error
	GetControlCallerMethod() (string, bool)
	SetControlCallerMethod(string) error
}

type ContextHeaderCreate func() ContextHeaderSetter

var factory ContextHeaderCreate = func() ContextHeaderSetter {
	return nil
}

func InitContextHeaderFactory(create ContextHeaderCreate) {
	factory = create
}

// GetControlRouteGroup ...
func GetControlRouteGroup(ctx context.Context) (group string, ok bool) {
	value := ctx.Value(ContextKeyControl)
	if isNil(value) {
		ok = false
		return
	}
	control, ok := value.(ContextControlRouter)
	if ok == false {
		return
	}
	return control.GetControlRouteGroup()
}

// SetControlRouteGroup ...
func SetControlRouteGroup(ctx context.Context, group string) (context.Context, error) {
	value := ctx.Value(ContextKeyControl)
	if isNil(value) {
		return ctx, ErrInvalidContext
	}
	control, ok := value.(ContextControlRouter)
	if !ok {
		return ctx, ErrInvalidContext
	}

	err := control.SetControlRouteGroup(group)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, control), nil
}

// GetControlRouteGroupWithDefault ...
func GetControlRouteGroupWithDefault(ctx context.Context, dv string) string {
	if group, ok := GetControlRouteGroup(ctx); ok {
		return group
	}
	return dv
}

// GetControlRouteGroupWithMasterDefault 主干泳道由传入的字符串命名，默认为主干
func GetControlRouteGroupWithMasterName(ctx context.Context, master string) string {
	// 除了主干的其他分支返回分支名
	if group, ok := GetControlRouteGroup(ctx); ok && group != "" {
		return group
	}
	return master
}

func getHeaderByKey(ctx context.Context, key string) (val interface{}, ok bool) {
	var header ContextHeader
	if header, ok = getHeader(ctx); ok {
		val, ok = header.ToKV()[key]
	}
	return
}

// GetUID ...
func GetUID(ctx context.Context) (uid int64, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadUID)
	if ok {
		uid, ok = val.(int64)
	}
	return
}

// GetSource ...
func GetSource(ctx context.Context) (source int32, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadSource)
	if ok {
		source, ok = val.(int32)
	}
	return
}

// GetIP ...
func GetIP(ctx context.Context) (ip string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadIP)
	if ok {
		ip, ok = val.(string)
	}
	return
}

// GetRegion ...
func GetRegion(ctx context.Context) (region string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadRegion)
	if ok {
		region, ok = val.(string)
	}
	return
}

// GetDt ...
func GetDt(ctx context.Context) (dt int32, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadDt)
	if ok {
		dt, ok = val.(int32)
	}
	return
}

// GetUnionID ...
func GetUnionID(ctx context.Context) (unionID string, ok bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadUnionID)
	if ok {
		unionID, ok = val.(string)
	}
	return
}

func GetDID(ctx context.Context) (string, bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadDID)
	if !ok {
		return "", false
	}
	did, ok := val.(string)
	return did, ok
}

func getHeader(ctx context.Context) (header ContextHeader, ok bool) {
	head := ctx.Value(ContextKeyHead)
	if isNil(head) {
		ok = false
		return
	}
	header, ok = head.(ContextHeader)
	return
}

func getHeaderPropertiesByKey(ctx context.Context, key string) (value string, ok bool) {
	properties, ok := getHeaderProperties(ctx)
	if !ok {
		return
	}
	value, ok = properties[key]
	return
}

func getHeaderProperties(ctx context.Context) (map[string]string, bool) {
	if header, ok := getHeader(ctx); ok {
		properties, ok := header.ToKV()[ContextKeyProperties]
		if !ok {
			return nil, ok
		}
		if isNil(properties) {
			return nil, false
		}
		data, ok := properties.(map[string]string)
		return data, ok
	}
	return nil, false
}

func GetPropertiesHLC(ctx context.Context) (string, bool) {
	return getHeaderPropertiesByKey(ctx, ContextPropertiesKeyHLC)
}

func GetPropertiesZone(ctx context.Context) (int32, bool) {
	val, ok := getHeaderPropertiesByKey(ctx, ContextPropertiesKeyZone)
	if !ok {
		return 0, false
	}
	zone, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	return int32(zone), true
}

func GetPropertiesZoneName(ctx context.Context) (string, bool) {
	return getHeaderPropertiesByKey(ctx, ContextPropertiesKeyZoneName)
}

func GetPropertiesHiiiHeader(ctx context.Context) (string, bool) {
	return getHeaderPropertiesByKey(ctx, ContextPropertiesKeyHiiiHeader)
}

func SetPropertiesHiiiHeader(ctx context.Context, info string) context.Context {
	return setHeaderProperties(ctx, ContextPropertiesKeyHiiiHeader, info)
}

func getContextSetter(ctx context.Context) (setter ContextHeaderSetter, ok bool) {
	set := ctx.Value(ContextKeyHead)
	if isNil(set) {
		ok = false
		return
	}
	setter, ok = set.(ContextHeaderSetter)
	return
}

func setHeaderProperties(ctx context.Context, key string, value string) context.Context {
	setter, ok := getContextSetter(ctx)
	if !ok {
		setter = factory()
	}
	if setter == nil {
		return ctx
	}
	properties, ok := getHeaderProperties(ctx)
	if !ok {
		properties = make(map[string]string)
	}
	properties[key] = value
	setter.SetKV(ContextKeyProperties, properties)
	return context.WithValue(ctx, ContextKeyHead, setter)
}

func SetHeaderPropertiesALL(ctx context.Context, hlc string, zone int32, zoneName string) context.Context {
	ctx = SetHeaderPropertiesHLC(ctx, hlc)
	ctx = SetHeaderPropertiesZone(ctx, zone)
	ctx = SetHeaderPropertiesZoneName(ctx, zoneName)
	return ctx
}

func SetHeaderPropertiesHLC(ctx context.Context, hlc string) context.Context {
	return setHeaderProperties(ctx, ContextPropertiesKeyHLC, hlc)
}

func SetHeaderPropertiesZone(ctx context.Context, zone int32) context.Context {
	return setHeaderProperties(ctx, ContextPropertiesKeyZone, fmt.Sprintf("%d", zone))
}

func SetHeaderPropertiesZoneName(ctx context.Context, zoneName string) context.Context {
	return setHeaderProperties(ctx, ContextPropertiesKeyZoneName, zoneName)
}

func GetZone(ctx context.Context) (int32, bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadZone)
	if !ok {
		return 0, false
	}
	zone, ok := val.(int32)
	return zone, ok
}

func GetZoneName(ctx context.Context) (string, bool) {
	val, ok := getHeaderByKey(ctx, ContextKeyHeadZoneName)
	if !ok {
		return "", false
	}
	zoneName, ok := val.(string)
	return zoneName, ok
}

// will get a new caller
func getControlCaller(ctx context.Context) (ContextControlCaller, error) {
	value := ctx.Value(ContextKeyControl)
	if value == nil {
		return nil, ErrInvalidContext
	}
	caller, ok := value.(ContextControlCaller)
	if !ok {
		return nil, ErrInvalidContext
	}
	return createControlCaller(caller)
}

func createControlCaller(caller ContextControlCaller) (ContextControlCaller, error) {
	callerData, err := json.Marshal(caller)
	if err != nil {
		return nil, fmt.Errorf("marshal invalid caller, err: %v", err)
	}
	newCaller := new(thriftutil.Control)
	err = json.Unmarshal(callerData, newCaller)
	if err != nil {
		return nil, fmt.Errorf("unmarshal invalid caller, err: %v", err)
	}

	return newCaller, nil
}

// GetControlCallerServerName ...
func GetControlCallerServerName(ctx context.Context) (serverName string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerServerName()
}

// SetControlCallerServerName ...
func SetControlCallerServerName(ctx context.Context, serverName string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerServerName(serverName)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

// GetControlCallerServerID ...
func GetControlCallerServerID(ctx context.Context) (serverID string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerServerId()
}

// SetControlCallerServerID ...
func SetControlCallerServerID(ctx context.Context, serverID string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerServerId(serverID)
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

// GetControlCallerMethod ...
func GetControlCallerMethod(ctx context.Context) (method string, ok bool) {
	caller, ok := ctx.Value(ContextKeyControl).(ContextControlCaller)
	if !ok {
		return
	}
	return caller.GetControlCallerMethod()
}

// SetControlCallerMethod ...
func SetControlCallerMethod(ctx context.Context, method string) (context.Context, error) {
	caller, err := getControlCaller(ctx)
	if err != nil {
		return ctx, err
	}
	err = caller.SetControlCallerMethod(method)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, ContextKeyControl, caller), nil
}

func getCaller(ctx context.Context) ContextCaller {
	value := ctx.Value(ContextKeyCaller)
	caller, ok := value.(ContextCaller)
	if !ok {
		return ContextCaller{}
	}

	return caller
}

// SetCallerMethod ...
func SetCallerMethod(ctx context.Context, method string) context.Context {
	caller := getCaller(ctx)
	caller.Method = method
	return context.WithValue(ctx, ContextKeyCaller, caller)
}

// GetCallerMethod ...
func GetCallerMethod(ctx context.Context) (method string, ok bool) {
	caller, ok := ctx.Value(ContextKeyCaller).(ContextCaller)
	if !ok {
		return
	}
	return caller.Method, true
}

type ValueContext struct {
	ctx context.Context
}

func (c ValueContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c ValueContext) Done() <-chan struct{}             { return nil }
func (c ValueContext) Err() error                        { return nil }
func (c ValueContext) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// NewValueContext returns a context that is never canceled.
func NewValueContext(ctx context.Context) context.Context {
	return ValueContext{ctx: ctx}
}

// 判断是否为空指针
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr || vi.Kind() == reflect.Map {
		return vi.IsNil()
	}
	return false
}
