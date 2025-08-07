package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestD1(t *testing.T) {
	yuan := gb.Int64FenToDecimalYuan(99999)
	t.Log(yuan)
	s := gb.Int64FenToDecimalYuanString(99999)
	t.Log(s)
	fen := gb.DecimalYuanToInt64Fen(yuan)
	t.Log(fen)
}
