package gb

import (
	"fmt"
	"strconv"
	"strings"
)

// GetLastNChars 函数用于处理GetLastNChars相关逻辑。
func GetLastNChars(str string, n int) string {
	runes := []rune(str)
	length := len(runes)

	if n <= 0 {
		return ""
	}

	if n >= length {
		return str
	}

	return string(runes[length-n:])
}

// GetFirstNChars 函数用于处理GetFirstNChars相关逻辑。
func GetFirstNChars(str string, n int) string {
	runes := []rune(str)
	length := len(runes)

	if n <= 0 {
		return ""
	}

	if n >= length {
		return str
	}

	return string(runes[:n])
}

// ConvertStringToUint32 函数用于处理ConvertStringToUint32相关逻辑。
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

// ConvertStringToUint32Simple 函数用于处理ConvertStringToUint32Simple相关逻辑。
func ConvertStringToUint32Simple(str string) uint32 {
	result, err := ConvertStringToUint32(str)
	if err != nil {
		return 0
	}
	return result
}

// GetGenderFormIDCard 函数用于处理GetGenderFormIDCard相关逻辑。
func GetGenderFormIDCard(idcard string) string {
	if !ValidateChineseIDCard(idcard) {
		return "未知"
	}
	// 获取第17位数字(索引为16)
	// 中国身份证第17位数字表示性别：奇数为男性，偶数为女性
	genderDigit := idcard[16:17]

	// 将字符转换为数字
	digit, err := strconv.Atoi(genderDigit)
	if err != nil {
		return "无效身份证号"
	}

	// 判断奇偶性
	if digit%2 == 0 {
		return "女"
	}
	return "男"
}

// KeywordAssembly 函数用于处理KeywordAssembly相关逻辑。
func KeywordAssembly(keyword string) string {
	return fmt.Sprintf("%%%s%%", keyword)
}
