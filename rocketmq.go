package gb

import (
	"github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"golang.org/x/net/context"
	"time"
)

var (
	InsRocketMQ = new(RocketMQClient)
)

type RocketMQConfig struct {
	endpoint      string
	consumerGroup string
	accessKey     string
	accessSecret  string
	topics        []string
}

// InitRocketMQConfig 创建RocketMQ配置
func InitRocketMQConfig(endpoint, consumerGroup, accessKey, accessSecret string, topics []string) *RocketMQConfig {
	return &RocketMQConfig{
		endpoint:      endpoint,
		consumerGroup: consumerGroup,
		accessKey:     accessKey,
		accessSecret:  accessSecret,
		topics:        topics,
	}
}

type RocketMQClient struct {
	Product  golang.Producer
	Consumer golang.SimpleConsumer
}

// GetProduct 获取生产者
func (conf *RocketMQConfig) GetProduct() error {
	config := &golang.Config{
		Endpoint:      conf.endpoint,
		ConsumerGroup: conf.consumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    conf.accessKey,
			AccessSecret: conf.accessSecret,
		},
	}

	producer, err := golang.NewProducer(config, golang.WithTopics(conf.topics...))
	if err != nil {
		return err
	}
	err = producer.Start()
	if err != nil {
		return err
	}
	InsRocketMQ.Product = producer
	return nil
}

// GetConsumer 获取消费者
func (conf *RocketMQConfig) GetConsumer(topic string, tag string, withAwaitDuration time.Duration) error {
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		Endpoint:      conf.endpoint,
		ConsumerGroup: conf.consumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    conf.accessKey,
			AccessSecret: conf.accessSecret,
		},
	},
		golang.WithAwaitDuration(withAwaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			topic: golang.NewFilterExpression(tag),
		}),
	)
	if err != nil {
		return err
	}
	err = simpleConsumer.Start()
	if err != nil {
		return err
	}
	InsRocketMQ.Consumer = simpleConsumer
	return nil
}

// SendMsg 发送同步消息
func (r *RocketMQClient) SendMsg(ctx context.Context, msg golang.Message) ([]*golang.SendReceipt, error) {
	return r.Product.Send(ctx, &msg)
}

// SendSyncMsg 发送异步消息
func (r *RocketMQClient) SendSyncMsg(ctx context.Context, msg golang.Message, syncHandler func(context.Context, []*golang.SendReceipt, error)) {
	r.Product.SendAsync(ctx, &msg, syncHandler)
}

// CaptureMsg 获取消息,调用 CaptureMsg 前要time.Sleep(time.Second * 1),不停顿一下有很大可能会报错空指针,原因未知
func (r *RocketMQClient) CaptureMsg(ctx context.Context, maxMessageNum int32, invisibleDuration time.Duration) ([]*golang.MessageView, error) {
	receive, err := r.Consumer.Receive(ctx, maxMessageNum, invisibleDuration)
	if err != nil {
		return nil, err
	}
	return receive, nil
}
