package msg

import (
	"context"

	"github.com/ArtisanCloud/PowerLibs/v3/logger/contract"
	"github.com/ArtisanCloud/PowerLibs/v3/logger/drivers"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/basicService/subscribeMessage/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

type WXMiniConfig struct {
	AppID     string
	Secret    string
	HttpDebug bool
	Log       *WXMiniConfigLog
	Cache     kernel.CacheInterface
}

type WXMiniConfigLog struct {
	Driver contract.LoggerInterface
	Level  string
	File   string
	Error  string
	ENV    string
	Stdout bool
}

var (
	InsWxMiniMsg *WxMiniMsgClient
)

type WxMiniMsgClient struct {
	App *miniProgram.MiniProgram
}

func InitWxMini(config WXMiniConfig) error {
	cnf := &miniProgram.UserConfig{
		AppID:     config.AppID,
		Secret:    config.Secret,
		HttpDebug: config.HttpDebug,
		Cache:     config.Cache,
	}
	if config.Log != nil {
		cnf.Log = miniProgram.Log{
			Driver: config.Log.Driver,
			Level:  config.Log.Level,
			File:   config.Log.File,
			Error:  config.Log.Error,
			ENV:    config.Log.ENV,
			Stdout: config.Log.Stdout,
		}
	} else {
		cnf.Log = miniProgram.Log{
			Driver: &drivers.DummyLogger{},
			Stdout: false,
		}
	}

	program, err := miniProgram.NewMiniProgram(cnf)
	if err != nil {
		return err
	}
	InsWxMiniMsg = &WxMiniMsgClient{
		App: program,
	}
	return nil
}

type MiniProgramStateType string

const (
	// MiniProgramStateDeveloper 体验版
	MiniProgramStateDeveloper MiniProgramStateType = "developer"
	// MiniProgramStateTrial 开发板
	MiniProgramStateTrial MiniProgramStateType = "trial"
	// MiniProgramStateFormal 正式版
	MiniProgramStateFormal MiniProgramStateType = "formal"
)

type WxMiniMsgContent struct {
	ToUserOpenID     string
	TemplateID       string
	Page             string
	MiniProgramState MiniProgramStateType
	Data             map[string]map[string]any
}

func (w *WxMiniMsgClient) SubscribeMessageSend(ctx context.Context, content WxMiniMsgContent) (*response.ResponseMiniProgram, error) {
	var data = make(power.HashMap)
	for k, v := range content.Data {
		for k1, v1 := range v {
			data[k] = power.HashMap{
				k1: v1,
			}
		}

	}
	resp, err := w.App.SubscribeMessage.Send(ctx, &request.RequestSubscribeMessageSend{
		ToUser:           content.ToUserOpenID,
		TemplateID:       content.TemplateID,
		Page:             content.Page,
		MiniProgramState: string(content.MiniProgramState),
		Lang:             "zh_CN",
		Data:             &data,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
