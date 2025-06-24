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
 @Time    : 2025/4/28 -- 13:51
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: parser xcron/parser/parser.go
*/

package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CronField 表示 cron 表达式中的一个字段
type CronField struct {
	Name      string
	Min       int
	Max       int
	Mandatory bool
	Values    []int
}

// CronExpression 表示一个完整的 cron 表达式
type CronExpression struct {
	Seconds    *CronField
	Minutes    *CronField
	Hours      *CronField
	DayOfMonth *CronField
	Month      *CronField
	DayOfWeek  *CronField
	Year       *CronField
	Expression string
	NumFields  int
}

// 定义错误
var (
	ErrInvalidExpression = errors.New("[xcron] 无效的 cron 表达式")
	ErrInvalidField      = errors.New("[xcron] 无效的字段值")
	ErrMandatoryField    = errors.New("[xcron] 必填字段缺失")
)

// NewCronExpression 创建新的 CronExpression
func NewCronExpression(expr string) (*CronExpression, error) {
	// 初始化 cron 表达式对象
	cronExpr := &CronExpression{
		Expression: expr,
		Seconds: &CronField{
			Name:      "Seconds",
			Min:       0,
			Max:       59,
			Mandatory: false,
		},
		Minutes: &CronField{
			Name:      "Minutes",
			Min:       0,
			Max:       59,
			Mandatory: true,
		},
		Hours: &CronField{
			Name:      "Hours",
			Min:       0,
			Max:       23,
			Mandatory: true,
		},
		DayOfMonth: &CronField{
			Name:      "DayOfMonth",
			Min:       1,
			Max:       31,
			Mandatory: true,
		},
		Month: &CronField{
			Name:      "Month",
			Min:       1,
			Max:       12,
			Mandatory: true,
		},
		DayOfWeek: &CronField{
			Name:      "DayOfWeek",
			Min:       0,
			Max:       7, // 0 和 7 都表示星期日
			Mandatory: true,
		},
		Year: &CronField{
			Name:      "Year",
			Min:       1970,
			Max:       2099,
			Mandatory: false,
		},
	}

	// 解析表达式
	return cronExpr.Parse()
}

// Parse 解析 cron 表达式
func (c *CronExpression) Parse() (*CronExpression, error) {
	// 分割表达式
	fields := strings.Fields(c.Expression)
	c.NumFields = len(fields)

	// 根据字段数量处理
	switch c.NumFields {
	case 5:
		// 5个字段: 分钟 小时 日 月 星期
		// 添加默认秒(0)和默认年(*)
		fields = append([]string{"0"}, fields...)
		fields = append(fields, "*")
	case 6:
		// 6个字段: 秒 分钟 小时 日 月 星期
		// 添加默认年(*)
		fields = append(fields, "*")
	case 7:
		// 7个字段: 秒 分钟 小时 日 月 星期 年
		// 完整表达式，不需要处理
	default:
		return nil, fmt.Errorf("%w: 字段数量错误", ErrInvalidExpression)
	}

	// 解析每个字段
	var err error
	if c.Seconds.Values, err = c.parseField(fields[0], c.Seconds); err != nil {
		return nil, err
	}
	if c.Minutes.Values, err = c.parseField(fields[1], c.Minutes); err != nil {
		return nil, err
	}
	if c.Hours.Values, err = c.parseField(fields[2], c.Hours); err != nil {
		return nil, err
	}
	if c.DayOfMonth.Values, err = c.parseField(fields[3], c.DayOfMonth); err != nil {
		return nil, err
	}
	if c.Month.Values, err = c.parseField(fields[4], c.Month); err != nil {
		return nil, err
	}
	if c.DayOfWeek.Values, err = c.parseField(fields[5], c.DayOfWeek); err != nil {
		return nil, err
	}
	if c.Year.Values, err = c.parseField(fields[6], c.Year); err != nil {
		return nil, err
	}

	return c, nil
}

