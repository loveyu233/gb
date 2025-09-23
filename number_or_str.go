package gb

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertStringToUint32 将带前导零的字符串转换为uint32
// 参数: str - 输入的字符串，如 "01", "02", "03", "11", "101", "1110"
// 返回: uint32值和错误信息
func ConvertStringToUint32(str string) (uint32, error) {
	// 去除空格
	str = strings.TrimSpace(str)

	// 检查空字符串
	if str == "" {
		return 0, fmt.Errorf("输入字符串为空")
	}

	// 检查是否全为数字
	for _, char := range str {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("字符串包含非数字字符: %s", str)
		}
	}

	// 使用strconv.ParseUint自动处理前导零
	// ParseUint会自动去除前导零并转换为数字
	result, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("转换失败: %v", err)
	}

	return uint32(result), nil
}

// ConvertStringToUint32Simple 简化版本，不返回错误
// 如果转换失败则返回0
func ConvertStringToUint32Simple(str string) uint32 {
	result, err := ConvertStringToUint32(str)
	if err != nil {
		return 0
	}
	return result
}
