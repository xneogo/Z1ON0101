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
 @Time    : 2025/7/1 -- 16:52
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xscanner xsql/xscanner/x_impl.go
*/

package xscanner

import (
	"encoding/json"
	"fmt"
	"github.com/qiguanzhu/infra/nerv/magi/xreflect"
	"github.com/qiguanzhu/infra/pkg"
	"github.com/qiguanzhu/infra/pkg/consts"
	"github.com/qiguanzhu/infra/seele/zsql"
	"github.com/qiguanzhu/infra/seele/zsql/sqlutils"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"
)

var userDefinedTagName string

// SetTagName can be set only once
func SetTagName(name string) {
	if userDefinedTagName != "" {
		return
	}
	userDefinedTagName = name
}

type XScanner struct{}

func (s XScanner) Scan(rows zsql.Rows, target interface{}, _ zsql.BindFunc) error {
	if nil == target || reflect.ValueOf(target).IsNil() || reflect.TypeOf(target).Kind() != reflect.Ptr {
		return pkg.ErrScannerTargetNotSettable
	}

	data, err := sqlutils.ResolveDataFromRows(rows)
	if nil != err {
		return err
	}

	switch reflect.TypeOf(target).Elem().Kind() {
	case reflect.Slice:
		if nil == data {
			return nil
		}
		err = bindSlice(data, target)
	default:
		if nil == data {
			return pkg.ErrScannerEmptyResult
		}
		err = bind(data[0], target)
	}

	return err
}

func (s XScanner) ScanMap(rows zsql.Rows) ([]map[string]interface{}, error) {
	return sqlutils.ResolveDataFromRows(rows)
}

func (s XScanner) ScanMapDecode(rows zsql.Rows) ([]map[string]interface{}, error) {
	results, err := sqlutils.ResolveDataFromRows(rows)
	if nil != err {
		return nil, err
	}
	for i := 0; i < len(results); i++ {
		for k, v := range results[i] {
			rv, ok := v.([]uint8)
			if !ok {
				continue
			}
			s := string(rv)
			// convert to int
			intVal, err := strconv.Atoi(s)
			if err == nil {
				results[i][k] = intVal
				continue
			}
			// convert to float64
			floatVal, err := strconv.ParseFloat(s, 64)
			if err == nil {
				results[i][k] = floatVal
				continue
			}
			// convert to string
			results[i][k] = s
		}
	}
	return results, nil
}

func (s XScanner) ScanMapDecodeClose(rows zsql.Rows) ([]map[string]interface{}, error) {
	result, err := s.ScanMapDecode(rows)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return result, err
}

// ScanMapClose is the same as ScanMap and close the rows
func (s XScanner) ScanMapClose(rows zsql.Rows) ([]map[string]interface{}, error) {
	result, err := s.ScanMap(rows)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return result, err
}

// ScanClose is the same as Scan and helps you Close the rows
// Don't exec the rows.Close after calling this
func (s XScanner) ScanClose(rows zsql.Rows, target interface{}, f zsql.BindFunc) error {
	err := s.Scan(rows, target, f)
	if nil != rows {
		errClose := rows.Close()
		if err == nil {
			err = pkg.NewCloseErr(errClose)
		}
	}
	return err
}

// caller must guarantee to pass a &slice as the second param
func bindSlice(arr []map[string]interface{}, target interface{}) error {
	targetObj := reflect.ValueOf(target)
	if !targetObj.Elem().CanSet() {
		return pkg.ErrScannerTargetNotSettable
	}
	length := len(arr)
	valueArrObj := reflect.MakeSlice(targetObj.Elem().Type(), 0, length)
	typeObj := valueArrObj.Type().Elem()
	var err error
	for i := 0; i < length; i++ {
		newObj := reflect.New(typeObj)
		newObjInterface := newObj.Interface()
		err = bind(arr[i], newObjInterface)
		if nil != err {
			return err
		}
		valueArrObj = reflect.Append(valueArrObj, newObj.Elem())
	}
	targetObj.Elem().Set(valueArrObj)
	return nil
}

