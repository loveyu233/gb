package gb

import "github.com/shopspring/decimal"

// Float64ToDecimal 将 float64 转换为 decimal.Decimal
func Float64ToDecimal(value float64) decimal.Decimal {
	return decimal.NewFromFloat(value)
}

// DecimalToFloat64 将 decimal.Decimal 转换为 float64
func DecimalToFloat64(value decimal.Decimal) float64 {
	return value.InexactFloat64()
}

// DecimalYuanToFen 将decimal类型的金额（元）转换为int64类型的分
func DecimalYuanToFen(value decimal.Decimal) int64 {
	return value.Mul(decimal.NewFromInt(100)).IntPart()
}

// FenToDecimalYuan 将int类型的金额（分）转换为decimal.Decimal类型的金额（元）
func FenToDecimalYuan(value int64) decimal.Decimal {
	return decimal.NewFromInt(value).Div(decimal.NewFromInt(100))
}

// DecimalPercent 计算百分比
func DecimalPercent(value decimal.Decimal, percent decimal.Decimal) decimal.Decimal {
	return value.Mul(percent).Div(decimal.NewFromInt(100))
}
