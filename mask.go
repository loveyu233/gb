package gb

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateChineseMobile 校验中国手机号格式是否正确
func ValidateChineseMobile(mobile string) bool {
	// 去除空格和特殊字符
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")

	// 中国手机号正则表达式
	// 1开头，第二位是3-9，总共11位数字
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, mobile)
	return matched
}

// MaskMobileCustom 自定义脱敏格式的手机号处理
// prefixLen: 前面保留的位数
// suffixLen: 后面保留的位数
// maskChar: 用于遮蔽的字符
func MaskMobileCustom(mobile string, prefixLen, suffixLen int, maskChar rune) string {
	// 去除空格和特殊字符
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")

	// 如果不是有效的手机号格式，返回原字符串
	if !ValidateChineseMobile(mobile) {
		return mobile
	}

	// 验证参数有效性
	if prefixLen < 0 || suffixLen < 0 || prefixLen+suffixLen >= len(mobile) {
		return mobile
	}

	// 计算中间需要遮蔽的位数
	maskLen := len(mobile) - prefixLen - suffixLen

	// 构建遮蔽字符串
	maskStr := strings.Repeat(string(maskChar), maskLen)

	// 返回脱敏后的手机号
	return mobile[:prefixLen] + maskStr + mobile[len(mobile)-suffixLen:]
}

// ValidateChineseIDCard 校验中国身份证号格式是否正确
func ValidateChineseIDCard(idCard string) bool {
	// 去除空格
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ToUpper(idCard)

	// 检查长度，必须是18位
	if len(idCard) != 18 {
		return false
	}

	// 检查前17位是否都是数字
	for i := 0; i < 17; i++ {
		if idCard[i] < '0' || idCard[i] > '9' {
			return false
		}
	}

	// 检查最后一位（校验码）
	lastChar := idCard[17]
	if lastChar != 'X' && (lastChar < '0' || lastChar > '9') {
		return false
	}

	// 计算校验码
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	for i := 0; i < 17; i++ {
		digit, _ := strconv.Atoi(string(idCard[i]))
		sum += digit * weights[i]
	}

	expectedCheckCode := checkCodes[sum%11]
	return byte(lastChar) == expectedCheckCode
}

// MaskMobile 对手机号进行脱敏处理
// 显示前3位和后3位，中间5位用*号替换
func MaskMobile(mobile string) string {
	// 去除空格和特殊字符
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")

	// 如果不是有效的手机号格式，返回原字符串
	if !ValidateChineseMobile(mobile) {
		return mobile
	}

	// 对手机号进行脱敏：前3位 + 5个* + 后3位
	return mobile[:3] + "*****" + mobile[8:]
}

// MaskIDCardCustom 自定义脱敏格式的身份证号处理
// prefixLen: 前面保留的位数
// suffixLen: 后面保留的位数
// maskChar: 用于遮蔽的字符
func MaskIDCardCustom(idCard string, prefixLen, suffixLen int, maskChar rune) string {
	// 去除空格
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ToUpper(idCard)

	// 如果不是有效的身份证号格式，返回原字符串
	if !ValidateChineseIDCard(idCard) {
		return idCard
	}

	// 验证参数有效性
	if prefixLen < 0 || suffixLen < 0 || prefixLen+suffixLen >= len(idCard) {
		return idCard
	}

	// 计算中间需要遮蔽的位数
	maskLen := len(idCard) - prefixLen - suffixLen

	// 构建遮蔽字符串
	maskStr := strings.Repeat(string(maskChar), maskLen)

	// 返回脱敏后的身份证号
	return idCard[:prefixLen] + maskStr + idCard[len(idCard)-suffixLen:]
}

// MaskIDCardBirthday 专门隐藏生日信息的身份证脱敏
// 保留地区码（前6位）和校验码（最后1位），隐藏生日和顺序码
func MaskIDCardBirthday(idCard string) string {
	// 去除空格
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ToUpper(idCard)

	// 如果不是有效的身份证号格式，返回原字符串
	if !ValidateChineseIDCard(idCard) {
		return idCard
	}

	// 身份证结构：前6位地区码 + 8位生日 + 3位顺序码 + 1位校验码
	// 保留地区码和校验码，隐藏生日和顺序码
	return idCard[:6] + "***********" + idCard[17:]
}

// MaskIDCard 对身份证号进行脱敏处理
// 显示前6位（地区码）和后4位（出生年份），中间8位用*号替换
func MaskIDCard(idCard string) string {
	// 去除空格
	idCard = strings.ReplaceAll(idCard, " ", "")
	idCard = strings.ToUpper(idCard)

	// 如果不是有效的身份证号格式，返回原字符串
	if !ValidateChineseIDCard(idCard) {
		return idCard
	}

	// 对身份证号进行脱敏：前6位 + 8个* + 后4位
	return idCard[:6] + "********" + idCard[14:]
}

// ValidateCustomGUIDRegex 使用正则表达式校验自定义GUID格式
// 要求：长度必须为20位，只能包含数字和大小写字母
func ValidateCustomGUIDRegex(guid string) bool {
	// 去除空格
	guid = strings.ReplaceAll(guid, " ", "")

	// 正则表达式：^[a-zA-Z0-9]{20}$
	// ^ 表示开始，$ 表示结束，[a-zA-Z0-9] 表示数字和大小写字母，{20} 表示恰好20位
	pattern := `^[a-zA-Z0-9]{20}$`
	matched, _ := regexp.MatchString(pattern, guid)
	return matched
}

func MaskUsername(username string) string {
	return GetFirstNChars(username, 1) + "*"
}
