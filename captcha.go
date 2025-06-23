package gb

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/rotate"
	"log"
	"sync"
	"time"
)

type CaptchaClient struct {
	client            *redis.Client
	storeMap          sync.Map
	rotateCapt        rotate.Captcha
	ttl               int64
	rotateCaptPadding int
	maxRetriesNumber  int
}

var Captcha = new(CaptchaClient)

type CaptchaClientOption func(*CaptchaClient)

// WithCaptchaRedisClient 设置redis客户端, ttl过期时间单位秒
func WithCaptchaRedisClient(client *redis.Client, ttl int64, maxRetriesNumber int) CaptchaClientOption {
	return func(c *CaptchaClient) {
		c.client = client
		c.ttl = ttl
		c.maxRetriesNumber = maxRetriesNumber
	}
}

// WithCaptchaRotateCapt rotateCaptPadding误差范围,初始化旋转验证码 默认为20-330
func WithCaptchaRotateCapt(rotateCaptPadding int, val ...option.RangeVal) CaptchaClientOption {
	var rangeVal []option.RangeVal
	if len(val) > 0 {
		rangeVal = val
	} else {
		rangeVal = []option.RangeVal{
			{Min: 20, Max: 330},
		}
	}
	return func(cc *CaptchaClient) {
		builder := rotate.NewBuilder(rotate.WithRangeAnglePos(rangeVal))
		imgs, err := imagesv2.GetImages()
		if err != nil {
			panic(err)
		}

		builder.SetResources(
			rotate.WithImages(imgs),
		)
		cc.rotateCapt = builder.Make()
		cc.rotateCaptPadding = rotateCaptPadding
	}
}

// InitCaptchaClient 如果不设置WithCaptchaRedisClient则使用内存保存key,过期时间不设置默认为120s
func InitCaptchaClient(opts ...CaptchaClientOption) {
	for _, opt := range opts {
		opt(Captcha)
	}
	if Captcha.client == nil {
		panic("redis client is nil")
	}

	if Captcha.ttl == 0 {
		Captcha.ttl = 120
	}

	if Captcha.maxRetriesNumber == 0 {
		Captcha.maxRetriesNumber = 3
	}
}

type RotateCaptInfo struct {
	Block       *rotate.Block
	MasterImage string
	ThumbImage  string
	Key         string
}

func (c *CaptchaClient) RotateCaptKey(userID string) string {
	return fmt.Sprintf("rotate-captkey:%s:%d:%s", userID, time.Now().UnixMicro(), GetXID())
}

func (c *CaptchaClient) RotateCaptErrKey(key string) string {
	return fmt.Sprintf("%s:err", key)
}

func (c *CaptchaClient) saveKV(key string, value any) error {
	return c.client.Set(context.Background(), key, value, time.Second*time.Duration(c.ttl)).Err()
}

func (c *CaptchaClient) delK(key string) {
	c.client.Del(context.Background(), key)
	c.client.Del(context.Background(), c.RotateCaptErrKey(key))
}

func (c *CaptchaClient) getK(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *CaptchaClient) RotateCapt(userID string) (*RotateCaptInfo, error) {
	var resp = new(RotateCaptInfo)
	captData, err := c.rotateCapt.Generate()
	if err != nil {
		log.Fatalln(err)
	}

	dotData := captData.GetData()
	if dotData == nil {
		return nil, errors.New("generate err")
	}
	resp.Block = dotData

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	tBase64, err = captData.GetThumbImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}

	resp.MasterImage = mBase64
	resp.ThumbImage = tBase64

	resp.Key = c.RotateCaptKey(userID)
	if err = c.client.Set(context.Background(), resp.Key, resp.Block.Angle, time.Second*time.Duration(c.ttl)).Err(); err != nil {
		return nil, err
	}

	return resp, nil
}

// RotateCaptVerify 校验是否正确,verify:是否正确,beyondErr:是否超出最大重试次数,err:错误
func (c *CaptchaClient) RotateCaptVerify(key string, srcAngle int) (verify bool, beyondErr bool, err error) {
	result, err := c.getK(c.RotateCaptErrKey(key))
	if err != nil {
		return false, false, err
	}

	var number = cast.ToInt(result)
	if number > c.maxRetriesNumber {
		c.delK(key)
		return false, true, nil
	}

	value, err := c.client.Get(context.Background(), key).Int()
	if err != nil || value == 0 {
		return false, false, err
	}

	validate := rotate.Validate(srcAngle, value, c.rotateCaptPadding)
	if !validate {
		if err = c.saveKV(c.RotateCaptErrKey(key), number+1); err != nil {
			return false, false, err
		}
		return false, false, nil
	} else {
		c.delK(key)
	}

	return true, false, nil
}
