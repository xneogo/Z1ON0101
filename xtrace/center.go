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
 @Time    : 2024/10/18 -- 16:22
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: center.go
*/

package xtrace

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/xneogo/Z1ON0101/xconfig"
	"github.com/xneogo/Z1ON0101/xconfig/xapollo"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/mconfig"
)

var initApolloLock = &sync.Mutex{}

var apolloCenter mconfig.ConfigCenter

var apolloSpanFilterConfig *spanFilterConfig

func initApolloCenter(ctx context.Context) error {
	if apolloCenter != nil {
		return nil
	}

	initApolloLock.Lock()
	defer initApolloLock.Unlock()

	if apolloCenter != nil {
		return nil
	}

	xlog.Infof(ctx, "initApolloCenter --> xtrace apollo center not found, init")

	return initGlobalApolloCenterWithoutLock(ctx)
}

func initGlobalApolloCenterWithoutLock(ctx context.Context) error {
	namespaceList := []string{xapollo.DefaultApolloTraceNamespace}
	newCenter, err := xconfig.NewConfigCenter(ctx, xapollo.ConfigTypeApollo, xapollo.DefaultApolloMiddlewareService, namespaceList)
	if err != nil {
		return fmt.Errorf("init apollo with service %s namespace %s error, %s",
			xapollo.DefaultApolloMiddlewareService, strings.Join(namespaceList, " "), err.Error())
	}

	apolloCenter = newCenter
	xlog.Infof(ctx, "initApolloCenter --> initGlobalApolloCenter success")
	return nil
}
