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
 @Time    : 2024/11/5 -- 18:30
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: pipeline.go
*/

package xredis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// RedisPipelineClient by RedisClient get pipeline
type RedisPipelineClient struct {
	*RedisCmd
	client redis.Pipeliner
}

func NewRedisPipelineClient(ctx context.Context, redisConf *RedisConfig) (*RedisPipelineClient, error) {
	redisClient, err := NewRedisClient(ctx, redisConf)
	if err != nil || redisClient == nil {
		return nil, err
	}
	return redisClient.Pipeline(ctx)
}

func (m *RedisPipelineClient) Exec(ctx context.Context) (cmds []redis.Cmder, err error) {
	command := "RedisPipelineClient.Exec"
	_ = m.WrapperLog(ctx, command, func() error {
		cmds, err = m.client.Exec(ctx)
		return err
	})
	return
}

func (m *RedisPipelineClient) Discard(ctx context.Context) (err error) {
	command := "RedisPipelineClient.Discard"
	_ = m.WrapperLog(ctx, command, func() error {
		m.client.Discard()
		return nil
	})
	return
}

func (m *RedisPipelineClient) GetRawClient(ctx context.Context) redis.Pipeliner {
	return m.client
}