// parseField 解析单个字段
func (c *CronExpression) parseField(fieldExpr string, field *CronField) ([]int, error) {
	// 检查必填字段
	if field.Mandatory && (fieldExpr == "" || fieldExpr == "?") {
		return nil, fmt.Errorf("%w: %s 是必填字段 %+v, %s", ErrMandatoryField, field.Name, field, fieldExpr)
	}

	// 处理月份和星期的字母表示
	if field.Name == "Month" {
		fieldExpr = replaceMonthNames(fieldExpr)
	} else if field.Name == "DayOfWeek" {
		fieldExpr = replaceDayOfWeekNames(fieldExpr)
	}

	// 处理特殊字符
	if fieldExpr == "*" {
		// 星号表示所有可能的值
		return generateRange(field.Min, field.Max), nil
	}

	// 处理复杂表达式
	values := make(map[int]bool)

	// 分割逗号分隔的部分
	parts := strings.Split(fieldExpr, ",")
	for _, part := range parts {
		// 处理每个部分
		nums, err := c.parsePart(part, field)
		if err != nil {
			return nil, err
		}

		// 添加到结果集
		for _, num := range nums {
			values[num] = true
		}
	}

	// 转换为切片
	result := make([]int, 0, len(values))
	for value := range values {
		result = append(result, value)
	}

	return result, nil
}

// parsePart 解析表达式的一部分
func (c *CronExpression) parsePart(part string, field *CronField) ([]int, error) {
	// 处理范围和步长
	if strings.Contains(part, "/") {
		// 处理步长
		return c.parseStepValue(part, field)
	} else if strings.Contains(part, "-") {
		// 处理范围
		return c.parseRangeValue(part, field)
	} else if part == "L" && field.Name == "DayOfMonth" {
		// 处理 L (最后一天)
		return []int{-1}, nil // 使用 -1 表示最后一天
	} else if strings.HasSuffix(part, "L") && field.Name == "DayOfWeek" {
		// 处理 nL (每月最后一个星期n)
		dayStr := strings.TrimSuffix(part, "L")
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < field.Min || day > field.Max {
			return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, part)
		}
		return []int{-day - 1}, nil // 使用负数表示特殊情况
	} else if strings.Contains(part, "#") && field.Name == "DayOfWeek" {
		// 处理 n#m (每月第m个星期n)
		return c.parseHashValue(part, field)
	}

	// 处理单个数值
	num, err := strconv.Atoi(part)
	if err != nil {
		return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, part)
	}

	// 验证范围
	if num < field.Min || num > field.Max {
		return nil, fmt.Errorf("%w: %s 中的 %d 超出范围 [%d,%d]", ErrInvalidField, field.Name, num, field.Min, field.Max)
	}

	return []int{num}, nil
}

// parseStepValue 解析步长表达式 (如 */5, 2-10/2)
func (c *CronExpression) parseStepValue(expr string, field *CronField) ([]int, error) {
	parts := strings.Split(expr, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s 中的 %s 格式错误", ErrInvalidField, field.Name, expr)
	}

	// 解析步长
	step, err := strconv.Atoi(parts[1])
	if err != nil || step <= 0 {
		return nil, fmt.Errorf("%w: %s 中的步长 %s 无效", ErrInvalidField, field.Name, parts[1])
	}

	// 解析范围
	var start, end int
	if parts[0] == "*" {
		start = field.Min
		end = field.Max
	} else if strings.Contains(parts[0], "-") {
		rangeParts := strings.Split(parts[0], "-")
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("%w: %s 中的 %s 格式错误", ErrInvalidField, field.Name, parts[0])
		}

		var err error
		start, err = strconv.Atoi(rangeParts[0])
		if err != nil || start < field.Min || start > field.Max {
			return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, rangeParts[0])
		}

		end, err = strconv.Atoi(rangeParts[1])
		if err != nil || end < field.Min || end > field.Max || end < start {
			return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, rangeParts[1])
		}
	} else {
		var err error
		start, err = strconv.Atoi(parts[0])
		if err != nil || start < field.Min || start > field.Max {
			return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, parts[0])
		}
		end = field.Max
	}

	// 生成结果
	result := make([]int, 0)
	for i := start; i <= end; i += step {
		result = append(result, i)
	}

	return result, nil
}

