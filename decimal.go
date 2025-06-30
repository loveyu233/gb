package gb

import (
	"fmt"
	"github.com/shopspring/decimal"
)

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

// DecimalADD values相加
func DecimalADD(values ...any) (decimal.Decimal, error) {
	if len(values) == 0 {
		return decimal.Zero, nil
	}

	result, err := convertToDecimal(values[0])
	if err != nil {
		return decimal.Zero, err
	}
	for i := 1; i < len(values); i++ {
		d, err := convertToDecimal(values[i])
		if err != nil {
			return decimal.Zero, nil
		}
		result = result.Add(d)
	}

	return result, nil
}

// DecimalSUB values相减
func DecimalSUB(values ...any) (decimal.Decimal, error) {
	if len(values) == 0 {
		return decimal.Zero, nil
	}

	result, err := convertToDecimal(values[0])
	if err != nil {
		return decimal.Zero, err
	}
	for i := 1; i < len(values); i++ {
		d, err := convertToDecimal(values[i])
		if err != nil {
			return decimal.Zero, nil
		}
		result = result.Sub(d)
	}

	return result, nil
}

// convertToDecimal 将任意类型转换为decimal.Decimal的辅助函数
func convertToDecimal(value any) (decimal.Decimal, error) {
	switch v := value.(type) {
	case string:
		d, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero, err
		}
		return d, nil
	case int:
		return decimal.NewFromInt(int64(v)), nil
	case int8:
		return decimal.NewFromInt(int64(v)), nil
	case int16:
		return decimal.NewFromInt(int64(v)), nil
	case int32:
		return decimal.NewFromInt(int64(v)), nil
	case int64:
		return decimal.NewFromInt(v), nil
	case uint:
		return decimal.NewFromInt(int64(v)), nil
	case uint8:
		return decimal.NewFromInt(int64(v)), nil
	case uint16:
		return decimal.NewFromInt(int64(v)), nil
	case uint32:
		return decimal.NewFromInt(int64(v)), nil
	case uint64:
		return decimal.NewFromInt(int64(v)), nil
	case float32:
		return decimal.NewFromFloat32(v), nil
	case float64:
		return decimal.NewFromFloat(v), nil
	case decimal.Decimal:
		return v, nil
	default:
		// 尝试转换为字符串再解析
		str := fmt.Sprintf("%v", v)
		d, err := decimal.NewFromString(str)
		if err != nil {
			return decimal.Zero, err
		}
		return d, nil
	}
}
