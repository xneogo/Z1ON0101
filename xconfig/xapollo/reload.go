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
 @Time    : 2024/11/1 -- 09:53
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: reload.go
*/

package xapollo

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"

	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/mconfig/mobserver"
	"github.com/xneogo/matrix/mentity"
	"github.com/xneogo/matrix/mserv"
	"github.com/xneogo/saferun"
)

// ResetApolloConfig
// 全量重新加载阿波罗配置
// type ApolloConfig struct {
//	Apollotest       bool `apollokey:"apollotest"`
//	ApollotestStruct struct {
//		Age  int64  `json:"age"`
//		Name string `json:"name"`
//	} `apollokey:"apolloteststruct"`
// }
// ctx为上下文信息,sb为框架基础base结构,cfgPrtPtr为配置结构体指针的指针,例: conf := &ApolloConfig 应传入&conf,因为需要指针原子替换
// 注意::配置结构体不要嵌套结构体指针
func ResetApolloConfig(ctx context.Context, sb mserv.ServerSessionProxy[mentity.ServInfo], cfgPtrPtr interface{}) error {
	return resetApolloConfigByTag(ctx, sb, cfgPtrPtr, "json")
}

// ResetApolloConfigByTag 根据自定义的tag加载阿波罗配置
func ResetApolloConfigByTag(ctx context.Context, sb mserv.ServerSessionProxy[mentity.ServInfo], cfgPtrPtr interface{}, tag string) error {
	return resetApolloConfigByTag(ctx, sb, cfgPtrPtr, tag)
}

// LoadApolloConfigByEvent 根据变更的event,增量更新配置
func LoadApolloConfigByEvent(ctx context.Context, sb mserv.ServerSessionProxy[mentity.ServInfo], cfgPtrPtr interface{}, event *mobserver.ChangeEvent) error {
	return resetApolloConfigByTag(ctx, sb, cfgPtrPtr, "json")
}
func LoadApolloConfigByEventTag(ctx context.Context, sb mserv.ServerSessionProxy[mentity.ServInfo], cfgPtrPtr interface{}, tag string, event *mobserver.ChangeEvent) error {
	return resetApolloConfigByTag(ctx, sb, cfgPtrPtr, tag)
}

// TODO::当前只支持简单的类型,且cfg应为结构体嵌套,不能为指针。
func resetApolloConfigByTag(ctx context.Context, sb mserv.ServerSessionProxy[mentity.ServInfo], cfgPtrPtr interface{}, tag string) (reserr error) {
	fun := "util.LoadApolloConfig -->"
	defer func() {
		if panicErr := recover(); panicErr != nil {
			reserr = saferun.DumpStack(panicErr)
			return
		}
	}()
	if cfgPtrPtr == nil {
		return fmt.Errorf("config is nil")
	}
	// 传进来的是cfg结构体的指针的指针,因为需要原子操作替换指针且需要反射设置值 cfgPtrPtr为 **Config
	cfgPtrPtrValue := reflect.ValueOf(cfgPtrPtr) // 拿到反射值
	if cfgPtrPtrValue.Kind() != reflect.Ptr {
		return fmt.Errorf("need config ptr ptr~")
	}
	if cfgPtrPtrValue.Elem().Kind() != reflect.Ptr {
		return fmt.Errorf("need config ptr ptr~")
	}

	// 配置结构体的反射信息
	cfgValue := cfgPtrPtrValue.Elem().Elem()
	cfgType := cfgValue.Type()

	// 这个oldCfgPtr用来最后原子赋值
	oldCfgPtr := (*unsafe.Pointer)(unsafe.Pointer(cfgPtrPtrValue.Pointer()))
	// 创建一个新的结构体用来替换
	newCfgPtr := reflect.New(cfgType)

	// 迭代写入空结构体配置数据
	fieldNum := newCfgPtr.Elem().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldValue := newCfgPtr.Elem().Field(i)
		// fixme::要是结构体嵌套结构体指针可能出现问题
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		fieldType := newCfgPtr.Elem().Type().Field(i)
		cfgKey := fieldType.Tag.Get(tag)
		if cfgKey == "" {
			continue
		}
		// 目前只支持字符串 bool int 和struct map
		switch fieldType.Type.Kind() {
		case reflect.String:
			v, ok := sb.ConfigCenter(ctx).GetString(ctx, cfgKey)
			if !ok {
				xlog.Errorf(ctx, "%s get apollo key:%s fail", fun, cfgKey)
				continue
			}
			fieldValue.SetString(v)
		case reflect.Bool:
			v, ok := sb.ConfigCenter(ctx).GetBool(ctx, cfgKey)
			if !ok {
				xlog.Errorf(ctx, "%s get apollo key:%s fail", fun, cfgKey)
				continue
			}
			fieldValue.SetBool(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, ok := sb.ConfigCenter(ctx).GetInt(ctx, cfgKey)
			if !ok {
				xlog.Errorf(ctx, "%s get apollo key:%s fail", fun, cfgKey)
				continue
			}
			fieldValue.SetInt(int64(v))
		case reflect.Struct, reflect.Array, reflect.Map, reflect.Slice:
			v, ok := sb.ConfigCenter(ctx).GetString(ctx, cfgKey)
			if !ok {
				xlog.Errorf(ctx, "%s get apollo key:%s fail", fun, cfgKey)
				continue
			}
			newCal := reflect.New(fieldType.Type)
			cal := newCal.Interface()
			err := json.Unmarshal([]byte(v), &cal)
			if err != nil {
				xlog.Errorf(ctx, "%s get apollo key:%s fail,Unmarshal err:%s", fun, cfgKey, err.Error())
				continue
			}
			fieldValue.Set(newCal.Elem())
		default:
			xlog.Errorf(ctx, "%s get apollo key:%s fail,unsupport type %s", fun, cfgKey, fieldType.Type.String())
			continue
		}
	}
	atomic.StorePointer(
		oldCfgPtr,
		unsafe.Pointer(newCfgPtr.Pointer()),
	)
	return nil
}
