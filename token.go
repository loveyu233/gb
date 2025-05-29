package gb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Claims 是JWT的标准声明加自定义字段
type Claims struct {
	jwt.RegisteredClaims
	User interface{} `json:"user"`
}

// TokenService 提供Token相关操作的接口
type TokenService interface {
	Generate(user interface{}, expiration time.Duration) (string, error)
	Validate(tokenStr string) (*Claims, error)
	Invalidate(tokenID string) error
	IsInvalidated(tokenID string) bool
}

// JWTTokenService 实现TokenService接口
type JWTTokenService struct {
	secret                    string
	redisClient               *redis.Client
	signingMethod             jwt.SigningMethod
	enableRedisCheckBlacklist bool // 是否启动redis黑名单
	blacklistKeyFn            func(tokenID string) string
	enableRedisCheck          bool                        // 是否启用Redis校验token
	redisTokenKeyFn           func(tokenID string) string // Redis中存储有效token的key生成函数
}

// TokenServiceOption 提供配置JWTTokenService的函数选项
type TokenServiceOption func(*JWTTokenService)

// WithRedisClient 设置redis客户端
func WithRedisClient(client *redis.Client) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.redisClient = client
	}
}

// WithRedisBlacklist 启用Redis黑名单
func WithRedisBlacklist(enabled bool, keyFn func(tokenID string) string) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.enableRedisCheckBlacklist = enabled
		if keyFn != nil {
			service.blacklistKeyFn = keyFn
		} else {
			service.blacklistKeyFn = func(tokenID string) string {
				return fmt.Sprintf("token:blacklist:%s", tokenID)
			}
		}
	}
}

// WithSigningMethod 设置签名方法
func WithSigningMethod(method jwt.SigningMethod) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.signingMethod = method
	}
}

// WithRedisTokenCheck 启用Redis校验token存在功能
func WithRedisTokenCheck(enabled bool, keyFn func(tokenID string) string) TokenServiceOption {
	return func(service *JWTTokenService) {
		service.enableRedisCheck = enabled
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
func (s *JWTTokenService) Generate(user interface{}, expiration time.Duration) (string, error) {
	now := time.Now()
	tokenID := generateUUID()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
		User: user,
	}

	token := jwt.NewWithClaims(s.signingMethod, claims)
	signedToken, err := token.SignedString([]byte(s.secret))

	// 如果启用了Redis校验并且有Redis客户端，则将token存入Redis
	if err == nil && s.enableRedisCheck && s.redisClient != nil {
		// 存储token到Redis，过期时间与token一致
		redisErr := s.redisClient.Set(context.Background(),
			s.redisTokenKeyFn(tokenID),
			1,
			expiration).Err()

		if redisErr != nil {
			return "", fmt.Errorf("存储token到Redis失败: %w", redisErr)
		}
	}

	return signedToken, err
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
	if s.enableRedisCheck && s.redisClient != nil {
		// 检查Redis中是否存在该token
		exists, err := s.redisClient.Exists(context.Background(), s.redisTokenKeyFn(claims.ID)).Result()
		if err != nil {
			return nil, fmt.Errorf("校验token存在性失败: %w", err)
		}
		if exists == 0 {
			return nil, errors.New("令牌不存在或已过期")
		}
	}

	return claims, nil
}

// Invalidate 将令牌添加到黑名单
func (s *JWTTokenService) Invalidate(tokenID string) error {
	if s.redisClient == nil {
		return errors.New("未配置Redis，无法使用撤销功能")
	}

	var err error

	// 如果启用了Redis校验，同时删除有效token记录
	if s.enableRedisCheckBlacklist {
		// 使用pipeline批量执行操作
		pipe := s.redisClient.Pipeline()
		pipe.Set(context.Background(), s.blacklistKeyFn(tokenID), 1, 7*24*time.Hour)
		pipe.Del(context.Background(), s.redisTokenKeyFn(tokenID))
		_, err = pipe.Exec(context.Background())
	} else {
		// 将令牌加入黑名单，设置一个合理的过期时间（如7天）
		err = s.redisClient.Set(context.Background(), s.blacklistKeyFn(tokenID), 1, 7*24*time.Hour).Err()
	}

	return err
}

// IsInvalidated 检查令牌是否在黑名单中
func (s *JWTTokenService) IsInvalidated(tokenID string) bool {
	if s.redisClient == nil {
		return false
	}

	exists, err := s.redisClient.Exists(context.Background(), s.blacklistKeyFn(tokenID)).Result()
	return err == nil && exists > 0
}

// GinAuthConfig 配置Gin认证中间件
type GinAuthConfig struct {
	TokenService    TokenService
	GetTokenStrFunc func(c *gin.Context) string
	HandleError     func(c *gin.Context, err error)
}

var defaultTokenSecret = "abcdef123456..."

// DefaultGinConfig 默认配置
var DefaultGinConfig = &GinAuthConfig{
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
	HandleError: func(c *gin.Context, err error) {
		ResponseError(c, ErrTokenInvalid)
		c.Abort()
	},
}

// GinAuth 创建一个Gin认证中间件,config为空则使用默认
func GinAuth(userPtr interface{}, config *GinAuthConfig) gin.HandlerFunc {
	// 合并默认配置
	if config == nil {
		config = &GinAuthConfig{}
	}

	if config.GetTokenStrFunc == nil {
		config.GetTokenStrFunc = DefaultGinConfig.GetTokenStrFunc
	}

	if config.HandleError == nil {
		config.HandleError = DefaultGinConfig.HandleError
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

		// 解析用户信息
		bytes, err := json.Marshal(claims.User)
		if err != nil {
			config.HandleError(c, errors.New("用户数据序列化失败"))
			return
		}

		if err = json.Unmarshal(bytes, userPtr); err != nil {
			config.HandleError(c, errors.New("用户数据解析失败"))
			return
		}

		// 将用户信息存入上下文
		c.Set("user", userPtr)
		c.Set("claims", claims)
		c.Next()
	}
}

// generateUUID 生成唯一标识符
func generateUUID() string {
	return uuid.NewString()
}

// 以下是辅助功能

// ExtractTokenClaims 从上下文中提取Claims
func ExtractTokenClaims(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	if tokenClaims, ok := claims.(*Claims); ok {
		return tokenClaims, true
	}

	return nil, false
}

// CombineMiddlewares 组合多个中间件
func CombineMiddlewares(middlewares ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, middleware := range middlewares {
			middleware(c)
			if c.IsAborted() {
				return
			}
		}
	}
}
