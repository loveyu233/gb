package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"testing"
)

func TestR(t *testing.T) {
	gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}))
	manager := gb.NewCaptchaManager(gb.RedisClient)
	generate, err := manager.Generate(gb.TypeSlide)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf("%+v", generate)
}
