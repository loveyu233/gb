package examples

import (
	"fmt"
	"testing"

	"github.com/loveyu233/gb"
)

func TestR(t *testing.T) {
	gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}))
	gb.InitCaptchaManager(gb.InsRedis)
	generate, err := gb.InsCaptcha.Generate(gb.TypeSlide)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf("%+v", generate)
}
