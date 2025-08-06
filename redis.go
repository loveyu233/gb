package gb

import (
	"context"
	"crypto/tls"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"net"
	"sort"
	"sync"
	"time"
)

var InsRedis *RedisConfig

type RedisConfig struct {
	redis.UniversalClient
	lock *redsync.Redsync
	once sync.Once
}

type WithRedisOption func(*redis.UniversalOptions)

// WithRedisAddressOption 设置Redis服务器地址列表
func WithRedisAddressOption(address []string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Addrs = address
	}
}

// WithRedisClientNameOption 设置客户端名称
func WithRedisClientNameOption(clientName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ClientName = clientName
	}
}

// WithRedisDBOption 设置Redis数据库索引
func WithRedisDBOption(db int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DB = db
	}
}

// WithRedisDialerOption 设置自定义拨号器
func WithRedisDialerOption(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Dialer = dialer
	}
}

// WithRedisOnConnectOption 设置连接建立时的回调函数
func WithRedisOnConnectOption(onConnect func(ctx context.Context, cn *redis.Conn) error) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.OnConnect = onConnect
	}
}

// WithRedisProtocolOption 设置Redis协议版本
func WithRedisProtocolOption(protocol int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Protocol = protocol
	}
}

// WithRedisUsernameOption 设置Redis用户名
func WithRedisUsernameOption(username string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Username = username
	}
}

// WithRedisPasswordOption 设置Redis密码
func WithRedisPasswordOption(password string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Password = password
	}
}

// WithRedisSentinelUsernameOption 设置哨兵用户名
func WithRedisSentinelUsernameOption(sentinelUsername string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelUsername = sentinelUsername
	}
}

// WithRedisSentinelPasswordOption 设置哨兵密码
func WithRedisSentinelPasswordOption(sentinelPassword string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelPassword = sentinelPassword
	}
}

// WithRedisMaxRetriesOption 设置最大重试次数
func WithRedisMaxRetriesOption(maxRetries int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetries = maxRetries
	}
}

// WithRedisMinRetryBackoffOption 设置最小重试退避时间
func WithRedisMinRetryBackoffOption(minRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinRetryBackoff = minRetryBackoff
	}
}

// WithRedisMaxRetryBackoffOption 设置最大重试退避时间
func WithRedisMaxRetryBackoffOption(maxRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetryBackoff = maxRetryBackoff
	}
}

// WithRedisDialTimeoutOption 设置连接超时时间
func WithRedisDialTimeoutOption(dialTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DialTimeout = dialTimeout
	}
}

// WithRedisReadTimeoutOption 设置读取超时时间
func WithRedisReadTimeoutOption(readTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadTimeout = readTimeout
	}
}

// WithRedisWriteTimeoutOption 设置写入超时时间
func WithRedisWriteTimeoutOption(writeTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.WriteTimeout = writeTimeout
	}
}

// WithRedisContextTimeoutEnabledOption 设置是否启用上下文超时
func WithRedisContextTimeoutEnabledOption(enabled bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ContextTimeoutEnabled = enabled
	}
}

// WithRedisPoolFIFOOption 设置连接池是否使用FIFO模式
func WithRedisPoolFIFOOption(poolFIFO bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolFIFO = poolFIFO
	}
}

// WithRedisPoolSizeOption 设置连接池大小
func WithRedisPoolSizeOption(poolSize int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolSize = poolSize
	}
}

// WithRedisPoolTimeoutOption 设置从连接池获取连接的超时时间
func WithRedisPoolTimeoutOption(poolTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolTimeout = poolTimeout
	}
}

// WithRedisMinIdleConnsOption 设置连接池中最小空闲连接数
func WithRedisMinIdleConnsOption(minIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinIdleConns = minIdleConns
	}
}

// WithRedisMaxIdleConnsOption 设置连接池中最大空闲连接数
func WithRedisMaxIdleConnsOption(maxIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxIdleConns = maxIdleConns
	}
}

// WithRedisMaxActiveConnsOption 设置连接池中最大活跃连接数
func WithRedisMaxActiveConnsOption(maxActiveConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxActiveConns = maxActiveConns
	}
}