// parseRangeValue 解析范围表达式 (如 1-5)
func (c *CronExpression) parseRangeValue(expr string, field *CronField) ([]int, error) {
	parts := strings.Split(expr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s 中的 %s 格式错误", ErrInvalidField, field.Name, expr)
	}

	// 解析范围
	start, err := strconv.Atoi(parts[0])
	if err != nil || start < field.Min || start > field.Max {
		return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil || end < field.Min || end > field.Max || end < start {
		return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, parts[1])
	}

	// 生成结果
	return generateRange(start, end), nil
}

// parseHashValue 解析 # 表达式 (如 1#3 - 每月第三个星期一)
func (c *CronExpression) parseHashValue(expr string, field *CronField) ([]int, error) {
	parts := strings.Split(expr, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s 中的 %s 格式错误", ErrInvalidField, field.Name, expr)
	}

	// 解析星期几
	dayOfWeek, err := strconv.Atoi(parts[0])
	if err != nil || dayOfWeek < field.Min || dayOfWeek > field.Max {
		return nil, fmt.Errorf("%w: %s 中的 %s 无效", ErrInvalidField, field.Name, parts[0])
	}

	// 解析第几个
	nth, err := strconv.Atoi(parts[1])
	if err != nil || nth < 1 || nth > 5 {
		return nil, fmt.Errorf("%w: %s 中的 %s 无效，必须在 1-5 之间", ErrInvalidField, field.Name, parts[1])
	}

	// 使用特殊编码表示 n#m
	// 100 * dayOfWeek + nth
	return []int{100*dayOfWeek + nth}, nil
}

// 生成范围内的所有整数
func generateRange(start, end int) []int {
	result := make([]int, end-start+1)
	for i := range result {
		result[i] = start + i
	}
	return result
}

// 替换月份名称为数字
func replaceMonthNames(expr string) string {
	monthMap := map[string]string{
		"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4",
		"MAY": "5", "JUN": "6", "JUL": "7", "AUG": "8",
		"SEP": "9", "OCT": "10", "NOV": "11", "DEC": "12",
	}

	for name, num := range monthMap {
		expr = regexp.MustCompile("(?i)"+name).ReplaceAllString(expr, num)
	}

	return expr
}

// 替换星期名称为数字
func replaceDayOfWeekNames(expr string) string {
	dayMap := map[string]string{
		"SUN": "0", "MON": "1", "TUE": "2", "WED": "3",
		"THU": "4", "FRI": "5", "SAT": "6",
	}

	for name, num := range dayMap {
		expr = regexp.MustCompile("(?i)"+name).ReplaceAllString(expr, num)
	}

	return expr
}

// Validate 验证 cron 表达式是否有效
func (c *CronExpression) Validate() error {
	// 检查必填字段
	if c.Minutes == nil || len(c.Minutes.Values) == 0 {
		return fmt.Errorf("%w: Minutes", ErrMandatoryField)
	}
	if c.Hours == nil || len(c.Hours.Values) == 0 {
		return fmt.Errorf("%w: Hours", ErrMandatoryField)
	}
	if c.DayOfMonth == nil || len(c.DayOfMonth.Values) == 0 {
		return fmt.Errorf("%w: DayOfMonth", ErrMandatoryField)
	}
	if c.Month == nil || len(c.Month.Values) == 0 {
		return fmt.Errorf("%w: Month", ErrMandatoryField)
	}
	if c.DayOfWeek == nil || len(c.DayOfWeek.Values) == 0 {
		return fmt.Errorf("%w: DayOfWeek", ErrMandatoryField)
	}

	return nil
}

// GetDescription 获取 cron 表达式的描述
func (c *CronExpression) GetDescription() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Cron 表达式: %s\n", c.Expression))
	sb.WriteString(fmt.Sprintf("秒: %v\n", c.Seconds.Values))
	sb.WriteString(fmt.Sprintf("分钟: %v\n", c.Minutes.Values))
	sb.WriteString(fmt.Sprintf("小时: %v\n", c.Hours.Values))
	sb.WriteString(fmt.Sprintf("日期: %v\n", c.DayOfMonth.Values))
	sb.WriteString(fmt.Sprintf("月份: %v\n", c.Month.Values))
	sb.WriteString(fmt.Sprintf("星期: %v\n", c.DayOfWeek.Values))
	sb.WriteString(fmt.Sprintf("年份: %v\n", c.Year.Values))

	return sb.String()
}
