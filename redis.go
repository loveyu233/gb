package gb

import (
	"context"
	"crypto/tls"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

var InsRedis *RedisConfig

type RedisConfig struct {
	redis.UniversalClient
	lock *redsync.Redsync
	once sync.Once
}

type WithRedisOption func(*redis.UniversalOptions)

// WithRedisAddressOption 用来设置 Redis 的节点地址。
func WithRedisAddressOption(address []string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Addrs = address
	}
}

// WithRedisClientNameOption 用来指定客户端名称。
func WithRedisClientNameOption(clientName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ClientName = clientName
	}
}

// WithRedisDBOption 用来设置默认数据库序号。
func WithRedisDBOption(db int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DB = db
	}
}

// WithRedisDialerOption 用来自定义底层拨号逻辑。
func WithRedisDialerOption(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Dialer = dialer
	}
}

// WithRedisOnConnectOption 用来注册连接建立后的回调。
func WithRedisOnConnectOption(onConnect func(ctx context.Context, cn *redis.Conn) error) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.OnConnect = onConnect
	}
}

// WithRedisProtocolOption 用来指定使用的协议版本。
func WithRedisProtocolOption(protocol int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Protocol = protocol
	}
}

// WithRedisUsernameOption 用来设置访问用户名。
func WithRedisUsernameOption(username string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Username = username
	}
}

// WithRedisPasswordOption 用来设置访问密码。
func WithRedisPasswordOption(password string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Password = password
	}
}

// WithRedisSentinelUsernameOption 用来设置哨兵用户名。
func WithRedisSentinelUsernameOption(sentinelUsername string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelUsername = sentinelUsername
	}
}

// WithRedisSentinelPasswordOption 用来设置哨兵密码。
func WithRedisSentinelPasswordOption(sentinelPassword string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelPassword = sentinelPassword
	}
}

// WithRedisMaxRetriesOption 用来配置最大重试次数。
func WithRedisMaxRetriesOption(maxRetries int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetries = maxRetries
	}
}

// WithRedisMinRetryBackoffOption 用来设置最小退避时长。
func WithRedisMinRetryBackoffOption(minRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinRetryBackoff = minRetryBackoff
	}
}

// WithRedisMaxRetryBackoffOption 用来设置最大退避时长。
func WithRedisMaxRetryBackoffOption(maxRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetryBackoff = maxRetryBackoff
	}
}

// WithRedisDialTimeoutOption 用来设置连接建立超时时间。
func WithRedisDialTimeoutOption(dialTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DialTimeout = dialTimeout
	}
}

// WithRedisReadTimeoutOption 用来设置读操作超时时间。
func WithRedisReadTimeoutOption(readTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadTimeout = readTimeout
	}
}

// WithRedisWriteTimeoutOption 用来设置写操作超时时间。
func WithRedisWriteTimeoutOption(writeTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.WriteTimeout = writeTimeout
	}
}

// WithRedisContextTimeoutEnabledOption 用来启用 Context 级超时控制。
func WithRedisContextTimeoutEnabledOption(enabled bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ContextTimeoutEnabled = enabled
	}
}

// WithRedisPoolFIFOOption 用来指定连接池使用 FIFO 策略。
func WithRedisPoolFIFOOption(poolFIFO bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolFIFO = poolFIFO
	}
}

// WithRedisPoolSizeOption 用来配置连接池大小。
func WithRedisPoolSizeOption(poolSize int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolSize = poolSize
	}
}

// WithRedisPoolTimeoutOption 用来设置等待连接的超时时间。
func WithRedisPoolTimeoutOption(poolTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolTimeout = poolTimeout
	}
}

// WithRedisMinIdleConnsOption 用来设置最小空闲连接数量。
func WithRedisMinIdleConnsOption(minIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinIdleConns = minIdleConns
	}
}

