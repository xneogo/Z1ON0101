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
 @Time    : 2024/10/12 -- 15:47
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: observer.go
*/

package xapollo

import (
	"fmt"
	"github.com/eva-nigouki/apogo"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

type aObserver struct {
	observer *mobserver.ConfigObserver
}

func (o *aObserver) HandleChangeEvent(ce *apogo.ChangeEvent) {
	o.observer.HandleChangeEvent(transApogoChangeEvent(ce))
}

func transApogoChangeEvent(ace *apogo.ChangeEvent) *mobserver.ChangeEvent {
	var changes = map[string]*mobserver.Change{}
	for k, ac := range ace.Changes {
		if c, err := transApogoChange(ac); err == nil {
			changes[k] = c
		}
	}
	return &mobserver.ChangeEvent{
		Source:    mobserver.Apollo,
		Namespace: ace.Namespace,
		Changes:   changes,
	}
}

func transApogoChange(ac *apogo.Change) (change *mobserver.Change, err error) {
	ct, err := transApogoChangeType(ac.ChangeType)
	if err != nil {
		fmt.Printf("transApogoChange err:%s", err.Error())
		return
	}

	change = &mobserver.Change{
		OldValue:   ac.OldValue.(string),
		NewValue:   ac.NewValue.(string),
		ChangeType: ct,
	}
	return
}

func transApogoChangeType(act apogo.ChangeType) (ct mobserver.ChangeType, err error) {
	switch act {
	case apogo.ADD:
		ct = mobserver.ADD
	case apogo.MODIFY:
		ct = mobserver.MODIFY
	case apogo.DELETE:
		ct = mobserver.DELETE
	default:
		err = fmt.Errorf("invalid apollo change type:%+v", act)
	}

	return
}
