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
 @Time    : 2025/4/15 -- 13:40
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xid xid/uuid.go
*/

package xid

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var uuidMu sync.Mutex

type UUID string

func GetUUID() (UUID, error) {
	uuidMu.Lock()
	defer uuidMu.Unlock()

	uuidGen, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return UUID(uuidGen.String()), nil
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Md5() string {
	h := md5.Sum([]byte(u))
	return fmt.Sprintf("%x", h)
}

func (u UUID) Sha256() string {
	h := sha1.Sum([]byte(u))
	return fmt.Sprintf("%x", h)
}

func (u UUID) Sha1() string {
	h := sha1.Sum([]byte(u))
	return fmt.Sprintf("%x", h)
}
