package examples

import (
	"github.com/loveyu233/gb"
	"sync"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	redis, err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		t.Fatal(err)
		return
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lock := redis.NewLock("123")
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

func TestBit(t *testing.T) {
	redis, err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		t.Fatal(err)
		return
	}
	value, err := redis.FindAllBitMapByTargetValue("bit1", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(value)
}
