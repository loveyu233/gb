package examples

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

func TestCorn(t *testing.T) {
	gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}))

	err := gb.InitCornJob()
	if err != nil {
		t.Fatal(err)
		return
	}
	for i := 0; i < 10; i++ {
		go func() {
			_, err := gb.InsCornJob.RunJobEveryDurationTheOne(1, 2*time.Second, gocron.NewTask(func() {
				t.Log(gb.NowString())
			}))
			t.Log(err)
		}()
	}
	time.Sleep(time.Second)
	gb.InsCornJob.Start()
	select {}
}
