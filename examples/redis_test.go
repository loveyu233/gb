package examples

import (
	"github.com/loveyu233/gb"
	"sync"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		t.Fatal(err)
		return
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lock := gb.RedisClient.NewLock("123")
			for {
				if err := lock.Lock(); err != nil {
					t.Log("加锁失败", index, err)
				} else {
					break
				}
			}
			defer lock.Unlock()

			time.Sleep(2 * time.Second)
			t.Log("加锁成功", index, time.Now().String())
		}(i)
	}
	wg.Wait()
}
