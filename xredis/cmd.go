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
 @Time    : 2024/11/5 -- 18:19
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: cmd.go
*/

package xredis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xneogo/Z1ON0101/xtrace"
	"github.com/xneogo/extensions/xtime"
)

// RedisConfig cache & redis apollo config
type RedisConfig struct {
	CacheConfig
	Addr      string `json:"addr"`
	PoolSize  int    `json:"poolSize"`
	TimeoutMs int64  `json:"timeout_ms"`
	Password  string `json:"password"`
	Cluster   string `json:"cluster"`
}
type RedisCmd struct {
	RedisConfig
	client redis.Cmdable
}

func (m *RedisCmd) Close(ctx context.Context) (err error) {
	return nil
}

func NewRedisCmd(ctx context.Context, redisConf *RedisConfig, client redis.Cmdable) *RedisCmd {
	return &RedisCmd{
		RedisConfig: *redisConf,
		client:      client,
	}
}

type Z struct {
	Score  float64
	Member interface{}
}

func (z Z) toRedisZ() redis.Z {
	return redis.Z{
		Score:  z.Score,
		Member: z.Member,
	}
}

func fromRedisZ(rz redis.Z) Z {
	return Z{
		Score:  rz.Score,
		Member: rz.Member,
	}
}

func toRedisZSlice(zs []Z) (rzs []redis.Z) {
	for _, z := range zs {
		rzs = append(rzs, z.toRedisZ())
	}
	return
}

func fromRedisZSlice(rzs []redis.Z) (zs []Z) {
	for _, rz := range rzs {
		zs = append(zs, fromRedisZ(rz))
	}
	return
}

type ZRangeBy struct {
	Min, Max      string
	Offset, Count int64
}

func toRedisZRangeBy(by ZRangeBy) *redis.ZRangeBy {
	return &redis.ZRangeBy{
		Min:    by.Min,
		Max:    by.Max,
		Offset: by.Offset,
		Count:  by.Count,
	}
}

func (m *RedisCmd) WrapKey(key string) string {
	return fmt.Sprintf("%s.%s", m.buildPre(), key)
}
func (m *RedisCmd) buildPre() string {
	var pre string
	if len(m.Prefix) > 0 {
		pre = fmt.Sprintf("%s.%s", m.Namespace, m.Prefix)
	} else {
		pre = m.Namespace
	}
	return pre
}
func (m *RedisCmd) UnWrapKey(key string) string {
	var pre = m.buildPre()
	if strings.HasPrefix(key, pre) {
		return key[len(pre):]
	}
	return key
}

func (m *RedisCmd) GetString(ctx context.Context, key string) (s string, err error) {
	command := "RedisCmd.Get"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.Get(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

// Get
// 返回的虽然是interface，但是redis是明确的string类型。
// 存入数据时，会对数据进行序列化成字符串，因此返回的都是string
func (m *RedisCmd) Get(ctx context.Context, key string) (s interface{}, err error) {
	command := "RedisCmd.Get"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.Get(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) MGet(ctx context.Context, keys ...string) (v []interface{}, err error) {
	command := "RedisCmd.MGet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var prefixKey = make([]string, 0, len(keys))
		for _, v := range keys {
			prefixKey = append(prefixKey, m.WrapKey(v))
		}
		v, err = m.client.MGet(ctx, prefixKey...).Result()
		return err
	})
	return
}

func (m *RedisCmd) Set(ctx context.Context, key string, val interface{}, exp time.Duration) (s string, err error) {
	command := "RedisCmd.Set"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.Set(ctx, m.WrapKey(key), val, exp).Result()
		return err
	})
	return
}

func (m *RedisCmd) Append(ctx context.Context, key, val string) (n int64, err error) {
	command := "RedisCmd.Append"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.Append(ctx, m.WrapKey(key), val).Result()
		return err
	})
	return
}

