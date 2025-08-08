package gb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims 是JWT的标准声明加自定义字段
type Claims struct {
	jwt.RegisteredClaims
	Data any `json:"data"`
}

// TokenService 提供Token相关操作的接口
type TokenService interface {
	Generate(user any, expiration time.Duration) (string, error)
	DeleteRedisToken(token string) error
	Validate(tokenStr string) (*Claims, error)
}

// JWTTokenService 实现TokenService接口
type JWTTokenService struct {
	secret                    string
	signingMethod             jwt.SigningMethod
	enableRedisCheckBlacklist bool // 是否启动redis黑名单
	blacklistKeyFn            func(tokenID string) string
	enableRedisCheck          bool                        // 是否启用Redis校验token
	redisTokenKeyFn           func(tokenID string) string // Redis中存储有效token的key生成函数
}

// TokenServiceOption 提供配置JWTTokenService的函数选项
type TokenServiceOption func(*JWTTokenService)

// WithSigningMethod 设置签名方法
func WithSigningMethod(method jwt.SigningMethod) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.signingMethod = method
	}
}

// WithRedisTokenKey 启用Redis校验token存在功能
func WithRedisTokenKey(keyFn func(tokenID string) string) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.enableRedisCheck = true
		if keyFn != nil {
			service.redisTokenKeyFn = keyFn
		} else {
			service.redisTokenKeyFn = func(tokenID string) string {
				return fmt.Sprintf("token:valid:%s", tokenID)
			}
		}
	}
}

// NewJWTTokenService 创建一个新的JWTTokenService
func NewJWTTokenService(secret string, options ...TokenServiceOption) *JWTTokenService {
	if secret == "" {
		secret = defaultTokenSecret
	}
	service := &JWTTokenService{
		secret:           secret,
		signingMethod:    jwt.SigningMethodHS256,
		enableRedisCheck: false, // 默认不启用Redis校验
	}

	for _, option := range options {
		option(service)
	}

	return service
}

// Generate 生成JWT令牌
func (s *JWTTokenService) Generate(data any, expiration time.Duration) (string, error) {
	now := Now()
	tokenID := GetUUID()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
		Data: data,
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)
	signedToken, err := token.SignedString([]byte(s.secret))

	// 如果启用了Redis校验并且有Redis客户端，则将token存入Redis
	if err == nil && s.enableRedisCheck {
		if InsRedis == nil {
			return "", redisClientNilErr()
		}
		// 存储token到Redis，过期时间与token一致
		redisErr := InsRedis.Set(context.Background(),
			s.redisTokenKeyFn(tokenID),
			1,
			expiration).Err()

		if redisErr != nil {
			return "", fmt.Errorf("存储token到Redis失败: %w", redisErr)
		}
	}

	return signedToken, err
}

func (s *JWTTokenService) DeleteRedisToken(token string) error {
	validate, err := s.Validate(token)
	if err != nil {
		return err
	}
	return InsRedis.Del(Context(), s.redisTokenKeyFn(validate.ID)).Err()
}

// Validate 验证JWT令牌并返回声明
func (s *JWTTokenService) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("加密方法错误: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("令牌无效")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("转换JWT声明失败")
	}

	// 如果启用了Redis校验
	if s.enableRedisCheck {
		if InsRedis == nil {
			return nil, redisClientNilErr()
		}
		// 检查Redis中是否存在该token
		exists, err := InsRedis.Exists(context.Background(), s.redisTokenKeyFn(claims.ID)).Result()
		if err != nil {
			return nil, fmt.Errorf("校验token存在性失败: %w", err)
		}
		if exists == 0 {
			return nil, errors.New("令牌不存在或已过期")
		}
	}

	return claims, nil
}

// GinAuthConfig 配置Gin认证中间件
type GinAuthConfig struct {
	TokenService    TokenService
	GetTokenStrFunc func(c *gin.Context) string
	HandleError     func(c *gin.Context, err error)
}

var defaultTokenSecret = "abcdef123456..."

var DefaultGInTokenErrHandler = func(c *gin.Context, err error) {
	ResponseError(c, ErrTokenInvalid.WithMessage(err.Error()))
	c.Abort()
}

// DefaultGinTokenConfig 默认配置,生产一定不可以使用
var DefaultGinTokenConfig = &GinAuthConfig{
	TokenService: NewJWTTokenService(defaultTokenSecret),
	GetTokenStrFunc: func(c *gin.Context) string {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			// 尝试从查询参数获取
			auth = c.Query("token")
		}
		// 尝试从cookie获取
		if auth == "" {
			if cookie, err := c.Cookie("token"); err == nil {
				auth = cookie
			}
		}
		return strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	},
	HandleError: DefaultGInTokenErrHandler,
}

type TokenDefaultData struct {
	ID string `json:"id"`
}

// GinAuth 创建一个Gin认证中间件,config为空则使用默认
func GinAuth(config *GinAuthConfig) gin.HandlerFunc {
	if config == nil {
		panic("config不能为nil")
	}

	if config.GetTokenStrFunc == nil {
		config.GetTokenStrFunc = DefaultGinTokenConfig.GetTokenStrFunc
	}

	if config.HandleError == nil {
		config.HandleError = DefaultGinTokenConfig.HandleError
	}

	if config.TokenService == nil {
		panic("必须提供TokenService")
	}

	return func(c *gin.Context) {
		// 获取令牌
		tokenStr := config.GetTokenStrFunc(c)
		if tokenStr == "" {
			config.HandleError(c, errors.New("令牌不存在"))
			return
		}

		// 验证令牌
		claims, err := config.TokenService.Validate(tokenStr)
		if err != nil {
			config.HandleError(c, err)
			return
		}

		c.Set("token_info", claims)
		c.Set("token", tokenStr)
		c.Next()
	}
}
