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
 @Time    : 2024/10/12 -- 11:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: label.go
*/

package xprometheus

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
)

// 常用labels
const (
	// 组名
	LabelGroupName = "group_name"
	// 服务名
	LabelServiceName = "servname"
	// 服务ID
	LabelServiceID = "servid"
	// 实例名称
	LabelInstance = "instance"
	// 来源
	LabelSource = "source"
	LabelAPI    = "api"
	LabelType   = "type"
	// 自定义状态码
	LabelErrCode = "errcode"
	// apm
	// 调用方服务名
	LabelCallerService = "caller_service"
	// 被调方服务名
	LabelCalleeService = "callee_service"
	// 调用方接入点名
	LabelCallerEndpoint = "caller_endpoint"
	// 被调方接入点名
	LabelCalleeEndpoint = "callee_endpoint"
	// 调用方服务 id
	LabelCallerServiceID = "caller_service_id"
	// 调用结果状态
	LabelCallStatus = "call_status"
	// 消息总线类型
	LabelMessageBusType = "message_bus_type"
)

var forbiddenChars = regexp.MustCompile("[ .=\\-/]")

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}

func SafePrometheusValue(v string) string {
	return forbiddenChars.ReplaceAllString(v, "_")
}