func (m *RedisCmd) MSet(ctx context.Context, pairs ...interface{}) (s string, err error) {
	command := "RedisCmd.MSet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var prefixPairs = make([]interface{}, 0, len(pairs))
		for k, v := range pairs {
			if (k & 1) == 0 {
				prefixPairs = append(prefixPairs, m.WrapKey(v.(string)))
			} else {
				prefixPairs = append(prefixPairs, v)
			}
		}
		s, err = m.client.MSet(ctx, prefixPairs...).Result()
		return err
	})
	return
}

func (m *RedisCmd) GetBit(ctx context.Context, key string, offset int64) (n int64, err error) {
	command := "RedisCmd.GetBit"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.GetBit(ctx, m.WrapKey(key), offset).Result()
		return err
	})
	return
}

func (m *RedisCmd) SetBit(ctx context.Context, key string, offset int64, value int) (n int64, err error) {
	command := "RedisCmd.SetBit"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.SetBit(ctx, m.WrapKey(key), offset, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) Incr(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.Incr"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.Incr(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) IncrBy(ctx context.Context, key string, val int64) (n int64, err error) {
	command := "RedisCmd.IncrBy"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.IncrBy(ctx, m.WrapKey(key), val).Result()
		return err
	})
	return
}

func (m *RedisCmd) Decr(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.Decr"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.Decr(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) DecrBy(ctx context.Context, key string, val int64) (n int64, err error) {
	command := "RedisCmd.DecrBy"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.DecrBy(ctx, m.WrapKey(key), val).Result()
		return err
	})
	return
}

func (m *RedisCmd) SetNX(ctx context.Context, key string, val interface{}, exp time.Duration) (b bool, err error) {
	command := "RedisCmd.SetNX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.SetNX(ctx, m.WrapKey(key), val, exp).Result()
		return err
	})
	return
}

func (m *RedisCmd) Exists(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.Exists"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.Exists(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) Del(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.Del"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.Del(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) Type(ctx context.Context, key string) (s string, err error) {
	command := "RedisCmd.Type"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.Type(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) Expire(ctx context.Context, key string, expiration time.Duration) (b bool, err error) {
	command := "RedisCmd.Expire"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.Expire(ctx, m.WrapKey(key), expiration).Result()
		return err
	})
	return
}

// HSet hashes apis
func (m *RedisCmd) HSet(ctx context.Context, key string, field string, value interface{}) (b int64, err error) {
	command := "RedisCmd.HSet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.HSet(ctx, m.WrapKey(key), field, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) HDel(ctx context.Context, key string, fields ...string) (n int64, err error) {
	command := "RedisCmd.HDel"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.HDel(ctx, m.WrapKey(key), fields...).Result()
		return err
	})
	return
}

func (m *RedisCmd) HExists(ctx context.Context, key string, field string) (b bool, err error) {
	command := "RedisCmd.HExists"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.HExists(ctx, m.WrapKey(key), field).Result()
		return err
	})
	return
}

func (m *RedisCmd) HGet(ctx context.Context, key string, field string) (s string, err error) {
	command := "RedisCmd.HGet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.HGet(ctx, m.WrapKey(key), field).Result()
		return err
	})
	return
}

func (m *RedisCmd) HGetAll(ctx context.Context, key string) (sm map[string]string, err error) {
	command := "RedisCmd.HGetAll"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		sm, err = m.client.HGetAll(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) HIncrBy(ctx context.Context, key string, field string, incr int64) (n int64, err error) {
	command := "RedisCmd.HIncrBy"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.HIncrBy(ctx, m.WrapKey(key), field, incr).Result()
		return err
	})
	return
}

func (m *RedisCmd) HIncrByFloat(ctx context.Context, key string, field string, incr float64) (f float64, err error) {
	command := "RedisCmd.HIncrByFloat"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		f, err = m.client.HIncrByFloat(ctx, m.WrapKey(key), field, incr).Result()
		return err
	})
	return
}

func (m *RedisCmd) HKeys(ctx context.Context, key string) (ss []string, err error) {
	command := "RedisCmd.HKeys"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.HKeys(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) HLen(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.HLen"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.HLen(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) HMGet(ctx context.Context, key string, fields ...string) (vs []interface{}, err error) {
	command := "RedisCmd.HMGet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		vs, err = m.client.HMGet(ctx, m.WrapKey(key), fields...).Result()
		return err
	})
	return
}

