package captcha

import (
	"encoding/json"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"image"
	"log"
)

// RotateCaptcha 旋转验证码
type RotateCaptcha struct {
	cache CacheImpl
	imgs  []image.Image
}

// NewRotateCaptcha 创建旋转验证码生成器
func NewRotateCaptcha(cache CacheImpl) *RotateCaptcha {
	// 加载背景图片
	imgs, err := imagesv2.GetImages()
	if err != nil {
		log.Fatalf("加载背景图片失败: %v", err)
	}

	return &RotateCaptcha{
		cache: cache,
		imgs:  imgs,
	}
}

// Generate 生成验证码
func (r *RotateCaptcha) Generate(cacheKey string, config CaptchaConfig) (*CaptchaResponse, error) {
	var rotateOption []rotate.Option
	if config.SquareSize != 0 {
		rotateOption = append(rotateOption, rotate.WithImageSquareSize(config.SquareSize))
	}
	rotateOption = append(rotateOption, rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: 20, Max: 330},
	}))
	builder := rotate.NewBuilder(rotateOption...)
	// 设置资源
	builder.SetResources(rotate.WithImages(r.imgs))

	captcha := builder.Make()

	captData, err := captcha.Generate()
	if err != nil {
		return nil, err
	}

	dotData := captData.GetData()
	if dotData == nil {
		return nil, ErrCaptchaGenerate
	}

	// 序列化数据
	dots, err := json.Marshal(dotData)
	if err != nil {
		return nil, err
	}

	// 生成Base64图片
	mBase64, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, err
	}

	tBase64, err := captData.GetThumbImage().ToBase64()
	if err != nil {
		return nil, err
	}

	// 创建响应
	options := captcha.GetOptions()
	response := &CaptchaResponse{
		CaptchaKey:        cacheKey,
		MasterImageBase64: mBase64,
		ThumbImageBase64:  tBase64,
		MasterWidth:       options.GetImageSize(),
		MasterHeight:      options.GetImageSize(),
		ThumbWidth:        dotData.Width,
		ThumbHeight:       dotData.Height,
		ThumbSize:         dotData.Width,
	}

	// 缓存验证数据
	if err := r.cache.SetCaptcha(response.CaptchaKey, string(dots), config.CacheTimeout); err != nil {
		return nil, err
	}

	return response, nil
}

// Verify 验证验证码
func (r *RotateCaptcha) Verify(key string, captchaData string, validationTolerance int) error {
	// 获取缓存数据
	cachedData, err := r.cache.GetCaptcha(key)
	if err != nil {
		return ErrCaptchaNotFound
	}

	// 解析缓存数据
	var dotData rotate.Block
	if err := json.Unmarshal([]byte(cachedData), &dotData); err != nil {
		return ErrInvalidData
	}

	// 解析验证数据
	var verifyData rotate.Block
	if err := json.Unmarshal([]byte(captchaData), &verifyData); err != nil {
		return ErrInvalidData
	}

	// 验证角度
	if !rotate.Validate(verifyData.Angle, dotData.Angle, validationTolerance) {
		return ErrCaptchaVerify
	}

	return nil
}

// GetType 获取验证码类型
func (r *RotateCaptcha) GetType() CaptchaType {
	return "rotate"
}
