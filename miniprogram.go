package gb

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

type WXMini struct {
	MiniProgramApp *miniProgram.MiniProgram
	isExistsUser   func(UnionID string) (user any, exists bool, err error)                                 // user:返回用户信息和token,exists:是否存在该用户,err错误
	createUser     func(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error) // 返回创建的用户信息
	generateToken  func(user any, sessionKey string) (data any, err error)
}

type MiniProgram struct {
	AppID             string                `json:"appID,omitempty"`
	Secret            string                `json:"secret,omitempty"`
	RedisAddr         string                `json:"redisAddr,omitempty"`
	MessageToken      string                `json:"messageToken,omitempty"`
	MessageAesKey     string                `json:"messageAesKey,omitempty"`
	VirtualPayAppKey  string                `json:"virtualPayAppKey,omitempty"`
	VirtualPayOfferID string                `json:"virtualPayOfferID,omitempty"`
	Env               string                `json:"env,omitempty"`
	Cache             kernel.CacheInterface `json:"cache,omitempty"`
	Log               miniProgram.Log       `json:"log"`
}

// WXMiniImp 定义用户相关操作接口
type WXMiniImp interface {
	// IsExistsUser 检查用户是否存在
	IsExistsUser(unionID string) (user any, exists bool, err error)

	// CreateUser 创建新用户
	CreateUser(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error)

	// GenerateToken 生成用户token
	GenerateToken(user any, sessionKey string) (data any, err error)
}

// MiniProgramConfig 小程序配置结构体
type MiniProgramConfig struct {
	// 基础配置
	AppID  string
	Secret string

	// 消息相关配置
	MessageToken  string
	MessageAESKey string

	// 虚拟支付相关配置
	VirtualPayAppKey  string
	VirtualPayOfferID string

	// 其他配置
	Log       miniProgram.Log
	Cache     kernel.CacheInterface
	HTTPDebug bool
	Debug     bool
}

// MiniProgramServiceConfig 小程序服务配置
type MiniProgramServiceConfig struct {
	MiniProgram MiniProgramConfig
	WXMiniImp   WXMiniImp
}

func WXNewMiniProgramService(config MiniProgramServiceConfig) (*miniProgram.MiniProgram, error) {
	app, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:        config.MiniProgram.AppID,
		Secret:       config.MiniProgram.Secret,
		ResponseType: response.TYPE_MAP,
		Token:        config.MiniProgram.MessageToken,
		AESKey:       config.MiniProgram.MessageAESKey,
		AppKey:       config.MiniProgram.VirtualPayAppKey,
		OfferID:      config.MiniProgram.VirtualPayOfferID,
		Log:          config.MiniProgram.Log,
		Cache:        config.MiniProgram.Cache,
		HttpDebug:    config.MiniProgram.HTTPDebug,
		Debug:        config.MiniProgram.Debug,
	})

	if err != nil {
		return nil, err
	}

	WX.WXMini.MiniProgramApp = app
	WX.WXMini.isExistsUser = config.WXMiniImp.IsExistsUser
	WX.WXMini.createUser = config.WXMiniImp.CreateUser
	WX.WXMini.generateToken = config.WXMiniImp.GenerateToken

	return WX.WXMini.MiniProgramApp, nil
}
