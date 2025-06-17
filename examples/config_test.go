package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"testing"
)

type Config struct {
	Username string
	Password string
	User     struct {
		Age     int
		Address string
	}
}

func TestConfig(t *testing.T) {
	var cfg = new(Config)
	t.Log(gb.InitConfig("./config/config.json", cfg))
	fmt.Printf("%#v\n", cfg)

	t.Log(gb.InitConfig("./config/config.yml", cfg))
	fmt.Printf("%#v\n", cfg)

	t.Log(gb.InitConfig("./config/config.yaml", cfg))
	fmt.Printf("%#v\n", cfg)
}
