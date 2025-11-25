package gb

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// DecimalYuanToInt64Fen 用来将金额从元转换成以分为单位的整数。
func DecimalYuanToInt64Fen(value decimal.Decimal) int64 {
	return value.Mul(decimal.NewFromInt(100)).IntPart()
}

// Int64FenToDecimalYuan 用来把以分表示的金额转换为 decimal 元。
func Int64FenToDecimalYuan(value int64) decimal.Decimal {
	return decimal.NewFromInt(value).Div(decimal.NewFromInt(100))
}

// Int64FenToDecimalYuanString 用来生成带“元”单位的金额字符串。
func Int64FenToDecimalYuanString(value int64) string {
	return fmt.Sprintf("%v元", decimal.NewFromInt(value).Div(decimal.NewFromInt(100)))
}