// WithRedisMaxIdleConnsOption 用来限制最大空闲连接数量。
func WithRedisMaxIdleConnsOption(maxIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxIdleConns = maxIdleConns
	}
}

// WithRedisMaxActiveConnsOption 用来设置最大活跃连接数。
func WithRedisMaxActiveConnsOption(maxActiveConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxActiveConns = maxActiveConns
	}
}

// WithRedisConnMaxIdleTimeOption 用来限定连接最长空闲时间。
func WithRedisConnMaxIdleTimeOption(connMaxIdleTime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxIdleTime = connMaxIdleTime
	}
}

// WithRedisConnMaxLifetimeOption 用来限定连接最大生命周期。
func WithRedisConnMaxLifetimeOption(connMaxLifetime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxLifetime = connMaxLifetime
	}
}

// WithRedisTLSConfigOption 用来配置 TLS 相关参数。
func WithRedisTLSConfigOption(tlsConfig *tls.Config) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.TLSConfig = tlsConfig
	}
}

// WithRedisMaxRedirectsOption 用来限定集群重定向次数。
func WithRedisMaxRedirectsOption(maxRedirects int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRedirects = maxRedirects
	}
}

// WithRedisReadOnlyOption 用来在集群模式下启用只读节点。
func WithRedisReadOnlyOption(readOnly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadOnly = readOnly
	}
}

// WithRedisRouteByLatencyOption 用来按延迟路由请求。
func WithRedisRouteByLatencyOption(routeByLatency bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteByLatency = routeByLatency
	}
}

// WithRedisRouteRandomlyOption 用来在多个节点间随机路由。
func WithRedisRouteRandomlyOption(routeRandomly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteRandomly = routeRandomly
	}
}

// WithRedisMasterNameOption 用来指定哨兵主节点名称。
func WithRedisMasterNameOption(masterName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MasterName = masterName
	}
}

// WithRedisDisableIdentityOption 用来禁用连接标识附加。
func WithRedisDisableIdentityOption(disableIdentity bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DisableIdentity = disableIdentity
	}
}

// WithRedisIdentitySuffixOption 用来设置连接标识后缀。
func WithRedisIdentitySuffixOption(identitySuffix string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IdentitySuffix = identitySuffix
	}
}

// WithRedisUnstableResp3Option 用来开启 resp3 实验性支持。
func WithRedisUnstableResp3Option(unstableResp3 bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.UnstableResp3 = unstableResp3
	}
}

// WithRedisIsClusterModeOption 用来声明是否以集群模式连接。
func WithRedisIsClusterModeOption(isClusterMode bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IsClusterMode = isClusterMode
	}
}

// InitRedis 用来依据提供的选项初始化全局 Redis 客户端。
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

// NewLock 用来基于 redsync 创建分布式锁。
func (r *RedisConfig) NewLock(key string, options ...redsync.Option) *redsync.Mutex {
	r.once.Do(func() {
		r.lock = redsync.New(goredis.NewPool(InsRedis))
	})

	return r.lock.NewMutex(key, options...)
}

// FindAllBitMapByTargetValue 用来返回位图中匹配目标值的位索引。
func (r *RedisConfig) FindAllBitMapByTargetValue(key string, targetValue byte) ([]int64, error) {
	ctx, cancel := Context()
	defer cancel()
	value, err := r.Get(ctx, key).Result()
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

// SetCaptcha 用来在 Redis 中缓存验证码。
func (r *RedisConfig) SetCaptcha(key string, value any, expiration time.Duration) error {
	return r.SetNX(context.Background(), key, value, expiration).Err()
}

// GetCaptcha 用来读取缓存的验证码。
func (r *RedisConfig) GetCaptcha(key string) (string, error) {
	return r.Get(context.Background(), key).Result()
}

// DelCaptcha 用来删除验证码缓存。
func (r *RedisConfig) DelCaptcha(key string) error {
	return r.Del(context.Background(), key).Err()
}
