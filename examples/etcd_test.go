package examples

import (
	"github.com/loveyu233/gb"
	"golang.org/x/net/context"
	"sync"
	"testing"
	"time"
)

func init() {
	err := gb.InitEtcd(gb.WithEtcdEndpointsOpt([]string{"127.0.0.1:2379"}))
	if err != nil {
		panic(err)
	}
}

func TestEtcd(t *testing.T) {
	put, err := gb.ETCDClient.Put(gb.Context(), "key", "value")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(put)
}

func TestLock(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lock, err := gb.ETCDClient.NewLock("123")
			if err != nil {
				panic(err)
			}
			// 设置等待时间为3s超时则报错
			if err = lock.Lock(gb.Context(3)); err != nil {
				t.Log("加锁失败:", index)
				return
			}
			defer lock.Unlock(context.Background())
			time.Sleep(2 * time.Second)
			t.Log("success", index, time.Now().String())
		}(i)
	}
	wg.Wait()
}
