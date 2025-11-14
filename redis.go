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

// WithRedisAddressOption 函数用于处理WithRedisAddressOption相关逻辑。
func WithRedisAddressOption(address []string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Addrs = address
	}
}

// WithRedisClientNameOption 函数用于处理WithRedisClientNameOption相关逻辑。
func WithRedisClientNameOption(clientName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ClientName = clientName
	}
}

// WithRedisDBOption 函数用于处理WithRedisDBOption相关逻辑。
func WithRedisDBOption(db int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DB = db
	}
}

// WithRedisDialerOption 函数用于处理WithRedisDialerOption相关逻辑。
func WithRedisDialerOption(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Dialer = dialer
	}
}

// WithRedisOnConnectOption 函数用于处理WithRedisOnConnectOption相关逻辑。
func WithRedisOnConnectOption(onConnect func(ctx context.Context, cn *redis.Conn) error) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.OnConnect = onConnect
	}
}

// WithRedisProtocolOption 函数用于处理WithRedisProtocolOption相关逻辑。
func WithRedisProtocolOption(protocol int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Protocol = protocol
	}
}

// WithRedisUsernameOption 函数用于处理WithRedisUsernameOption相关逻辑。
func WithRedisUsernameOption(username string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Username = username
	}
}

// WithRedisPasswordOption 函数用于处理WithRedisPasswordOption相关逻辑。
func WithRedisPasswordOption(password string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.Password = password
	}
}

// WithRedisSentinelUsernameOption 函数用于处理WithRedisSentinelUsernameOption相关逻辑。
func WithRedisSentinelUsernameOption(sentinelUsername string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelUsername = sentinelUsername
	}
}

// WithRedisSentinelPasswordOption 函数用于处理WithRedisSentinelPasswordOption相关逻辑。
func WithRedisSentinelPasswordOption(sentinelPassword string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.SentinelPassword = sentinelPassword
	}
}

// WithRedisMaxRetriesOption 函数用于处理WithRedisMaxRetriesOption相关逻辑。
func WithRedisMaxRetriesOption(maxRetries int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetries = maxRetries
	}
}

// WithRedisMinRetryBackoffOption 函数用于处理WithRedisMinRetryBackoffOption相关逻辑。
func WithRedisMinRetryBackoffOption(minRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinRetryBackoff = minRetryBackoff
	}
}

// WithRedisMaxRetryBackoffOption 函数用于处理WithRedisMaxRetryBackoffOption相关逻辑。
func WithRedisMaxRetryBackoffOption(maxRetryBackoff time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRetryBackoff = maxRetryBackoff
	}
}

// WithRedisDialTimeoutOption 函数用于处理WithRedisDialTimeoutOption相关逻辑。
func WithRedisDialTimeoutOption(dialTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DialTimeout = dialTimeout
	}
}

// WithRedisReadTimeoutOption 函数用于处理WithRedisReadTimeoutOption相关逻辑。
func WithRedisReadTimeoutOption(readTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadTimeout = readTimeout
	}
}

// WithRedisWriteTimeoutOption 函数用于处理WithRedisWriteTimeoutOption相关逻辑。
func WithRedisWriteTimeoutOption(writeTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.WriteTimeout = writeTimeout
	}
}

// WithRedisContextTimeoutEnabledOption 函数用于处理WithRedisContextTimeoutEnabledOption相关逻辑。
func WithRedisContextTimeoutEnabledOption(enabled bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ContextTimeoutEnabled = enabled
	}
}

// WithRedisPoolFIFOOption 函数用于处理WithRedisPoolFIFOOption相关逻辑。
func WithRedisPoolFIFOOption(poolFIFO bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolFIFO = poolFIFO
	}
}

// WithRedisPoolSizeOption 函数用于处理WithRedisPoolSizeOption相关逻辑。
func WithRedisPoolSizeOption(poolSize int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolSize = poolSize
	}
}

// WithRedisPoolTimeoutOption 函数用于处理WithRedisPoolTimeoutOption相关逻辑。
func WithRedisPoolTimeoutOption(poolTimeout time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.PoolTimeout = poolTimeout
	}
}

// WithRedisMinIdleConnsOption 函数用于处理WithRedisMinIdleConnsOption相关逻辑。
func WithRedisMinIdleConnsOption(minIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MinIdleConns = minIdleConns
	}
}

