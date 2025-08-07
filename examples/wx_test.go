package examples

import (
	"testing"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/user/response"
	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

type PayImp struct {
}

func (p PayImp) PayNotify(orderId string, attach string) error {
	//TODO implement me
	panic("implement me")
}

func (p PayImp) RefundNotify(orderId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PayImp) Pay(c *gin.Context) (*gb.PayRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p PayImp) Refund(c *gin.Context) (*gb.RefundRequest, error) {
	//TODO implement me
	panic("implement me")
}

func TestWXPay(t *testing.T) {
	gb.InitWXWXPaymentApp(gb.WXPaymentAppConfig{
		Payment:  gb.Payment{},
		WXPayImp: PayImp{},
	})

	// 支付
	gb.InsWX.WXPay.Pay(&gb.PayRequest{
		Price:       0,
		Description: "",
		Ip:          "",
		Openid:      "",
		Attach:      "",
		NotifyUrl:   "",
		OutTradeNo:  "",
	})

	// 订单查询
	gb.InsWX.WXPay.QueryOrder("")

	// 退款
	gb.InsWX.WXPay.Refund(&gb.RefundRequest{
		OrderId:    "",
		TotalFee:   0,
		RefundFee:  0,
		RefundDesc: "",
		NotifyUrl:  "",
	})

	// 退款查询
	gb.InsWX.WXPay.QueryRefundOrder("")

	engine := gin.Default()
	group := engine.Group("/wx")
	// 注册微信支付和微信退款回调api
	gb.InsWX.WXPay.RegisterHandlers(group)
}

type WXMiniImp struct {
}

func (W WXMiniImp) IsExistsUser(unionID string) (user any, exists bool, err error) {
	//TODO 查询当前用户是否存在
	panic("implement me")
}

func (W WXMiniImp) CreateUser(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error) {
	//TODO 创建用户
	panic("implement me")
}

func (W WXMiniImp) GenerateToken(user any, sessionKey string) (data any, err error) {
	//TODO 根据用户信息生成token返回信息
	panic("implement me")
}

func TestWXMini(t *testing.T) {
	gb.InitWXMiniProgramService(gb.MiniProgramServiceConfig{
		MiniProgram: gb.MiniProgramConfig{},
		WXMiniImp:   WXMiniImp{},
	})

	engine := gin.Default()
	group := engine.Group("/wx")
	// 注册微信小程序登录回调api
	gb.InsWX.WXMini.RegisterHandlers(group)
}

type WXOfficialImp struct {
}

func (W WXOfficialImp) Subscribe(rs *response.ResponseGetUserInfo, event contract.EventInterface) error {
	//TODO 订阅后回调
	panic("implement me")
}

func (W WXOfficialImp) UnSubscribe(rs *response.ResponseGetUserInfo, event contract.EventInterface) error {
	//TODO 取消订阅回调
	panic("implement me")
}

func (W WXOfficialImp) PushHandler(c *gin.Context) (toUsers []string, message string) {
	//TODO 发送消息
	panic("implement me")
}

func TestOffia(t *testing.T) {
	gb.InitWXOfficialAccountAppService(gb.OfficialAccountAppServiceConfig{
		OfficialAccount: gb.OfficialAccount{},
		WXOfficialImp:   WXOfficialImp{},
	})

	// 发送消息
	gb.InsWX.WXOfficial.PushTemplateMessage("", "", "")
}
