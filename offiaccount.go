package gb

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/gin-gonic/gin"
)

type WXOfficial struct {
	OfficialAccountApp *officialAccount.OfficialAccount
	subscribe          func(unionID, openID string) error
	unSubscribe        func(unionID, openID string) error
	pushHandler        func(c *gin.Context)
}

type WXOfficialImp interface {
	Subscribe(unionID, openID string) error
	UnSubscribe(unionID, openID string) error
	PushHandler(c *gin.Context)
}

type OfficialAccount struct {
	AppID         string                `json:"appID,omitempty"`
	AppSecret     string                `json:"appSecret,omitempty"`
	MessageToken  string                `json:"messageToken,omitempty"`
	MessageAesKey string                `json:"messageAesKey,omitempty"`
	ResponseType  string                `json:"responseType,omitempty"`
	Cache         kernel.CacheInterface `json:"cache,omitempty"`
	HttpDebug     bool                  `json:"httpDebug,omitempty"`
}

type OfficialAccountAppServiceConfig struct {
	OfficialAccount OfficialAccount
	WXOfficialImp   WXOfficialImp
}

func WXNewOfficialAccountAppService(conf OfficialAccountAppServiceConfig) (*officialAccount.OfficialAccount, error) {
	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:        conf.OfficialAccount.AppID,
		Secret:       conf.OfficialAccount.AppSecret,
		Token:        conf.OfficialAccount.MessageToken,
		AESKey:       conf.OfficialAccount.MessageAesKey,
		ResponseType: conf.OfficialAccount.ResponseType,
		Cache:        conf.OfficialAccount.Cache,
		HttpDebug:    conf.OfficialAccount.HttpDebug,
	})
	WX.WXOfficial.OfficialAccountApp = app
	WX.WXOfficial.subscribe = conf.WXOfficialImp.Subscribe
	WX.WXOfficial.unSubscribe = conf.WXOfficialImp.UnSubscribe
	WX.WXOfficial.pushHandler = conf.WXOfficialImp.PushHandler
	return app, err
}
