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
 @Time    : 2024/10/12 -- 15:41
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: driver.go
*/

package xconfig

import (
	"context"
	"fmt"
	"github.com/xneogo/matrix/mconfig"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[mconfig.ConfigureType]mconfig.Driver)
)

// Register ...
func Register(cType mconfig.ConfigureType, driver mconfig.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("xconfig: driver is nil")
	}

	if _, dup := drivers[cType]; dup {
		panic("xconfig: driver can called Register only once")
	}

	drivers[cType] = driver
}

// Drivers returns a sorted list of the names of the registered driver.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()

	var list []string
	for name := range drivers {
		list = append(list, string(name))
	}

	sort.Strings(list)
	return list
}

// GetDriver returns a driver implement by config type.
func GetDriver(cType mconfig.ConfigureType) (mconfig.Driver, error) {
	driversMu.RLock()
	driveri, ok := drivers[cType]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unkown config type:%s", string(cType))
	}
	return driveri, nil
}

// NewConfigCenter ...
func NewConfigCenter(ctx context.Context, cType mconfig.ConfigureType, serviceName string, namespaceNames []string, options ...mconfig.Option) (mconfig.ConfigCenter, error) {
	driver, err := GetDriver(cType)
	if err != nil {
		return nil, fmt.Errorf("new config center err:%s", err.Error())
	}
	return driver.New(ctx, serviceName, namespaceNames, options...)
}