func bind(result map[string]interface{}, target interface{}) (resp error) {
	if nil != resp {
		return
	}
	defer func() {
		if r := recover(); nil != r {
			resp = fmt.Errorf("error:[%v], stack:[%s]", r, string(debug.Stack()))
		}
	}()
	valueObj := reflect.ValueOf(target).Elem()
	if !valueObj.CanSet() {
		return pkg.ErrScannerTargetNotSettable
	}
	typeObj := valueObj.Type()
	if typeObj.Kind() == reflect.Ptr {
		ptrType := typeObj.Elem()
		newObj := reflect.New(ptrType)
		newObjInterface := newObj.Interface()
		err := bind(result, newObjInterface)
		if nil == err {
			valueObj.Set(newObj)
		}
		return err
	}
	typeObjName := typeObj.Name()

	for i := 0; i < valueObj.NumField(); i++ {
		fieldTypeI := typeObj.Field(i)
		fieldName := fieldTypeI.Name

		// for convenience
		wrapErr := func(from, to reflect.Type) pkg.ScanErr {
			return pkg.NewScanErr(typeObjName, fieldName, from, to)
		}

		valuei := valueObj.Field(i)
		if !valuei.CanSet() {
			continue
		}
		tagName, ok := lookUpTagName(fieldTypeI)
		if !ok || "" == tagName {
			continue
		}
		mapValue, ok := result[tagName]
		if !ok || mapValue == nil {
			continue
		}
		// if one field is a pointer type, we must allocate memory for it first and json unmarshal
		// except for that the pointer type implements the interface ByteUnmarshaler
		if fieldTypeI.Type.Kind() == reflect.Ptr && !fieldTypeI.Type.Implements(_byteUnmarshalerType) {
			if fieldTypeI.Type.Elem().Kind() == reflect.Struct {
				err := defaultStructUnmarshal(&valuei, mapValue)
				if err == nil {
					continue
				}
			}
			// 老逻辑
			valuei.Set(reflect.New(fieldTypeI.Type.Elem()))
			valuei = valuei.Elem()
		}
		if fieldTypeI.Type.Kind() == reflect.Slice {
			err := defaultSliceUnmarshal(&valuei, mapValue)
			// 如果解析成功继续往下走，否则按照之前的逻辑解析
			if err == nil {
				continue
			}
		}
		// 结构体类型走unmarshal逻辑
		if fieldTypeI.Type.Kind() == reflect.Struct && fieldTypeI.Type.String() != "time.Time" {
			vPtr := reflect.New(valuei.Type())
			err := defaultStructUnmarshal(&vPtr, mapValue)
			if err == nil {
				valuei.Set(vPtr.Elem())
				continue
			}
		}
		err := convert(mapValue, valuei, wrapErr)
		if nil != err {
			return err
		}
	}
	return nil
}

var _byteUnmarshalerType = reflect.TypeOf(new(zsql.ByteUnmarshaler)).Elem()

type convertErrWrapper func(from, to reflect.Type) pkg.ScanErr

