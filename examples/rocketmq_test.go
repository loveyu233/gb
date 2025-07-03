package examples

import (
	"github.com/loveyu233/gb"
	"time"
)

func initRock() {
	config := gb.InitRocketMQConfig("127.0.0.1:8080", "cg", "rocketmq2", "12345678", []string{"test-1", "topic-time"})
	err, _ := config.GetProduct()
	if err != nil {
		panic(err)
	}
	consumer, err := config.GetConsumer("test-1", "user", 40*time.Second)
	if err != nil {
		panic(err)
	}
	_ = consumer
}