// WithRedisMaxIdleConnsOption 函数用于处理WithRedisMaxIdleConnsOption相关逻辑。
func WithRedisMaxIdleConnsOption(maxIdleConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxIdleConns = maxIdleConns
	}
}

// WithRedisMaxActiveConnsOption 函数用于处理WithRedisMaxActiveConnsOption相关逻辑。
func WithRedisMaxActiveConnsOption(maxActiveConns int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxActiveConns = maxActiveConns
	}
}

// WithRedisConnMaxIdleTimeOption 函数用于处理WithRedisConnMaxIdleTimeOption相关逻辑。
func WithRedisConnMaxIdleTimeOption(connMaxIdleTime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxIdleTime = connMaxIdleTime
	}
}

// WithRedisConnMaxLifetimeOption 函数用于处理WithRedisConnMaxLifetimeOption相关逻辑。
func WithRedisConnMaxLifetimeOption(connMaxLifetime time.Duration) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ConnMaxLifetime = connMaxLifetime
	}
}

// WithRedisTLSConfigOption 函数用于处理WithRedisTLSConfigOption相关逻辑。
func WithRedisTLSConfigOption(tlsConfig *tls.Config) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.TLSConfig = tlsConfig
	}
}

// WithRedisMaxRedirectsOption 函数用于处理WithRedisMaxRedirectsOption相关逻辑。
func WithRedisMaxRedirectsOption(maxRedirects int) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MaxRedirects = maxRedirects
	}
}

// WithRedisReadOnlyOption 函数用于处理WithRedisReadOnlyOption相关逻辑。
func WithRedisReadOnlyOption(readOnly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.ReadOnly = readOnly
	}
}

// WithRedisRouteByLatencyOption 函数用于处理WithRedisRouteByLatencyOption相关逻辑。
func WithRedisRouteByLatencyOption(routeByLatency bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteByLatency = routeByLatency
	}
}

// WithRedisRouteRandomlyOption 函数用于处理WithRedisRouteRandomlyOption相关逻辑。
func WithRedisRouteRandomlyOption(routeRandomly bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.RouteRandomly = routeRandomly
	}
}

// WithRedisMasterNameOption 函数用于处理WithRedisMasterNameOption相关逻辑。
func WithRedisMasterNameOption(masterName string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.MasterName = masterName
	}
}

// WithRedisDisableIdentityOption 函数用于处理WithRedisDisableIdentityOption相关逻辑。
func WithRedisDisableIdentityOption(disableIdentity bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.DisableIdentity = disableIdentity
	}
}

// WithRedisIdentitySuffixOption 函数用于处理WithRedisIdentitySuffixOption相关逻辑。
func WithRedisIdentitySuffixOption(identitySuffix string) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IdentitySuffix = identitySuffix
	}
}

// WithRedisUnstableResp3Option 函数用于处理WithRedisUnstableResp3Option相关逻辑。
func WithRedisUnstableResp3Option(unstableResp3 bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.UnstableResp3 = unstableResp3
	}
}

// WithRedisIsClusterModeOption 函数用于处理WithRedisIsClusterModeOption相关逻辑。
func WithRedisIsClusterModeOption(isClusterMode bool) WithRedisOption {
	return func(options *redis.UniversalOptions) {
		options.IsClusterMode = isClusterMode
	}
}

// InitRedis 函数用于处理InitRedis相关逻辑。
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

// NewLock 方法用于处理NewLock相关逻辑。
func (r *RedisConfig) NewLock(key string, options ...redsync.Option) *redsync.Mutex {
	r.once.Do(func() {
		r.lock = redsync.New(goredis.NewPool(InsRedis))
	})

	return r.lock.NewMutex(key, options...)
}

// FindAllBitMapByTargetValue 方法用于处理FindAllBitMapByTargetValue相关逻辑。
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

// SetCaptcha 方法用于处理SetCaptcha相关逻辑。
func (r *RedisConfig) SetCaptcha(key string, value any, expiration time.Duration) error {
	return r.SetNX(context.Background(), key, value, expiration).Err()
}

// GetCaptcha 方法用于处理GetCaptcha相关逻辑。
func (r *RedisConfig) GetCaptcha(key string) (string, error) {
	return r.Get(context.Background(), key).Result()
}

// DelCaptcha 方法用于处理DelCaptcha相关逻辑。
func (r *RedisConfig) DelCaptcha(key string) error {
	return r.Del(context.Background(), key).Err()
}