func (m *RedisCmd) HMSet(ctx context.Context, key string, fields map[string]interface{}) (s bool, err error) {
	command := "RedisCmd.HMSet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.HMSet(ctx, m.WrapKey(key), fields).Result()
		return err
	})
	return
}

func (m *RedisCmd) HSetNX(ctx context.Context, key string, field string, val interface{}) (b bool, err error) {
	command := "RedisCmd.HSetNX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.HSetNX(ctx, m.WrapKey(key), field, val).Result()
		return err
	})
	return
}

func (m *RedisCmd) HVals(ctx context.Context, key string) (ss []string, err error) {
	command := "RedisCmd.HVals"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.HVals(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

// ZAdd sorted set apis
func (m *RedisCmd) ZAdd(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "RedisCmd.ZAdd"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZAdd(ctx, m.WrapKey(key), toRedisZSlice(members)...).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZAddNX(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "RedisCmd.ZAddNX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZAddNX(ctx, m.WrapKey(key), toRedisZSlice(members)...).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZAddXX(ctx context.Context, key string, members []Z) (n int64, err error) {
	command := "RedisCmd.ZAddXX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZAddXX(ctx, m.WrapKey(key), toRedisZSlice(members)...).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZCard(ctx context.Context, key string) (n int64, err error) {
	command := "RedisCmd.ZCard"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZCard(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZCount(ctx context.Context, key, min, max string) (n int64, err error) {
	command := "RedisCmd.ZCount"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZCount(ctx, m.WrapKey(key), min, max).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	command := "RedisCmd.ZRange"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.ZRange(ctx, m.WrapKey(key), start, stop).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRangeByLex(ctx context.Context, key string, by ZRangeBy) (ss []string, err error) {
	command := "RedisCmd.ZRangeByLex"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.ZRangeByLex(ctx, m.WrapKey(key), toRedisZRangeBy(by)).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRangeByScore(ctx context.Context, key string, by ZRangeBy) (ss []string, err error) {
	command := "RedisCmd.ZRangeByScore"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.ZRangeByScore(ctx, m.WrapKey(key), toRedisZRangeBy(by)).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRangeByScoreWithScores(ctx context.Context, key string, by ZRangeBy) (zs []Z, err error) {
	command := "RedisCmd.ZRangeByScoreWithScores"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var rzs []redis.Z
		rzs, err = m.client.ZRangeByScoreWithScores(ctx, m.WrapKey(key), toRedisZRangeBy(by)).Result()
		zs = fromRedisZSlice(rzs)
		return err
	})
	return
}

func (m *RedisCmd) ZRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	command := "RedisCmd.ZRangeWithScores"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var rzs []redis.Z
		rzs, err = m.client.ZRangeWithScores(ctx, m.WrapKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
		return err
	})
	return
}

func (m *RedisCmd) ZRevRange(ctx context.Context, key string, start, stop int64) (ss []string, err error) {
	command := "RedisCmd.ZRevRange"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.ZRevRange(ctx, m.WrapKey(key), start, stop).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) (zs []Z, err error) {
	command := "RedisCmd.ZRevRangeWithScores"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var rzs []redis.Z
		rzs, err = m.client.ZRevRangeWithScores(ctx, m.WrapKey(key), start, stop).Result()
		zs = fromRedisZSlice(rzs)
		return err
	})
	return
}

func (m *RedisCmd) ZRevRangeByScoreWithScores(ctx context.Context, key string, by ZRangeBy) (zs []Z, err error) {
	command := "RedisCmd.ZRevRangeByScoreWithScores"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var rzs []redis.Z
		rzs, err = m.client.ZRevRangeByScoreWithScores(ctx, m.WrapKey(key), toRedisZRangeBy(by)).Result()
		zs = fromRedisZSlice(rzs)
		return err
	})
	return
}

