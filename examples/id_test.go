package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestUUID(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Log(gb.GetUUID())
	}
}

func TestXID(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Log(gb.GetXID())
	}
}

func TestSno(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Log(gb.GetSnowflakeID())
	}
}
