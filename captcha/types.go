package captcha

import (
	"errors"
	"github.com/wenlng/go-captcha/v2/base/option"
	"time"
)

// CaptchaType 验证码类型
type CaptchaType string

// 错误定义
var (
	ErrCaptchaVerify   = errors.New("验证码验证失败")
	ErrCaptchaGenerate = errors.New("验证码生成失败")
	ErrCaptchaNotFound = errors.New("验证码未找到")
	ErrInvalidData     = errors.New("验证码数据无效")
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	CacheTimeout        time.Duration `json:"cache_timeout"`
	ValidationTolerance int           `json:"validation_tolerance"`
	ImageSize           *option.Size  `json:"image_size"`       // 主图大小
	ThumbImageSize      *option.Size  `json:"thumb_image_size"` // 缩略图大小
	SquareSize          int           `json:"square_size"`      // 旋转验证码大小
}

// DefaultConfig 默认配置
var DefaultConfig = CaptchaConfig{
	CacheTimeout:        120 * time.Second,
	ValidationTolerance: 10,
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaKey        string `json:"captcha_key"`
	MasterImageBase64 string `json:"master_image_base64"`
	ThumbImageBase64  string `json:"thumb_image_base64"`
	MasterWidth       int    `json:"master_width"`
	MasterHeight      int    `json:"master_height"`
	ThumbWidth        int    `json:"thumb_width"`
	ThumbHeight       int    `json:"thumb_height"`
	ThumbSize         int    `json:"thumb_size,omitempty"`
	DisplayX          int    `json:"display_x,omitempty"`
	DisplayY          int    `json:"display_y,omitempty"`
}

// CaptchaVerifyRequest 验证请求
type CaptchaVerifyRequest struct {
	CaptchaKey          string `json:"captcha_key"`
	CaptchaData         string `json:"captcha_data"`
	ValidationTolerance int    `json:"validation_tolerance"`
}

// CaptchaInterface 验证码接口
type CaptchaInterface interface {
	Generate(cacheKey string, config CaptchaConfig) (*CaptchaResponse, error)
	Verify(key string, data string, validationTolerance int) error
	GetType() CaptchaType
}

type CacheImpl interface {
	SetCaptcha(key string, value any, expiration time.Duration) error
	GetCaptcha(key string) (string, error)
	DelCaptcha(key string) error
}
