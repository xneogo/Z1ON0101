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
 @Time    : 2025/4/15 -- 12:05
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xmq xmq/global.go
*/

package xmq

import "context"

type Handler interface {
	CommitMsg(ctx context.Context) error
}

type AckHandler interface {
	Ack(ctx context.Context) error
}

type QueueModel interface {
	// WriteMsg 写入消息
	WriteMsg(ctx context.Context, topic string, key string, val interface{}) error

	// ReadMsgByGroup 读完消息后会自动提交offset
	ReadMsgByGroup(ctx context.Context, topic, groupId string, value interface{}) (context.Context, error)
	// ReadMsgByPartition ...
	ReadMsgByPartition(ctx context.Context, topic string, partition int, value interface{}) (context.Context, error)
	// FetchMsgByGroup 读完消息后不会自动提交offset,需要手动调用Handle.CommitMsg方法来提交offset
	FetchMsgByGroup(ctx context.Context, topic, groupId string, value interface{}) (context.Context, Handler, error)
}

type KafkaQueueImpl struct{}

type DelayQueueImpl struct{}

type PulsarQueueImpl struct{}