func (m *RedisCmd) ZRevRangeByScore(ctx context.Context, key string, by ZRangeBy) (ss []string, err error) {
	command := "RedisCmd.ZRevRangeByScore"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		ss, err = m.client.ZRevRangeByScore(ctx, m.WrapKey(key), toRedisZRangeBy(by)).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRank(ctx context.Context, key string, member string) (n int64, err error) {
	command := "RedisCmd.ZRank"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZRank(ctx, m.WrapKey(key), member).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRevRank(ctx context.Context, key string, member string) (n int64, err error) {
	command := "RedisCmd.ZRevRank"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZRevRank(ctx, m.WrapKey(key), member).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRem(ctx context.Context, key string, members []interface{}) (n int64, err error) {
	command := "RedisCmd.ZRem"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		n, err = m.client.ZRem(ctx, m.WrapKey(key), members).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRemRangeByScore(ctx context.Context, key, min, max string) (i int64, err error) {
	command := "RedisCmd.ZRemRangeByScore"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		i, err = m.client.ZRemRangeByScore(ctx, m.WrapKey(key), min, max).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) (i int64, err error) {
	command := "RedisCmd.ZRemRangeByRank"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		i, err = m.client.ZRemRangeByRank(ctx, m.WrapKey(key), start, stop).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZIncrBy(ctx context.Context, key string, increment float64, member string) (f float64, err error) {
	command := "RedisCmd.ZIncrBy"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		f, err = m.client.ZIncrBy(ctx, m.WrapKey(key), increment, member).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZScore(ctx context.Context, key string, member string) (f float64, err error) {
	command := "RedisCmd.ZScore"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		f, err = m.client.ZScore(ctx, m.WrapKey(key), member).Result()
		return err
	})
	return
}

func (m *RedisCmd) TTL(ctx context.Context, key string) (d time.Duration, err error) {
	command := "RedisCmd.TTL"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		d, err = m.client.TTL(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, rcursor uint64, err error) {
	command := "RedisCmd.SScan"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		keys, rcursor, err = m.client.SScan(ctx, m.WrapKey(key), cursor, match, count).Result()
		return err
	})
	return
}

func (m *RedisCmd) SAdd(ctx context.Context, key string, members ...interface{}) (i int64, err error) {
	command := "RedisCmd.SAdd"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		i, err = m.client.SAdd(ctx, m.WrapKey(key), members...).Result()
		return err
	})
	return
}

func (m *RedisCmd) SPop(ctx context.Context, key string) (s string, err error) {
	command := "RedisCmd.SPop"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.SPop(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) SPopN(ctx context.Context, key string, count int64) (s []string, err error) {
	command := "RedisCmd.SPopN"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.SPopN(ctx, m.WrapKey(key), count).Result()
		return err
	})
	return
}

func (m *RedisCmd) SRem(ctx context.Context, key string, members ...interface{}) (i int64, err error) {
	command := "RedisCmd.SRem"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		i, err = m.client.SRem(ctx, m.WrapKey(key), members...).Result()
		return err
	})
	return
}

func (m *RedisCmd) SCard(ctx context.Context, key string) (i int64, err error) {
	command := "RedisCmd.SCard"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		i, err = m.client.SCard(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) SIsMember(ctx context.Context, key string, member interface{}) (b bool, err error) {
	command := "RedisCmd.SIsMember"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		b, err = m.client.SIsMember(ctx, m.WrapKey(key), member).Result()
		return err
	})
	return
}

func (m *RedisCmd) SMembers(ctx context.Context, key string) (s []string, err error) {
	command := "RedisCmd.SMembers"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.SMembers(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) SRandMember(ctx context.Context, key string) (s string, err error) {
	command := "RedisCmd.SRandMember"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.SRandMember(ctx, m.WrapKey(key)).Result()
		return err
	})
	return
}

func (m *RedisCmd) SRandMemberN(ctx context.Context, key string, count int64) (s []string, err error) {
	command := "RedisCmd.SRandMemberN"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		s, err = m.client.SRandMemberN(ctx, m.WrapKey(key), count).Result()
		return err
	})
	return
}

func (m *RedisCmd) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, rcursor uint64, err error) {
	command := "RedisCmd.ZScan"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		keys, rcursor, err = m.client.ZScan(ctx, m.WrapKey(key), cursor, match, count).Result()
		return err
	})
	return
}

