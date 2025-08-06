package examples

import (
	"github.com/loveyu233/gb"
	"golang.org/x/net/context"
	"sync"
	"testing"
)

func init() {
	err := gb.InitEtcd(gb.WithEtcdEndpointsOpt([]string{"127.0.0.1:2379"}))
	if err != nil {
		panic(err)
	}
}

func TestEtcd(t *testing.T) {
	put, err := gb.InsEtcd.Put(gb.Context(), "key", "value")
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
			lock, err := gb.InsEtcd.NewLock("lock-event-d09777dfa1ab37c5")
			if err != nil {
				panic(err)
			}
			// 设置等待时间为3s超时则报错
			if err = lock.Lock(context.Background()); err != nil {
				t.Log("加锁失败:", index, err)
				return
			}
			defer lock.Unlock(gb.Context())
			t.Log("success", index, gb.Now().String())
		}(i)
	}
	wg.Wait()
}
