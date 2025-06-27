package examples

import (
	"fmt"
	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/loveyu233/gb"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func initRock() {
	config := gb.InitRocketMQConfig("127.0.0.1:8080", "cg", "rocketmq2", "12345678", []string{"test-1", "topic-time"})
	err := config.GetProduct()
	if err != nil {
		panic(err)
	}
	err = config.GetConsumer("test-1", "user", 40*time.Second)
	if err != nil {
		panic(err)
	}
}

func TestProduct(t *testing.T) {
	msg, err := gb.RocketMQ.SendMsg(gb.Context(), golang.Message{
		Topic: "test-1",
		Body:  []byte("hello world2"),
		Tag:   nil,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	fmt.Printf("%+v\n", msg)
}

func TestConsumer(t *testing.T) {
	time.Sleep(time.Second * 1)
	msg, err := gb.RocketMQ.CaptureMsg(context.Background(), 16, 20*time.Second)
	if err != nil {
		t.Log(err.Error())
		return
	}
	for _, item := range msg {
		gb.RocketMQ.Consumer.Ack(context.TODO(), item)
		fmt.Printf("tga:%v key:%v body:%s \n", item.GetTag(), item.GetKeys(), string(item.GetBody()))
	}
}
