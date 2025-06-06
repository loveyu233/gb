package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestRedis(t *testing.T) {
	err := gb.NewRedisClient("127.0.0.1:6379", gb.WithRedisClientDB(18))
	if err != nil {
		panic(err)
	}
	t.Log(gb.Redis.Get("qrcode/2"))
}
