package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestA1(t *testing.T) {
	t.Log(gb.PasswordValidateStrength("123aAaa.", 6))
}
