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
 @Time    : 2024/10/25 -- 18:38
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: center.go
*/

package xconfig

import (
	"github.com/xneogo/matrix/mconfig"
)

type ConfigCenterInfo struct {
	// 配置类型
	CfgType mconfig.ConfigureType
	ConfigDriverInfo
}

type ConfigDriverInfo struct {
	// 配置信息的根路径，
	CfgRootPath string
	// 当前初始化配置中心的应用
	AppKey string
	// namespace
	NamespaceNames []string
}