func (m *RedisCmd) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) (keys []string, rcursor uint64, err error) {
	command := "RedisCmd.HScan"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		keys, rcursor, err = m.client.HScan(ctx, m.WrapKey(key), cursor, match, count).Result()
		return err
	})
	return
}

func (m *RedisCmd) SInter(ctx context.Context, keys ...string) (result []string, err error) {
	command := "RedisCmd.SInter"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var prefixKey = make([]string, 0, len(keys))
		for _, v := range keys {
			prefixKey = append(prefixKey, m.WrapKey(v))
		}
		result, err = m.client.SInter(ctx, prefixKey...).Result()
		return err
	})
	return
}

func (m *RedisCmd) PFAdd(ctx context.Context, key string, els ...interface{}) (result int64, err error) {
	command := "RedisCmd.PFAdd"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		result, err = m.client.PFAdd(ctx, m.WrapKey(key), els...).Result()
		return err
	})
	return
}

func (m *RedisCmd) PFMerge(ctx context.Context, key string, keys ...string) (result string, err error) {
	command := "RedisCmd.PFMerge"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var prefixKey = make([]string, 0, len(keys))
		for _, v := range keys {
			prefixKey = append(prefixKey, m.WrapKey(v))
		}
		result, err = m.client.PFMerge(ctx, m.WrapKey(key), prefixKey...).Result()
		return err
	})
	return
}

func (m *RedisCmd) PFCount(ctx context.Context, keys ...string) (result int64, err error) {
	command := "RedisCmd.PFCount"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		var prefixKey = make([]string, 0, len(keys))
		for _, v := range keys {
			prefixKey = append(prefixKey, m.WrapKey(v))
		}
		result, err = m.client.PFCount(ctx, prefixKey...).Result()
		return err
	})
	return
}

func (m *RedisCmd) LIndex(ctx context.Context, key string, index int64) (result string, err error) {

	command := "RedisCmd.LIndex"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LIndex(ctx, k, index).Result()
		return err
	})
	return
}

func (m *RedisCmd) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (result int64, err error) {
	command := "RedisCmd.LInsert"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LInsert(ctx, k, op, pivot, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) LLen(ctx context.Context, key string) (result int64, err error) {
	command := "RedisCmd.LLen"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LLen(ctx, k).Result()
		return err
	})
	return
}

func (m *RedisCmd) LPop(ctx context.Context, key string) (result string, err error) {
	command := "RedisCmd.LPop"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LPop(ctx, k).Result()
		return err
	})
	return
}
func (m *RedisCmd) BLPop(ctx context.Context, timeout time.Duration, key string) (result string, err error) {
	command := "RedisCmd.BLPop"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		results, er := m.client.BLPop(ctx, timeout, m.WrapKey(key)).Result()
		if er != nil {
			err = er
			return er
		}
		// 没有等到数据，不报错，返回空数据
		if len(results) == 0 {
			return nil
		}
		if len(results) > 1 {
			result = results[1]
			return nil
		}
		return nil
	})
	return
}

func (m *RedisCmd) LPush(ctx context.Context, key string, values ...interface{}) (result int64, err error) {
	command := "RedisCmd.LPush"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LPush(ctx, k, values...).Result()
		return err
	})
	return
}

func (m *RedisCmd) LPushX(ctx context.Context, key string, value interface{}) (result int64, err error) {
	command := "RedisCmd.LPushX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LPushX(ctx, k, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) LRange(ctx context.Context, key string, start, stop int64) (s []string, err error) {
	command := "RedisCmd.LRange"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		s, err = m.client.LRange(ctx, k, start, stop).Result()
		return err
	})
	return
}

func (m *RedisCmd) LRem(ctx context.Context, key string, count int64, value interface{}) (result int64, err error) {
	command := "RedisCmd.LRem"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LRem(ctx, k, count, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) LSet(ctx context.Context, key string, index int64, value interface{}) (result string, err error) {
	command := "RedisCmd.LSet"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LSet(ctx, k, index, value).Result()
		return err
	})
	return
}

