package util

import (
	"fmt"
	"strings"
)

// FormatInteger 格式化整数（千位分割）
func FormatInteger(value interface{}) string {
	var num int64
	switch v := value.(type) {
	case int:
		num = int64(v)
	case int32:
		num = int64(v)
	case int64:
		num = v
	case float64:
		num = int64(v)
	default:
		return "0"
	}

	str := fmt.Sprintf("%d", num)
	if len(str) <= 3 {
		return str
	}

	// 从右往左每三位插入逗号
	result := ""
	for i, ch := range strings.Split(str, "") {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += ch
	}
	return result
}

// FormatDecimal 格式化浮点数
func FormatDecimal(value float64, digits int) string {
	formatStr := fmt.Sprintf("%%.%df", digits)
	return fmt.Sprintf(formatStr, value)
}

// FormatPercent 格式化百分比
func FormatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value*100)
}

// FormatCNY 格式化 CNY 金额
func FormatCNY(value float64) string {
	return fmt.Sprintf("¥%.2f", value)
}
