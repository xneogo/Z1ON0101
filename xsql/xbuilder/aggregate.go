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
 @Time    : 2024/10/9 -- 16:19
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: aggregate.go
*/

package xbuilder

import (
	"reflect"
	"strconv"

	"github.com/xneogo/Z1ON0101/xsql/factory"
)

type resultResolve struct {
	data interface{}
}

func (r resultResolve) Int64() int64 {
	switch t := r.data.(type) {
	case int64:
		return t
	case int32:
		return int64(t)
	case int:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case []uint8:
		i64, err := strconv.ParseInt(string(t), 10, 64)
		if nil != err {
			return int64(r.Float64())
		}
		return i64
	default:
		return 0
	}
}

func (r resultResolve) Float64() float64 {
	switch t := r.data.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case []uint8:
		f64, _ := strconv.ParseFloat(string(t), 64)
		return f64
	default:
		return float64(r.Int64())
	}
}

type agBuilder string

func (a agBuilder) Symbol() string {
	return string(a)
}

// AggregateCount count(col)
func AggregateCount(col string) factory.AggregateSymbolBuilder {
	return agBuilder("count(" + col + ")")
}

// AggregateSum sum(col)
func AggregateSum(col string) factory.AggregateSymbolBuilder {
	return agBuilder("sum(" + col + ")")
}

// AggregateAvg avg(col)
func AggregateAvg(col string) factory.AggregateSymbolBuilder {
	return agBuilder("avg(" + col + ")")
}

// AggregateMax max(col)
func AggregateMax(col string) factory.AggregateSymbolBuilder {
	return agBuilder("max(" + col + ")")
}

// AggregateMin min(col)
func AggregateMin(col string) factory.AggregateSymbolBuilder {
	return agBuilder("min(" + col + ")")
}

// OmitEmpty is a helper function to clear where map zero value
func OmitEmpty(where map[string]interface{}, omitKey []string) map[string]interface{} {
	for _, key := range omitKey {
		v, ok := where[key]
		if !ok {
			continue
		}

		if isZero(reflect.ValueOf(v)) {
			delete(where, key)
		}
	}
	return where
}

// OmitAllEmpty clear all where/data map zero value
func OmitAllEmpty(data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		if isZero(reflect.ValueOf(v)) {
			delete(data, k)
		}
	}
	return data
}

// isZero reports whether a value is a zero value
// Including support: Bool, Array, String, Float32, Float64, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr
// Map, Slice, Interface, Struct
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Array, reflect.String:
		return v.Len() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Interface:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	// Traverse the Struct and only return true
	// if all of its fields return IsZero == true
	n := v.NumField()
	for i := 0; i < n; i++ {
		vf := v.Field(i)
		if !isZero(vf) {
			return false
		}
	}
	return true
}
