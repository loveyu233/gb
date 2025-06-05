package examples

import (
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	err := gb.NewRedisClient("127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	t.Log(gb.R.Del("a"))
	t.Log(gb.R.Get("a"))
	set, err := gb.R.Set("a", "b", 10*time.Second)
	if err != nil {
		panic(err)
	}
	t.Log(set)
	t.Log(gb.R.Get("a"))
}
