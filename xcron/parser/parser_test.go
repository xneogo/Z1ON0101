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
 @Time    : 2025/4/28 -- 14:13
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: parser xcron/parser/parser_test.go
*/

package parser

import (
	"fmt"
	"testing"
)

func TestCronExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		shouldError bool
	}{
		{"基本表达式", "0 0 12 * * *", false},
		{"带秒的表达式", "30 0 12 * * *", false},
		{"带年的表达式", "0 0 12 * * * 2022", false},
		{"星号表达式", "* * * * * *", false},
		{"步长表达式", "*/5 * * * * *", false},
		{"范围表达式", "10-30 * * * * *", false},
		{"列表表达式", "5,10,15 * * * * *", false},
		{"L表达式", "0 0 12 L * *", false},
		{"星期L表达式", "0 0 12 * * 5L", false},
		{"星期#表达式", "0 0 12 * * 5#3", false},
		{"月份名称", "0 0 12 * JAN-DEC *", false},
		{"星期名称", "0 0 12 * * MON-FRI", false},
		{"复杂表达式", "0 0/5 14,18 * * *", false},
		{"无效字段数", "0 * *", true},
		{"无效范围", "0 0 25 * * ?", true},
		{"无效步长", "0/0 * * * * ?", true},
		{"无效格式", "a b c d e f", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewCronExpression(tt.expression)
			if (err != nil) != tt.shouldError {
				t.Errorf("NewCronExpression(%q) error = %v, shouldError %v", tt.expression, err, tt.shouldError)
			}
			if err == nil {
				fmt.Println(e.GetDescription())
			}
		})
	}
}

func TestSpecificExpressions(t *testing.T) {
	tests := []struct {
		expression string
		seconds    []int
		minutes    []int
		hours      []int
	}{
		{"0 0 12 * * *", []int{0}, []int{0}, []int{12}},
		{"*/5 * * * * *", []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}},
		{"0 0/15 8-10 * * *", []int{0}, []int{0, 15, 30, 45}, []int{8, 9, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			expr, err := NewCronExpression(tt.expression)
			if err != nil {
				t.Fatalf("解析表达式失败: %v", err)
			}

			// 检查秒
			if !compareIntSlices(expr.Seconds.Values, tt.seconds) {
				t.Errorf("秒字段不匹配, 期望 %v, 得到 %v", tt.seconds, expr.Seconds.Values)
			}

			// 检查分钟
			if !compareIntSlices(expr.Minutes.Values, tt.minutes) {
				t.Errorf("分钟字段不匹配, 期望 %v, 得到 %v", tt.minutes, expr.Minutes.Values)
			}

			// 检查小时
			if !compareIntSlices(expr.Hours.Values, tt.hours) {
				t.Errorf("小时字段不匹配, 期望 %v, 得到 %v", tt.hours, expr.Hours.Values)
			}
		})
	}
}

// 比较两个整数切片是否包含相同的元素（不考虑顺序）
func compareIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[int]bool)
	for _, v := range a {
		aMap[v] = true
	}

	for _, v := range b {
		if !aMap[v] {
			return false
		}
	}

	return true
}
