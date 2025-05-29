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
	generateToken  func(user any) (data any, err error)
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

func WXNewMiniMiniProgramService(conf MiniProgram,
	isExistsUser func(UnionID string) (user any, exists bool, err error),
	createUser func(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error),
	generateToken func(user any) (data any, err error)) (*miniProgram.MiniProgram, error) {
	app, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:        conf.AppID,  // 小程序、公众号或者企业微信的appid
		Secret:       conf.Secret, // 商户号 appID
		ResponseType: response.TYPE_MAP,
		Token:        conf.MessageToken,
		AESKey:       conf.MessageAesKey,
		AppKey:       conf.VirtualPayAppKey,
		OfferID:      conf.VirtualPayOfferID,
		Log:          conf.Log,
		Cache:        conf.Cache,
		HttpDebug:    true,
		Debug:        false,
	})

	WX.WXMini.MiniProgramApp = app
	WX.WXMini.isExistsUser = isExistsUser
	WX.WXMini.createUser = createUser
	WX.WXMini.generateToken = generateToken
	return WX.WXMini.MiniProgramApp, err
}
