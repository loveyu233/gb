package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestGen(t *testing.T) {
	encryption, err := gb.PasswordEncryption("hello world")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(encryption)

	t.Log(gb.PasswordCompare(encryption, "hello world"))
}

func TestA1(t *testing.T) {
	t.Log(gb.PasswordValidateStrength("123aAaa.", 6, 12))
}
