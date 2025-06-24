package gb

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"sync"
)

var RedisClient = new(Redis)

type Redis struct {
	redis.UniversalClient
	lock *redsync.Redsync
	once sync.Once
}

type WithRedisOption func(*redis.UniversalOptions)

func WithRedisAddressOption(address []string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Addrs = address
	}
}

func WithRedisDBOption(dn int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DB = dn
	}
}

func WithRedisUsernameOption(username string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Username = username
	}
}

func WithRedisPasswordOption(password string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Password = password
	}
}

func InitRedis(ops ...WithRedisOption) error {
	RedisClient.once = sync.Once{}
	opts := &redis.UniversalOptions{
		RouteRandomly: true,
	}
	for _, op := range ops {
		op(opts)
	}
	RedisClient.UniversalClient = redis.NewUniversalClient(opts)
	return RedisClient.UniversalClient.Ping(context.Background()).Err()
}

// NewLock 使用: https://github.com/go-redsync/redsync, 不设置options会有默认的重试次数和时间,也就是会lock会有错误返回,示例:examples/redis_test.go
func (r *Redis) NewLock(key string, options ...redsync.Option) *redsync.Mutex {
	r.once.Do(func() {
		r.lock = redsync.New(goredis.NewPool(RedisClient))
	})

	return r.lock.NewMutex(key, options...)
}
