package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestStr(t *testing.T) {
	s := "你好HelloWorld你好!"
	t.Log(gb.GetLastNChars(s, 3))
	t.Log(gb.GetFirstNChars(s, 3))
}
