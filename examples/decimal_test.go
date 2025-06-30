package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestD1(t *testing.T) {
	decimal := gb.Float64ToDecimal(1.234)
	t.Log(decimal.String())

	decimalToFloat64 := gb.DecimalToFloat64(decimal)
	t.Log(decimalToFloat64)

	toFen := gb.DecimalYuanToFen(decimal)
	t.Log(toFen)

	decimalYuan := gb.FenToDecimalYuan(10099)
	t.Log(decimalYuan)

	percent := gb.DecimalPercent(gb.Float64ToDecimal(100), gb.Float64ToDecimal(3))
	t.Log(percent)
}

func TestD2(t *testing.T) {
	t.Log(gb.DecimalADD(float64(0.1), float64(0.2), float64(0.3)))
	t.Log(gb.DecimalADD("0.1", "0.2", 0.3))
	t.Log(gb.DecimalSUB(float64(0.1), float64(0.2), float64(0.3)))
	t.Log(gb.DecimalSUB("0.1", "0.2", "0.3"))
}
