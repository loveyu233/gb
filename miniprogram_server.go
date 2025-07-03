package gb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func (w *WXMini) RegisterHandlers(r *gin.RouterGroup) {
	r.Use(GinLogSetModuleName("微信小程序"))
	r.POST("/wx/login", GinLogSetOptionName("小程序登录", w.IsSaveHandlerLog), w.login)
}

type Phone struct {
	PhoneNumber string `json:"phoneNumber"`
}

func (w *WXMini) login(c *gin.Context) {
	var params struct {
		Code          string `binding:"required" json:"code"`
		EncryptedData string `json:"encrypted_data"`
		IvStr         string `json:"iv_str"`
	}
	if err := c.BindJSON(&params); err != nil {
		ResponseError(c, ErrInvalidParam)
		return
	}

	session, err := w.MiniProgramApp.Auth.Session(context.Background(), params.Code)
	if err != nil || session.ErrCode != 0 {
		ResponseError(c, ErrRequestWechat.WithMessage("获取微信小程序用户会话代码失败"))
		return
	}

	var (
		user   any
		exists bool
	)

	//检测用户是否注册
	user, exists, err = w.isExistsUser(session.UnionID)
	if err != nil {
		ResponseError(c, ErrDatabase.WithMessage("查询用户信息失败:%s", err.Error()))
		return
	}
	if !exists {
		if params.EncryptedData == "" {
			//如果是用户首次自动登录 没有授权手机号 就返回给用户open_id
			ResponseSuccess(c, map[string]interface{}{
				"open_id": session.OpenID,
			})
			return
		}
		//未注册,获取手机号
		data, _err := w.MiniProgramApp.Encryptor.DecryptData(params.EncryptedData, session.SessionKey, params.IvStr)
		if _err != nil {
			ResponseError(c, ErrRequestWechat.WithMessage("获取微信小程序用户数据失败"))
			return
		}
		var info Phone
		err = json.Unmarshal(data, &info)
		if err != nil || info.PhoneNumber == "" {
			ResponseError(c, ErrRequestWechat.WithMessage("获取微信小程序用户手机号失败"))
			return
		}

		if user, err = w.createUser(info.PhoneNumber, session.UnionID, session.OpenID, getAreaCodeByIp(c.ClientIP()), c.ClientIP()); err != nil {
			ResponseError(c, ErrDatabase.WithMessage("创建用户信息失败:%s", err.Error()))
			return
		}
	}

	data, err := w.generateToken(user, session.SessionKey)
	if err != nil {
		ResponseError(c, ErrServerBusy.WithMessage("token生成失败:%s", err.Error()))
		return
	}
	ResponseSuccess(c, data)
}

type AreaCode struct {
	Ip        string `json:"ip"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	ProvinceS string `json:"provinceS"`
	City      string `json:"city"`
	CityS     string `json:"cityS"`
	AdCode    string `json:"adCode"`
}

func getAreaCodeByIp(ip string) (adCode string) {
	var code AreaCode
	res, err := resty.New().R().
		Get(fmt.Sprintf("https://api.xtjzx.cn/geo-tool-pub/loc?ip=%s", ip))
	if err != nil {
		return "0"
	}
	json.Unmarshal(res.Body(), &code)
	return code.AdCode
}

// APIWXACodeCreateQRCode 获取小程序二维码，适用于需要的码数量较少的业务场景,pagePath:可携带query参数
func (w *WXMini) APIWXACodeCreateQRCode(ctx context.Context, pagePath string, width int64) (*http.Response, error) {
	rs, err := w.MiniProgramApp.WXACode.CreateQRCode(ctx, pagePath, width)

	if err != nil {
		return nil, err
	}

	return rs, nil
}

type MiniCode struct {
	ctx        context.Context
	pagePath   string // 扫码进入的小程序页面路径，最大长度 1024 个字符，不能为空
	width      int64
	r, g, b    int64
	envVersion string // 要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版。
	autoColor  bool
	isHyaline  bool
}

func NewMiniCode(ctx context.Context) *MiniCode {
	return &MiniCode{
		ctx: ctx,
	}
}

func (m *MiniCode) SetPagePath(pagePath string) *MiniCode {
	m.pagePath = pagePath
	return m
}

func (m *MiniCode) SetWidth(w int64) *MiniCode {
	m.width = w
	return m
}
func (m *MiniCode) SetRGB(r, g, b int64) *MiniCode {
	m.r, m.g, m.b = r, g, b
	return m
}
func (m *MiniCode) SetEnvVersion(version string) *MiniCode {
	m.envVersion = version
	return m
}
func (m *MiniCode) SetAutoColor(autoColor bool) *MiniCode {
	m.autoColor = autoColor
	return m
}
func (m *MiniCode) SetIsHyaline(isHyaline bool) *MiniCode {
	m.isHyaline = isHyaline
	return m
}

// APIWXACodeGet 获取小程序码，适用于需要的码数量较少的业务场景,pagePath:可携带query参数
func (w *WXMini) APIWXACodeGet(code MiniCode) (*http.Response, error) {
	rs, err := w.MiniProgramApp.WXACode.Get(
		code.ctx,
		code.pagePath,
		code.width,
		code.autoColor,
		&power.HashMap{
			"r": code.r,
			"g": code.g,
			"b": code.b,
		},
		code.isHyaline,
		code.envVersion,
	)

	if err != nil {
		return nil, err
	}

	return rs, nil
}

type MiniCodeMiniUnlimitedCode struct {
	ctx        context.Context
	pagePath   string
	scene      string
	width      int64
	r, g, b    int64
	envVersion string // 要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版。
	autoColor  bool
	isHyaline  bool
	checkPage  bool
}

func NewMiniUnlimitedCode(ctx context.Context) *MiniCodeMiniUnlimitedCode {
	return &MiniCodeMiniUnlimitedCode{
		ctx: ctx,
	}
}

func (m *MiniCodeMiniUnlimitedCode) SetPagePath(pagePath string) *MiniCodeMiniUnlimitedCode {
	m.pagePath = pagePath
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetScene(scene string) *MiniCodeMiniUnlimitedCode {
	m.scene = scene
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetWidth(w int64) *MiniCodeMiniUnlimitedCode {
	m.width = w
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetRGB(r, g, b int64) *MiniCodeMiniUnlimitedCode {
	m.r, m.g, m.b = r, g, b
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetEnvVersion(version string) *MiniCodeMiniUnlimitedCode {
	m.envVersion = version
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetAutoColor(autoColor bool) *MiniCodeMiniUnlimitedCode {
	m.autoColor = autoColor
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetIsHyaline(isHyaline bool) *MiniCodeMiniUnlimitedCode {
	m.isHyaline = isHyaline
	return m
}
func (m *MiniCodeMiniUnlimitedCode) SetCheckPage(checkPage bool) *MiniCodeMiniUnlimitedCode {
	m.checkPage = checkPage
	return m
}

// APIWXACodeGetUnlimited 获取小程序码，适用于需要的码数量极多的业务场景,scene:携带的参数
func (w *WXMini) APIWXACodeGetUnlimited(code *MiniCodeMiniUnlimitedCode) (*http.Response, error) {
	rs, err := w.MiniProgramApp.WXACode.GetUnlimited(
		code.ctx,
		code.scene,
		code.pagePath,
		code.checkPage,
		code.envVersion,
		code.width,
		code.autoColor,
		&power.HashMap{
			"r": code.r,
			"g": code.g,
			"b": code.b,
		},
		code.isHyaline,
	)

	if err != nil {
		return nil, err
	}

	return rs, nil
}
