package gb

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// DecimalYuanToInt64Fen 将decimal类型的金额（元）转换为int64类型的分
func DecimalYuanToInt64Fen(value decimal.Decimal) int64 {
	return value.Mul(decimal.NewFromInt(100)).IntPart()
}

// Int64FenToDecimalYuan 将int类型的金额（分）转换为decimal.Decimal类型的金额（元）
func Int64FenToDecimalYuan(value int64) decimal.Decimal {
	return decimal.NewFromInt(value).Div(decimal.NewFromInt(100))
}

// Int64FenToDecimalYuanString 将int类型的金额（分）转换为decimal.Decimal类型的金额（元）带单位元的字符串
func Int64FenToDecimalYuanString(value int64) string {
	return fmt.Sprintf("%v元", decimal.NewFromInt(value).Div(decimal.NewFromInt(100)))
}
