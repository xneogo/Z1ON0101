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
 @Time    : 2024/11/10 -- 17:15
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: remote.go
*/

package xviper

import "strings"

// ViperRemoteProvider copied from viper
type ViperRemoteProvider struct {
	ProviderType string
	Addrs        []string
	Location     string
	SecretKey    string
}

func (rp ViperRemoteProvider) Provider() string {
	return rp.ProviderType
}

func (rp ViperRemoteProvider) Endpoint() string {
	return strings.Join(rp.Addrs, ";")
}

func (rp ViperRemoteProvider) Path() string {
	return rp.Location
}

func (rp ViperRemoteProvider) SecretKeyring() string {
	return rp.SecretKey
}