// WithRedisConnMaxIdleTimeOption 设置连接最大空闲时间
func WithRedisConnMaxIdleTimeOption(connMaxIdleTime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxIdleTime = connMaxIdleTime
	}
}

// WithRedisConnMaxLifetimeOption 设置连接最大生命周期
func WithRedisConnMaxLifetimeOption(connMaxLifetime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxLifetime = connMaxLifetime
	}
}

// WithRedisTLSConfigOption 设置TLS配置
func WithRedisTLSConfigOption(tlsConfig *tls.Config) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.TLSConfig = tlsConfig
	}
}

// WithRedisMaxRedirectsOption 设置集群模式下最大重定向次数
func WithRedisMaxRedirectsOption(maxRedirects int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRedirects = maxRedirects
	}
}

// WithRedisReadOnlyOption 设置是否只读模式
func WithRedisReadOnlyOption(readOnly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadOnly = readOnly
	}
}

// WithRedisRouteByLatencyOption 设置是否按延迟路由请求
func WithRedisRouteByLatencyOption(routeByLatency bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteByLatency = routeByLatency
	}
}

// WithRedisRouteRandomlyOption 设置是否随机路由请求
func WithRedisRouteRandomlyOption(routeRandomly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteRandomly = routeRandomly
	}
}

// WithRedisMasterNameOption 设置哨兵模式下的主服务器名称
func WithRedisMasterNameOption(masterName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MasterName = masterName
	}
}

// WithRedisDisableIdentityOption 设置是否禁用身份验证
func WithRedisDisableIdentityOption(disableIdentity bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DisableIdentity = disableIdentity
	}
}

// WithRedisIdentitySuffixOption 设置身份后缀
func WithRedisIdentitySuffixOption(identitySuffix string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IdentitySuffix = identitySuffix
	}
}

// WithRedisUnstableResp3Option 设置是否启用不稳定的RESP3协议
func WithRedisUnstableResp3Option(unstableResp3 bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.UnstableResp3 = unstableResp3
	}
}

// WithRedisIsClusterModeOption 设置是否为集群模式
func WithRedisIsClusterModeOption(isClusterMode bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IsClusterMode = isClusterMode
	}
}

func InitRedis(ops ...WithRedisOption) error {
	InsRedis = new(RedisConfig)
	InsRedis.once = sync.Once{}
	opts := &redis.UniversalOptions{}
	for _, op := range ops {
		op(opts)
	}
	if len(opts.Addrs) == 0 {
		panic("redis address is empty")
	}
	InsRedis.UniversalClient = redis.NewUniversalClient(opts)
	return InsRedis.UniversalClient.Ping(context.Background()).Err()
}

// NewLock 使用: https://github.com/go-redsync/redsync, 不设置options会有默认的重试次数和时间,也就是会lock会有错误返回,示例:examples/redis_test.go
func (r *RedisConfig) NewLock(key string, options ...redsync.Option) *redsync.Mutex {
	r.once.Do(func() {
		r.lock = redsync.New(goredis.NewPool(InsRedis))
	})

	return r.lock.NewMutex(key, options...)
}

// FindAllBitMapByTargetValue 返回bitmap类型key的value中bit位值为targetValue的bit位置
func (r *RedisConfig) FindAllBitMapByTargetValue(key string, targetValue byte) ([]int64, error) {
	value, err := r.Get(Context(), key).Result()
	if err != nil {
		return nil, err
	}

	var setBits []int64
	for byteIndex, b := range []byte(value) {
		for bitIndex := 0; bitIndex < 8; bitIndex++ {
			if (b>>bitIndex)&1 == targetValue {
				bitPosition := int64(byteIndex*8 + (7 - bitIndex))
				setBits = append(setBits, bitPosition)
			}
		}
	}

	sort.Slice(setBits, func(i, j int) bool {
		return setBits[i] < setBits[j]
	})
	return setBits, nil
}
func (r *RedisConfig) SetCaptcha(key string, value any, expiration time.Duration) error {
	return r.SetNX(context.Background(), key, value, expiration).Err()
}

func (r *RedisConfig) GetCaptcha(key string) (string, error) {
	return r.Get(context.Background(), key).Result()
}

func (r *RedisConfig) DelCaptcha(key string) error {
	return r.Del(context.Background(), key).Err()
}
