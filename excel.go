package gb

import (
	"fmt"
	"unsafe"
)

// ExcelGetPosition 根据行列索引计算Excel位置
// row: 行索引（从0开始）
// col: 列索引（从0开始）
// 返回: Excel位置字符串，如 "A1", "B1", "AA10" 等
func ExcelGetPosition(row, col int64) string {
	if row < 0 || col < 0 {
		return ""
	}

	// 预分配足够大的缓冲区，避免重复分配
	// 最大列数约为18278 (ZZZ)，最大需要3个字符
	// 行数可能很大，预留20个字符应该足够
	buf := make([]byte, 0, 24)

	// 高效计算列字母
	buf = appendExcelColumn(buf, col)

	// 高效转换行号
	buf = appendInt64(buf, row+1)

	// 直接转换为字符串，避免额外拷贝
	return *(*string)(unsafe.Pointer(&buf))
}

// appendExcelColumn 高效地将列索引转换为Excel列字母并追加到缓冲区
func appendExcelColumn(buf []byte, col int64) []byte {
	if col < 26 {
		// 单字母情况，直接处理
		return append(buf, byte('A'+col))
	}

	// 多字母情况，使用栈来避免字符串反转
	var stack [4]byte // 最多4个字符 (AAAA对应列索引约为450,000)
	stackSize := 0

	for {
		remainder := col % 26
		stack[stackSize] = byte('A' + remainder)
		stackSize++

		col = col/26 - 1
		if col < 0 {
			break
		}
	}

	// 从栈顶开始追加到缓冲区
	for i := stackSize - 1; i >= 0; i-- {
		buf = append(buf, stack[i])
	}

	return buf
}

// appendInt64 高效地将int64追加到字节缓冲区
// 使用自定义实现避免strconv.Itoa的额外分配
func appendInt64(buf []byte, n int64) []byte {
	if n == 0 {
		return append(buf, '0')
	}

	// 处理负数（虽然在这个场景不会出现）
	if n < 0 {
		buf = append(buf, '-')
		n = -n
	}

	// 计算数字位数
	temp := n
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}

	// 预分配空间
	start := len(buf)
	buf = buf[:start+digits]

	// 从后往前填充数字
	for i := start + digits - 1; i >= start; i-- {
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	return buf
}

// ExcelGetPositionBatch 批量转换，适用于大量转换场景
// 复用缓冲区，进一步提高性能
func ExcelGetPositionBatch(positions []struct{ Row, Col int64 }) []string {
	results := make([]string, len(positions))
	buf := make([]byte, 0, 32) // 复用缓冲区

	for i, pos := range positions {
		buf = buf[:0] // 重置缓冲区长度，但保留容量

		if pos.Row < 0 || pos.Col < 0 {
			results[i] = ""
			continue
		}

		buf = appendExcelColumn(buf, pos.Col)
		buf = appendInt64(buf, pos.Row+1)

		// 创建字符串副本
		results[i] = string(buf)
	}

	return results
}

// ExcelColumnToIndex 将Excel列字母转换为索引（反向转换）
func ExcelColumnToIndex(col string) int64 {
	var result int64
	for _, char := range col {
		if char < 'A' || char > 'Z' {
			return -1 // 无效字符
		}
		result = result*26 + int64(char-'A'+1)
	}
	return result - 1
}

// ExcelParsePosition 将Excel位置字符串解析为行列索引
// 例如: "A1" -> (0, 0), "B10" -> (9, 1), "AA1" -> (0, 26)
// 返回: row, col, error
func ExcelParsePosition(position string) (row, col int64, err error) {
	if len(position) == 0 {
		return 0, 0, fmt.Errorf("空位置字符串")
	}

	// 分离字母部分和数字部分
	var colPart []byte
	var rowPart []byte

	i := 0
	// 提取列字母部分
	for i < len(position) {
		char := position[i]
		if char >= 'A' && char <= 'Z' {
			colPart = append(colPart, char)
			i++
		} else if char >= 'a' && char <= 'z' {
			// 支持小写字母，转换为大写
			colPart = append(colPart, char-'a'+'A')
			i++
		} else {
			break
		}
	}

	// 提取行数字部分
	for i < len(position) {
		char := position[i]
		if char >= '0' && char <= '9' {
			rowPart = append(rowPart, char)
			i++
		} else {
			return 0, 0, fmt.Errorf("行部分中的字符无效: %c", char)
		}
	}

	if len(colPart) == 0 {
		return 0, 0, fmt.Errorf("缺少列数据")
	}
	if len(rowPart) == 0 {
		return 0, 0, fmt.Errorf("缺少行数据")
	}

	// 解析列索引
	col = 0
	for _, char := range colPart {
		col = col*26 + int64(char-'A'+1)
	}
	col-- // 转换为0基索引

	// 解析行索引
	row = 0
	for _, char := range rowPart {
		digit := int64(char - '0')
		if row > (1<<63-1-digit)/10 { // 防止溢出
			return 0, 0, fmt.Errorf("行数太大")
		}
		row = row*10 + digit
	}
	row-- // 转换为0基索引

	if row < 0 || col < 0 {
		return 0, 0, fmt.Errorf("无效位置: 行=%d, 列=%d", row+1, col+1)
	}

	return row, col, nil
}

// ExcelParsePositionUnsafe 高性能版本，不进行错误检查
// 仅在确保输入格式正确时使用
func ExcelParsePositionUnsafe(position string) (row, col int64) {
	if len(position) == 0 {
		return 0, 0
	}

	i := 0
	col = 0

	// 解析列部分
	for i < len(position) {
		char := position[i]
		if char >= 'A' && char <= 'Z' {
			col = col*26 + int64(char-'A'+1)
			i++
		} else if char >= 'a' && char <= 'z' {
			col = col*26 + int64(char-'a'+1)
			i++
		} else {
			break
		}
	}
	col-- // 转换为0基索引

	// 解析行部分
	row = 0
	for i < len(position) {
		char := position[i]
		if char >= '0' && char <= '9' {
			row = row*10 + int64(char-'0')
			i++
		} else {
			break
		}
	}
	row-- // 转换为0基索引

	return row, col
}
