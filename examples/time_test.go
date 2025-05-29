package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestTime(t *testing.T) {
	t.Log(gb.FormatRelativeDate(gb.GetCurrentTime().AddDate(0, 0, -1)))    // 昨天
	t.Log(gb.FormatRelativeDate(gb.GetCurrentTime().AddDate(0, 0, 0)))     // 今天
	t.Log(gb.FormatRelativeDate(gb.GetCurrentTime().AddDate(0, 0, 1)))     // 明天
	t.Log(gb.FormatRelativeDate(gb.GetCurrentTime().AddDate(0, 0, 2)))     // 2025-05-31
	t.Log(gb.FormatRelativeDate(gb.GetCurrentTime().AddDate(0, 0, 2), "")) // 空字符串
}
