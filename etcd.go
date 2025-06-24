package gb

import (
	"crypto/tls"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var ETCDClient = new(Etcd)

type Etcd struct {
	*clientv3.Client
}

type WithEtcdOpt func(*clientv3.Config)

// WithEtcdEndpointsOpt 设置etcd集群端点
func WithEtcdEndpointsOpt(endpoints []string) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.Endpoints = endpoints
	}
}

// WithEtcdAutoSyncIntervalOpt 设置自动同步间隔
func WithEtcdAutoSyncIntervalOpt(interval time.Duration) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.AutoSyncInterval = interval
	}
}

// WithEtcdDialTimeoutOpt 设置连接超时时间
func WithEtcdDialTimeoutOpt(timeout time.Duration) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.DialTimeout = timeout
	}
}

// WithEtcdDialKeepAliveTimeOpt 设置连接保活时间
func WithEtcdDialKeepAliveTimeOpt(keepAliveTime time.Duration) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.DialKeepAliveTime = keepAliveTime
	}
}

// WithEtcdDialKeepAliveTimeoutOpt 设置连接保活超时时间
func WithEtcdDialKeepAliveTimeoutOpt(keepAliveTimeout time.Duration) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.DialKeepAliveTimeout = keepAliveTimeout
	}
}

// WithEtcdMaxCallSendMsgSizeOpt 设置最大发送消息大小
func WithEtcdMaxCallSendMsgSizeOpt(size int) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.MaxCallSendMsgSize = size
	}
}

// WithEtcdMaxCallRecvMsgSizeOpt 设置最大接收消息大小
func WithEtcdMaxCallRecvMsgSizeOpt(size int) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.MaxCallRecvMsgSize = size
	}
}

// WithEtcdTLSOpt 设置TLS配置
func WithEtcdTLSOpt(tlsConfig *tls.Config) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.TLS = tlsConfig
	}
}

// WithEtcdUsernameOpt 设置用户名
func WithEtcdUsernameOpt(username string) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.Username = username
	}
}

// WithEtcdPasswordOpt 设置密码
func WithEtcdPasswordOpt(password string) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.Password = password
	}
}

// WithEtcdRejectOldClusterOpt 设置是否拒绝旧集群
func WithEtcdRejectOldClusterOpt(reject bool) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.RejectOldCluster = reject
	}
}

// WithEtcdDialOptionsOpt 设置拨号选项
func WithEtcdDialOptionsOpt(dialOptions []grpc.DialOption) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.DialOptions = dialOptions
	}
}

// WithEtcdContextOpt 设置上下文
func WithEtcdContextOpt(ctx context.Context) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.Context = ctx
	}
}

// WithEtcdLoggerOpt 设置日志记录器
func WithEtcdLoggerOpt(logger *zap.Logger) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.Logger = logger
	}
}

// WithEtcdLogConfigOpt 设置日志配置
func WithEtcdLogConfigOpt(logConfig *zap.Config) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.LogConfig = logConfig
	}
}

// WithEtcdPermitWithoutStreamOpt 设置是否允许无流连接
func WithEtcdPermitWithoutStreamOpt(permit bool) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.PermitWithoutStream = permit
	}
}

// WithEtcdMaxUnaryRetriesOpt 设置最大一元调用重试次数
func WithEtcdMaxUnaryRetriesOpt(maxRetries uint) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.MaxUnaryRetries = maxRetries
	}
}

// WithEtcdBackoffWaitBetweenOpt 设置重试间隔等待时间
func WithEtcdBackoffWaitBetweenOpt(backoffWait time.Duration) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.BackoffWaitBetween = backoffWait
	}
}

// WithEtcdBackoffJitterFractionOpt 设置重试抖动分数
func WithEtcdBackoffJitterFractionOpt(jitterFraction float64) WithEtcdOpt {
	return func(config *clientv3.Config) {
		config.BackoffJitterFraction = jitterFraction
	}
}

// InitEtcd 初始化etcd客户端
func InitEtcd(opts ...WithEtcdOpt) error {
	config := &clientv3.Config{}
	for _, opt := range opts {
		opt(config)
	}
	client, err := clientv3.New(*config)
	if err != nil {
		return err
	}
	ETCDClient.Client = client
	return nil
}

// NewLock 初始化etcd分布式锁
func (e *Etcd) NewLock(key string, opts ...concurrency.SessionOption) (*concurrency.Mutex, error) {
	session, err := concurrency.NewSession(e.Client, opts...)
	if err != nil {
		return nil, err
	}
	mutex := concurrency.NewMutex(session, key)
	return mutex, nil
}
