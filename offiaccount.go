package gb

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/gin-gonic/gin"
	"os"
)

type WXOfficial struct {
	OfficialAccountApp *officialAccount.OfficialAccount
	subscribe          func(unionID, openID string) error
	unSubscribe        func(unionID, openID string) error
	pushHandler        func(c *gin.Context)
}

type OfficialAccount struct {
	AppID         string                `json:"appID,omitempty"`
	AppSecret     string                `json:"appSecret,omitempty"`
	MessageToken  string                `json:"messageToken,omitempty"`
	MessageAesKey string                `json:"messageAesKey,omitempty"`
	Cache         kernel.CacheInterface `json:"cache,omitempty"`
	HttpDebug     bool                  `json:"httpDebug,omitempty"`
}

func WXNewOfficialAccountAppService(conf OfficialAccount,
	subscribe func(unionID, openID string) error,
	unSubscribe func(unionID, openID string) error,
	pushHandler func(c *gin.Context)) (*officialAccount.OfficialAccount, error) {
	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:        conf.AppID,
		Secret:       conf.AppSecret,
		Token:        conf.MessageToken,
		AESKey:       conf.MessageAesKey,
		ResponseType: os.Getenv("response_type"),
		Cache:        conf.Cache,
		HttpDebug:    conf.HttpDebug,
	})
	WX.WXOfficial.OfficialAccountApp = app
	WX.WXOfficial.subscribe = subscribe
	WX.WXOfficial.unSubscribe = unSubscribe
	WX.WXOfficial.pushHandler = pushHandler
	return app, err
}
