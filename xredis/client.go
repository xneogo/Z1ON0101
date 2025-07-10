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
 @Time    : 2025/7/1 -- 18:23
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xredis xredis/client.go
*/

package xredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xneogo/Z1ON0101/xlog"
)

type RedisClient struct {
	*RedisCmd
	client *redis.Client
}

func NewRedisClient(ctx context.Context, redisConf *RedisConfig) (*RedisClient, error) {
	fun := "NewRedisClient -->"

	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		DialTimeout:  3 * time.Duration(redisConf.TimeoutMs) * time.Millisecond,
		ReadTimeout:  time.Duration(redisConf.TimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(redisConf.TimeoutMs) * time.Millisecond,
		PoolSize:     redisConf.PoolSize,
		PoolTimeout:  2 * time.Duration(redisConf.TimeoutMs) * time.Millisecond,
		Password:     redisConf.Password,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		xlog.Errorf(ctx, "%s Ping: %s err: %v", fun, pong, err)
	}

	return &RedisClient{
		RedisCmd: NewRedisCmd(ctx, redisConf, client),
		client:   client,
	}, err
}

/*Do
注意： 不支持mset、msetnx命令；
因为即使在这里实现，也需要调用方区别于其他命令，不如直接调用mset接口。
因为需要把key包装一下，所以这里把cmd、key拆出来。
命令的拼接逻辑是cmd、 keys、otherArgs按顺序拼接。
*/
func (m *RedisClient) Do(ctx context.Context, cmd string, keys []string, otherArgs ...interface{}) (s interface{}, err error) {
	command := "redisext.Do"
	_ = m.WrapperLog(ctx, command, func() error {

		var rcmd *redis.Cmd
		var prefixKey []string
		for _, v := range keys {
			prefixKey = append(prefixKey, m.WrapKey(v))
		}
		args := []interface{}{cmd}
		for _, key := range prefixKey {
			args = append(args, key)
		}
		args = append(args, otherArgs...)

		rcmd = m.client.Do(ctx, args...)
		s, err = rcmd.Result()
		return err
	})
	return
}

func (m *RedisClient) Pipeline(ctx context.Context) (pipe *RedisPipelineClient, err error) {
	command := "RedisCmd.Pipeline"
	_ = m.WrapperLog(ctx, command, func() error {
		p := m.client.Pipeline()
		pipe = &RedisPipelineClient{
			RedisCmd: &RedisCmd{
				RedisConfig: m.RedisConfig,
				client:      p,
			},
			client: p,
		}
		return err
	})
	return
}
func (m *RedisClient) Subscribe(ctx context.Context, channels ...string) (p *redis.PubSub, err error) {
	command := "RedisCmd.Subscribe"
	_ = m.WrapperLog(ctx, command, func() error {
		p = m.client.Subscribe(ctx, channels...)
		return err
	})
	return
}
func (m *RedisClient) PSubscribe(ctx context.Context, channels ...string) (p *redis.PubSub, err error) {
	command := "RedisCmd.PSubscribe"
	_ = m.WrapperLog(ctx, command, func() error {
		p = m.client.PSubscribe(ctx, channels...)
		return err
	})
	return
}

func (m *RedisClient) GetRawClient(ctx context.Context) *redis.Client {
	return m.client
}
