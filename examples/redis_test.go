package examples

import (
	"fmt"
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
			t.Log("加锁成功", index, gb.Now().String())
		}(i)
	}
	wg.Wait()
}

func TestBit(t *testing.T) {
	err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		t.Fatal(err)
		return
	}
	value, err := gb.RedisClient.FindAllBitMapByTargetValue("bit1", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(value)
}

func TestLua(t *testing.T) {
	err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}), gb.WithRedisDBOption(0))
	if err != nil {
		panic(err)
	}
	//value, err := gb.LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreDesc("rank_list1", 0, 3, "100002")
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	//if value == nil {
	//	fmt.Println("空数据")
	//	return
	//}
	//fmt.Printf("%+v\n", value.Target)
	//for i := range value.Range {
	//	fmt.Printf("range item:%+v\n", value.Range[i])
	//}

	//value, err := gb.LuaRedisZSetGetTargetKeyAndStartToEndRankByScoreAndGetHashValue("rank_list", "user_info", 0, 3, "100002", true)
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	//if value == nil {
	//	fmt.Println("空数据")
	//	return
	//}
	//fmt.Printf("%+v\n", value.Target)
	//for i := range value.Range {
	//	fmt.Printf("range item:%+v\n", value.Range[i])
	//}

	//desc, err := gb.LuaRedisZSetGetMemberScoreAndRankAndGetHashValue("rank_list", "user_info", "100002", true)
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//fmt.Printf("%+v\n", desc)

	desc, err := gb.LuaRedisZSetGetMultipleMembersScoreAndRankAndHashValuesDesc("rank_list", "user_info", []string{"100003", "100004", "100005"})
	if err != nil {
		t.Log(err)
		return
	}
	for _, item := range desc {
		fmt.Printf("%+v\n", item)
	}
}
