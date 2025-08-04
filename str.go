package gb

// GetLastNChars 获取字符串的后n位字符
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

// GetFirstNChars 获取字符串的前n位字符
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
