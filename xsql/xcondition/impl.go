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
 @Time    : 2025/7/1 -- 17:21
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcondition xsql/xcondition/impl.go
*/

package xcondition

import (
	"fmt"
	"strings"

	"github.com/xneogo/matrix/msql"
)

type Condition struct {
	cons map[string]interface{}
}

func New() msql.ConditionsProxy {
	return &Condition{
		cons: make(map[string]interface{}),
	}
}

// NewChanges is same as NewCondition used for updates
func NewChanges() msql.ConditionsProxy {
	return &Condition{
		cons: make(map[string]interface{}),
	}
}

func (c *Condition) Set(column string, value interface{}) msql.ConditionsProxy {
	c.cons[column] = value
	return c
}
func (c *Condition) Export() map[string]interface{} {
	if c.cons == nil {
		c.cons = make(map[string]interface{})
	}
	return c.cons
}
func (c *Condition) ToString() string {
	var cond []string
	for k, v := range c.cons {
		cond = append(cond, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(cond, "-")
}
