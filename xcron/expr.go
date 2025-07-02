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
 @Time    : 2025/4/28 -- 13:47
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: xcron xcron/expr.go
*/

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
 @Time    : 2025/4/28 -- 14:30
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: expr xcron/expr.go
*/

package xcron

import (
	"github.com/xneogo/Z1ON0101/xcron/parser"
	"time"
)

// Expr 接口定义了计算 cron 表达式下次执行时间的方法
type Expr interface {
	// Next 返回下一次执行时间（从给定时间开始）
	Next(t time.Time) time.Time

	// NextN 返回接下来 n 次执行时间（从给定时间开始）
	// n 不能超过 5, 超过只返回 5 个时间点
	NextN(t time.Time, n int) []time.Time
}

// ExprImpl 实现了 Expr 接口
type ExprImpl struct {
	expression *parser.CronExpression
}

// NewExpr 创建一个新的 Expr 实例
func NewExpr(cronExpression string) (Expr, error) {
	expression, err := parser.NewCronExpression(cronExpression)
	if err != nil {
		return nil, err
	}

	return &ExprImpl{
		expression: expression,
	}, nil
}

// Next 实现 Expr 接口的 Next 方法
func (c *ExprImpl) Next(t time.Time) time.Time {
	// 从当前时间的下一秒开始查找
	t = t.Add(time.Second)

	// 最多尝试查找未来一年内的时间点
	yearLimit := t.AddDate(1, 0, 0)

	for t.Before(yearLimit) {
		if c.matches(t) {
			return t
		}
		t = t.Add(time.Second)
	}

	// 如果一年内找不到匹配的时间点，返回零值
	return time.Time{}
}

func (c *ExprImpl) matches(t time.Time) bool {
	// TODO 实现匹配逻辑
	return false
}

// NextN 实现 Expr 接口的 NextN 方法
func (c *ExprImpl) NextN(t time.Time, n int) []time.Time {
	if n <= 0 || n > 5 {
		return []time.Time{}
	}

	result := make([]time.Time, 0, n)
	current := t

	for i := 0; i < n; i++ {
		next := c.Next(current)
		if next.IsZero() {
			break
		}
		result = append(result, next)
		current = next
	}

	return result
}

// getLastDayOfMonth 获取月份的最后一天
func getLastDayOfMonth(t time.Time) int {
	// 获取下个月的第一天，然后减去一天
	nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	lastDay := nextMonth.Add(-24 * time.Hour)
	return lastDay.Day()
}

// isLastDayOfWeekInMonth 检查是否是月份的最后一个指定星期几
func isLastDayOfWeekInMonth(t time.Time) bool {
	// 获取当前日期的星期几
	dayOfWeek := int(t.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	// 获取当月的最后一天
	lastDay := getLastDayOfMonth(t)

	// 检查从当前日期到月底是否还有相同的星期几
	for day := t.Day() + 7; day <= lastDay; day += 7 {
		return false
	}

	return true
}

// getNthDayOfWeekInMonth 获取当前日期是本月的第几个星期几
func getNthDayOfWeekInMonth(t time.Time) int {
	day := t.Day()
	dayOfWeek := int(t.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	// 计算本月第一个该星期几的日期
	firstDay := day % 7
	if firstDay == 0 {
		firstDay = 7
	}

	// 计算是第几个星期几
	return (day-firstDay)/7 + 1
}
