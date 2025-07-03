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
 @Description: const.go
*/

package xapollo

import (
	"github.com/xneogo/matrix/mconfig"
)

const (
	// DefaultApolloMiddlewareService ...
	DefaultApolloMiddlewareService = "middleware"
	// DefaultApolloMQNamespace ...
	DefaultApolloMQNamespace = "infra.mq"
	// DefaultApolloCacheNamespace ...
	DefaultApolloCacheNamespace = "infra.cache"
	// DefaultApolloMysqlNamespace ...
	DefaultApolloMysqlNamespace = "infra.mysql"
	// DefaultApolloTraceNamespace ...
	DefaultApolloTraceNamespace = "infra.trace"
	// ConfigTypeApollo ...
	ConfigTypeApollo mconfig.ConfigureType = "apollo"
)
