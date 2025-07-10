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
 @Time    : 2024/11/5 -- 18:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: watchdog.go
*/

package xredis

import (
	"context"
	"errors"
	"fmt"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix"
	"math/rand"
	"time"
)

const (
	durGap       = time.Millisecond * 20
	secretWindow = 1000 // 默认，如果有并发量特别大的接口，这个 Window 需要调大且支持自定义
)

type XWatchDog struct {
	cli    *RedisClient
	dur    time.Duration
	k      string
	from   string
	secret int64
	window int64
	ctx    context.Context
	cancel context.CancelFunc
}

func NewXWatchDog(client *RedisClient, window int64, from string) matrix.Woof {
	if window <= secretWindow {
		window = secretWindow
	}
	dog := &XWatchDog{
		cli:    client,
		window: window,
		from:   from,
	}

	return dog
}

func (x *XWatchDog) Lock(ctx context.Context, k string, _ interface{}, dur time.Duration) error {
	fun := x.from + "XWatchDog.Lock ->"
	x.dur = dur
	x.ctx, x.cancel = context.WithCancel(ctx)

	if x.secret == 0 {
		x.secret = rand.Int63n(x.window)
	}
	lock, err := x.cli.Lock(ctx, k, x.secret, dur)
	if err != nil {
		xlog.Errorf(ctx, "%s key:%s try lock fail:%s", fun, k, err)
		return err
	}
	if !lock {
		return errors.New("already locked")
	}
	go x.Watch(ctx, k, dur)
	return err
}

func (x *XWatchDog) Unlock(ctx context.Context, k string) error {
	fun := x.from + "XWatchDog.Unlock ->"
	r, err := x.cli.Unlock(x.ctx, x.k, x.secret)
	if err != nil {
		xlog.Warnf(x.ctx, "%s key:%s try unlock fail:%t,%s", fun, x.k, r, err)
	}
	x.cancel()
	return nil
}

func (x *XWatchDog) Watch(ctx context.Context, k string, dur time.Duration) {
	ticker := time.NewTicker(dur - durGap)
	defer ticker.Stop()

	for {
		select {
		case <-x.ctx.Done():
			// try unlock last time
			fmt.Println("exiting watch process")
			_ = x.Unlock(ctx, k)
			return
		case <-ticker.C:
			fmt.Println("relock")
			_ = x.Lock(ctx, k, x.secret, dur)
			ticker.Reset(dur - durGap)
		}
	}
}
