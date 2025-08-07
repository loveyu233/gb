package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		//t.Log(gb.Random(10, gb.RandomCharacterSetUpperStrExcludeCharIO(), gb.RandomCharacterSetNumberStrExcludeCharo1()))
		t.Log(gb.RandomExcludeErrorPronCharacters(3))
	}
}
