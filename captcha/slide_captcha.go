package captcha

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
	"image"
	"log"
)

// SlideCaptcha 滑动验证码
type SlideCaptcha struct {
	cache     CacheImpl
	newGraphs []*slide.GraphImage
	imgs      []image.Image
}

// NewSlideCaptcha 创建滑动验证码生成器
func NewSlideCaptcha(cache CacheImpl) *SlideCaptcha {
	// 加载背景图片
	imgs, err := imagesv2.GetImages()
	if err != nil {
		log.Fatalf("加载背景图片失败: %v", err)
	}

	// 加载拼图图片
	graphs, err := tiles.GetTiles()
	if err != nil {
		log.Fatalf("加载拼图图片失败: %v", err)
	}

	// 转换图片格式
	newGraphs := make([]*slide.GraphImage, 0, len(graphs))
	for _, graph := range graphs {
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	return &SlideCaptcha{
		cache:     cache,
		newGraphs: newGraphs,
		imgs:      imgs,
	}
}

// Generate 生成验证码
func (s *SlideCaptcha) Generate(cacheKey string, config CaptchaConfig) (*CaptchaResponse, error) {
	var slideOption []slide.Option
	if config.ImageSize != nil {
		slideOption = append(slideOption, slide.WithImageSize(*config.ImageSize))
	}

	builder := slide.NewBuilder(slideOption...)

	// 设置资源
	builder.SetResources(
		slide.WithGraphImages(s.newGraphs),
		slide.WithBackgrounds(s.imgs),
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

	tBase64, err := captData.GetTileImage().ToBase64()
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
		ThumbWidth:        dotData.Width,
		ThumbHeight:       dotData.Height,
		DisplayX:          dotData.DX,
		DisplayY:          dotData.DY,
	}

	// 缓存验证数据
	if err := s.cache.SetCaptcha(response.CaptchaKey, string(dots), config.CacheTimeout); err != nil {
		return nil, err
	}

	return response, nil
}

func compressBase64(data string) (string, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write([]byte(data))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// 将压缩后的数据再次base64编码
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Verify 验证验证码
func (s *SlideCaptcha) Verify(key string, captchaData string, validationTolerance int) error {
	// 获取缓存数据
	cachedData, err := s.cache.GetCaptcha(key)
	if err != nil {
		return ErrCaptchaNotFound
	}

	// 解析缓存数据
	var dotData slide.Block
	if err := json.Unmarshal([]byte(cachedData), &dotData); err != nil {
		return ErrInvalidData
	}

	// 解析验证数据
	var verifyData slide.Block
	if err := json.Unmarshal([]byte(captchaData), &verifyData); err != nil {
		return ErrInvalidData
	}

	tx, ty := dotData.X, dotData.Y
	// 验证滑动位置
	if !slide.Validate(verifyData.X, verifyData.Y, tx, ty, validationTolerance) {
		return ErrCaptchaVerify
	}

	return nil
}

// GetType 获取验证码类型
func (s *SlideCaptcha) GetType() CaptchaType {
	return "slide"
}