func (m *RedisCmd) LTrim(ctx context.Context, key string, start, stop int64) (result string, err error) {
	command := "RedisCmd.LTrim"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.LTrim(ctx, k, start, stop).Result()
		return err
	})
	return
}

func (m *RedisCmd) RPop(ctx context.Context, key string) (result string, err error) {
	command := "RedisCmd.RPop"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.RPop(ctx, k).Result()
		return err
	})
	return
}

func (m *RedisCmd) BRPop(ctx context.Context, timeout time.Duration, key string) (result string, err error) {
	command := "RedisCmd.BRPop"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		results, er := m.client.BRPop(ctx, timeout, m.WrapKey(key)).Result()
		if er != nil {
			err = er
			return er
		}
		// 没有等到数据，不报错，返回空数据
		if len(results) == 0 {
			return nil
		}
		if len(results) > 1 {
			result = results[1]
			return nil
		}
		return nil
	})
	return
}

func (m *RedisCmd) RPush(ctx context.Context, key string, values ...interface{}) (result int64, err error) {
	command := "RedisCmd.RPush"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.RPush(ctx, k, values...).Result()
		return err
	})
	return
}

func (m *RedisCmd) RPushX(ctx context.Context, key string, value interface{}) (result int64, err error) {
	command := "RedisCmd.RPushX"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		k := m.WrapKey(key)
		result, err = m.client.RPushX(ctx, k, value).Result()
		return err
	})
	return
}

// ScriptLoad load script to redis server
func (m *RedisCmd) ScriptLoad(ctx context.Context, script *Script) (r string, err error) {
	command := "RedisClient.ScriptLoad"
	_ = m.WrapperLog(ctx, command, func() error {
		r, err = m.client.ScriptLoad(ctx, script.src).Result()
		return err
	})
	return
}

// ScriptExists check if script exists in redis server
func (m *RedisCmd) ScriptExists(ctx context.Context, script *Script) (r bool, err error) {
	command := "RedisClient.ScriptExists"
	_ = m.WrapperLog(ctx, command, func() error {
		result, err := m.client.ScriptExists(ctx, script.src).Result()
		if err != nil {
			r = false
			return err
		}
		if len(result) > 0 {
			r = result[0]
		}
		return err
	})
	return
}

// Eval exec with script
func (m *RedisCmd) Eval(ctx context.Context, script *Script, keys []string, args ...interface{}) (r interface{}, err error) {
	command := "RedisClient.Eval"
	_ = m.WrapperLog(ctx, command, func() error {
		for i, key := range keys {
			keys[i] = m.WrapKey(key)
		}
		if err == nil {
			r, err = m.client.Eval(ctx, script.src, keys, args...).Result()
		}
		return err
	})
	return
}

// EvalSha exec with script hash
func (m *RedisCmd) EvalSha(ctx context.Context, script *Script, keys []string, args ...interface{}) (r interface{}, err error) {
	command := "RedisClient.EvalSha"
	_ = m.WrapperLog(ctx, command, func() error {

		for i, key := range keys {
			keys[i] = m.WrapKey(key)
		}
		if err == nil {
			r, err = m.client.EvalSha(ctx, script.hash, keys, args...).Result()
		}
		return err
	})
	return
}

func (m *RedisCmd) WrapperLog(ctx context.Context, command string, executor func() error) error {
	span, ctx := xtrace.StartSpanFromContext(ctx, command)
	st := xtime.NewTimeStat()
	defer func() {
		span.Finish()
		statReqDuration(ctx, m.Namespace, command, m.Cluster, st.Millisecond())
	}()
	err := executor()
	statReqErr(ctx, m.Namespace, command, m.Cluster, err)
	return err
}

func (m *RedisCmd) Publish(ctx context.Context, channel string, value interface{}) (result int64, err error) {
	command := "RedisCmd.Publish"
	_ = WrapMetric(ctx, RedisMetric{
		Command:   command,
		Namespace: m.Namespace,
		Cluster:   m.Cluster,
	}, func() error {
		result, err = m.client.Publish(ctx, channel, value).Result()
		return err
	})
	return
}