func isIntSeriesType(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

func isUintSeriesType(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uint64
}

func isFloatSeriesType(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}

func lookUpTagName(typeObj reflect.StructField) (string, bool) {
	var tName string
	if "" != userDefinedTagName {
		tName = userDefinedTagName
	} else {
		tName = consts.DefaultTagName
	}
	name, ok := typeObj.Tag.Lookup(tName)
	if !ok {
		return "", false
	}
	name = xreflect.ResolveTagName(name)
	return name, ok
}

func convert(mapValue interface{}, valuei reflect.Value, wrapErr convertErrWrapper) error {
	// vit: ValueI Type
	vit := valuei.Type()
	// mvt: MapValue Type
	mvt := reflect.TypeOf(mapValue)
	if nil == mvt {
		return nil
	}
	// []byte tp []byte && time.Time to time.Time
	if mvt.AssignableTo(vit) {
		valuei.Set(reflect.ValueOf(mapValue))
		return nil
	}
	// time.Time to string
	switch assertT := mapValue.(type) {
	case time.Time:
		return handleConvertTime(assertT, mvt, vit, &valuei, wrapErr)
	}

	// according to go-mysql-driver/mysql, driver.Value type can only be:
	// int64 or []byte(> maxInt64)
	// float32/float64
	// []byte
	// time.Time if parseTime=true or DATE type will be converted into []byte
	switch mvt.Kind() {
	case reflect.Int64:
		if isIntSeriesType(vit.Kind()) {
			valuei.SetInt(mapValue.(int64))
		} else if isUintSeriesType(vit.Kind()) {
			valuei.SetUint(uint64(mapValue.(int64)))
		} else if vit.Kind() == reflect.Bool {
			v := mapValue.(int64)
			if v > 0 {
				valuei.SetBool(true)
			} else {
				valuei.SetBool(false)
			}
		} else if vit.Kind() == reflect.String {
			valuei.SetString(strconv.FormatInt(mapValue.(int64), 10))
		} else {
			return wrapErr(mvt, vit)
		}
	case reflect.Float32:
		if isFloatSeriesType(vit.Kind()) {
			valuei.SetFloat(float64(mapValue.(float32)))
		} else {
			return wrapErr(mvt, vit)
		}
	case reflect.Float64:
		if isFloatSeriesType(vit.Kind()) {
			valuei.SetFloat(mapValue.(float64))
		} else {
			return wrapErr(mvt, vit)
		}
	case reflect.Slice:
		return handleConvertSlice(mapValue, mvt, vit, &valuei, wrapErr)
	default:
		return wrapErr(mvt, vit)
	}
	return nil
}

func handleConvertSlice(mapValue interface{}, mvt, vit reflect.Type, valuei *reflect.Value, wrapErr convertErrWrapper) error {
	mapValueSlice, ok := mapValue.([]byte)
	if !ok {
		return pkg.ErrScannerSliceToString
	}
	mapValueStr := string(mapValueSlice)
	vitKind := vit.Kind()
	switch {
	case vitKind == reflect.String:
		valuei.SetString(mapValueStr)
	case isIntSeriesType(vitKind):
		intVal, err := strconv.ParseInt(mapValueStr, 10, 64)
		if nil != err {
			return wrapErr(mvt, vit)
		}
		valuei.SetInt(intVal)
	case isUintSeriesType(vitKind):
		uintVal, err := strconv.ParseUint(mapValueStr, 10, 64)
		if nil != err {
			return wrapErr(mvt, vit)
		}
		valuei.SetUint(uintVal)
	case isFloatSeriesType(vitKind):
		floatVal, err := strconv.ParseFloat(mapValueStr, 64)
		if nil != err {
			return wrapErr(mvt, vit)
		}
		valuei.SetFloat(floatVal)
	case vitKind == reflect.Bool:
		intVal, err := strconv.ParseInt(mapValueStr, 10, 64)
		if nil != err {
			return wrapErr(mvt, vit)
		}
		if intVal > 0 {
			valuei.SetBool(true)
		} else {
			valuei.SetBool(false)
		}
	default:
		if _, ok := valuei.Interface().(zsql.ByteUnmarshaler); ok {
			return byteUnmarshal(mapValueSlice, valuei, wrapErr)
		}
		return wrapErr(mvt, vit)
	}
	return nil
}

// valuei Here is the type of ByteUnmarshaler
func byteUnmarshal(mapValueSlice []byte, valuei *reflect.Value, wrapErr convertErrWrapper) error {
	var pt reflect.Value
	initFlag := false
	// init pointer
	if valuei.IsNil() {
		pt = reflect.New(valuei.Type().Elem())
		initFlag = true
	} else {
		pt = *valuei
	}
	err := pt.Interface().(zsql.ByteUnmarshaler).UnmarshalByte(mapValueSlice)
	if nil != err {
		structName := pt.Elem().Type().Name()
		return fmt.Errorf("[scanner]: %s.UnmarshalByte fail to unmarshal the bytes, err: %s", structName, err)
	}
	if initFlag {
		valuei.Set(pt)
	}
	return nil
}

func handleConvertTime(assertT time.Time, mvt, vit reflect.Type, valuei *reflect.Value, wrapErr convertErrWrapper) error {
	if vit.Kind() == reflect.String {
		sTime := assertT.Format(consts.CTimeFormat)
		valuei.SetString(sTime)
		return nil
	}
	return wrapErr(mvt, vit)
}

func defaultStructUnmarshal(valuei *reflect.Value, mapValue interface{}) error {
	var pt reflect.Value
	initFlag := false
	// init pointer
	if valuei.IsNil() {
		pt = reflect.New(valuei.Type().Elem())
		initFlag = true
	} else {
		pt = *valuei
	}
	err := json.Unmarshal(mapValue.([]byte), pt.Interface())
	if nil != err {
		structName := pt.Elem().Type().Name()
		return fmt.Errorf("[scanner]: %s.Unmarshal fail to unmarshal the bytes, err: %s", structName, err)
	}
	if initFlag {
		valuei.Set(pt)
	}
	return nil
}

func defaultSliceUnmarshal(valuei *reflect.Value, mapValue interface{}) error {
	var pt reflect.Value
	initFlag := false
	// init pointer
	if valuei.IsNil() {
		// 创建slice
		itemslice := reflect.MakeSlice(valuei.Type(), 0, 0)
		// 指针赋值
		pt = reflect.New(itemslice.Type())
		// 指针指向slice
		pt.Elem().Set(itemslice)
		initFlag = true
	} else {
		pt = *valuei
	}
	err := json.Unmarshal(mapValue.([]byte), pt.Interface())
	if nil != err {
		structName := pt.Elem().Type().Name()
		return fmt.Errorf("[scanner]: %s.Unmarshal fail to unmarshal the bytes, err: %s", structName, err)
	}
	if initFlag {
		valuei.Set(pt.Elem())
	}
	return nil
}
