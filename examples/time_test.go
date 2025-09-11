package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestTime(t *testing.T) {
	t.Log(gb.FormatDateRelativeDate(gb.Now().AddDate(0, 0, 0)))  // 今天
	t.Log(gb.FormatDateRelativeDate(gb.Now().AddDate(0, 0, -1))) // 昨天
	t.Log(gb.FormatDateRelativeDate(gb.Now().AddDate(0, -1, 0))) // 上月
	t.Log(gb.FormatDateRelativeDate(gb.Now().AddDate(0, -2, 2))) // 空字符串

	t.Log(gb.FormatTimeRelativeDate(gb.Now()))
}

func TestTimeTo(t *testing.T) {
	t.Log(gb.StringToGbDateTime("2025-01-11 09:30:20"))
	t.Log(gb.StringToGBDateOnly("2025-01-11"))
	t.Log(gb.StringToGBTimeOnly("09:30:20"))
	t.Log(gb.StringToGBTimeOnlyNoSec("09:30:20"))

	t.Log(gb.TimeToGBDateTime(gb.Now()))
	t.Log(gb.TimeToGBDateOnly(gb.Now()))
	t.Log(gb.TimeToGBTimeOnly(gb.Now()))
	t.Log(gb.TimeToGBTimeOnlyNoSec(gb.Now()))
}

func TestAAAA(t *testing.T) {
	t.Log(gb.TimeChineseWeekday(gb.Now()))
}
