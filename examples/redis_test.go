package examples

import (
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		t.Fatal(err)
		return
	}

	for i := 0; i < 10; i++ {
		go func(index int) {
			lock := gb.RedisClient.NewLock("123")
			for {
				if err := lock.Lock(); err != nil {
					t.Log("加锁失败", index, err)
				} else {
					break
				}
			}
			defer lock.Unlock()

			time.Sleep(1 * time.Second)
			t.Log("加锁成功", index, time.Now().String())
		}(i)
	}
	select {}
}
