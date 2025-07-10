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
 @Time    : 2024/11/5 -- 19:04
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: watchdog_test.go
*/

package xredis

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/xneogo/eins/oneredis"
)

var client *RedisClient

func init() {
	dsn := "redis://localhost:6379"
	cli := oneredis.MustConnect(dsn)
	client = &RedisClient{
		RedisCmd: NewRedisCmd(context.Background(), nil, cli),
		client:   cli,
	}
}

func TestWatchdog(t *testing.T) {
	ctx := context.Background()
	doLogic(ctx)
	time.Sleep(time.Second * 2)
	doLogic(ctx)
}

func doLogic(ctx context.Context) {
	dog := NewXWatchDog(client, 0, "test")
	err := dog.Lock(ctx, "doLogic", uuid.New(), time.Second*2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dog.Unlock(ctx, "doLogic")
	time.Sleep(time.Second * 3)
	return
}
