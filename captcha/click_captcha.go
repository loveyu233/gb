package captcha

import (
	"encoding/json"
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/click"
	"image"
	"log"
)

// ClickCaptcha 点击验证码
type ClickCaptcha struct {
	cache CacheImpl
	fonts *truetype.Font
	imgs  []image.Image
}

// NewClickCaptcha 创建点击验证码生成器
func NewClickCaptcha(cache CacheImpl) *ClickCaptcha {
	// 加载字体
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		log.Fatalf("加载字体失败: %v", err)
	}

	// 加载背景图片
	imgs, err := imagesv2.GetImages()
	if err != nil {
		log.Fatalf("加载背景图片失败: %v", err)
	}

	return &ClickCaptcha{
		cache: cache,
		fonts: fonts,
		imgs:  imgs,
	}
}

// Generate 生成验证码
func (c *ClickCaptcha) Generate(cacheKey string, config CaptchaConfig) (*CaptchaResponse, error) {
	var clickOption []click.Option
	if config.ImageSize != nil {
		clickOption = append(clickOption, click.WithImageSize(*config.ImageSize))
	}
	if config.ThumbImageSize != nil {
		clickOption = append(clickOption, click.WithRangeThumbImageSize(*config.ThumbImageSize))
	}
	builder := click.NewBuilder(clickOption...)
	// 设置资源
	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{c.fonts}),
		click.WithBackgrounds(c.imgs),
	)
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
		MasterWidth:       options.GetImageSize().Width,
		MasterHeight:      options.GetImageSize().Height,
		ThumbWidth:        options.GetThumbImageSize().Width,
		ThumbHeight:       options.GetThumbImageSize().Height,
	}

	// 缓存验证数据
	if err := c.cache.SetCaptcha(response.CaptchaKey, string(dots), config.CacheTimeout); err != nil {
		return nil, err
	}

	return response, nil
}

// Verify 验证验证码
func (c *ClickCaptcha) Verify(key string, captchaData string, validationTolerance int) error {
	// 获取缓存数据
	cachedData, err := c.cache.GetCaptcha(key)
	if err != nil {
		return ErrCaptchaNotFound
	}

	// 解析缓存数据
	var dotData map[int]*click.Dot
	if err := json.Unmarshal([]byte(cachedData), &dotData); err != nil {
		return ErrInvalidData
	}

	// 解析验证数据
	var verifyData map[int]*click.Dot
	if err := json.Unmarshal([]byte(captchaData), &verifyData); err != nil {
		return ErrInvalidData
	}

	// 验证点击数量
	if len(verifyData) != len(dotData) {
		return ErrCaptchaVerify
	}

	// 验证每个点击位置
	for k, expectedDot := range dotData {
		actualDot, exists := verifyData[k+1]
		if !exists {
			return ErrCaptchaVerify
		}

		if !click.Validate(
			actualDot.X, actualDot.Y,
			expectedDot.X, expectedDot.Y,
			expectedDot.Width, expectedDot.Height,
			validationTolerance,
		) {
			return ErrCaptchaVerify
		}
	}

	return nil
}

// GetType 获取验证码类型
func (c *ClickCaptcha) GetType() CaptchaType {
	return "click"
}
