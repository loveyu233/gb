package gb

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/user/response"
	"github.com/gin-gonic/gin"
)

type WXOfficial struct {
	OfficialAccountApp *officialAccount.OfficialAccount
	subscribe          func(rs *response.ResponseGetUserInfo, event contract.EventInterface) error
	unSubscribe        func(rs *response.ResponseGetUserInfo, event contract.EventInterface) error
	pushHandler        func(c *gin.Context) (toUsers []string, message string)
}

type WXOfficialImp interface {
	Subscribe(rs *response.ResponseGetUserInfo, event contract.EventInterface) error
	UnSubscribe(rs *response.ResponseGetUserInfo, event contract.EventInterface) error
	PushHandler(c *gin.Context) (toUsers []string, message string)
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

func WXNewOfficialAccountAppService(conf OfficialAccountAppServiceConfig) error {
	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:        conf.OfficialAccount.AppID,
		Secret:       conf.OfficialAccount.AppSecret,
		Token:        conf.OfficialAccount.MessageToken,
		AESKey:       conf.OfficialAccount.MessageAesKey,
		ResponseType: conf.OfficialAccount.ResponseType,
		Cache:        conf.OfficialAccount.Cache,
		HttpDebug:    conf.OfficialAccount.HttpDebug,
	})
	if err != nil {
		return err
	}
	WX.WXOfficial.OfficialAccountApp = app
	WX.WXOfficial.subscribe = conf.WXOfficialImp.Subscribe
	WX.WXOfficial.unSubscribe = conf.WXOfficialImp.UnSubscribe
	WX.WXOfficial.pushHandler = conf.WXOfficialImp.PushHandler
	return nil
}
