package gb

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// DecimalYuanToInt64Fen 函数用于处理DecimalYuanToInt64Fen相关逻辑。
func DecimalYuanToInt64Fen(value decimal.Decimal) int64 {
	return value.Mul(decimal.NewFromInt(100)).IntPart()
}

// Int64FenToDecimalYuan 函数用于处理Int64FenToDecimalYuan相关逻辑。
func Int64FenToDecimalYuan(value int64) decimal.Decimal {
	return decimal.NewFromInt(value).Div(decimal.NewFromInt(100))
}

// Int64FenToDecimalYuanString 函数用于处理Int64FenToDecimalYuanString相关逻辑。
func Int64FenToDecimalYuanString(value int64) string {
	return fmt.Sprintf("%v元", decimal.NewFromInt(value).Div(decimal.NewFromInt(100)))
}
