package examples

import (
	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
	"testing"
)

func TestWXPay(t *testing.T) {
	gb.WXNewWXPaymentApp(gb.Payment{}, func(orderId string, attach string) error {
		return nil
	}, func(orderId string) error {
		return nil
	})
	// 付款
	gb.WX.WXPay.Pay(&gb.PayRequest{
		Price:       0,
		Description: "",
		Ip:          "",
		Openid:      "",
		Attach:      "",
		NotifyUrl:   "",
	})
	// 退款
	gb.WX.WXPay.Refund(&gb.RefundRequest{})
	// 注册路由
	gb.PublicRoutes = append(gb.PublicRoutes, gb.WX.WXPay.WXPayHttpGroup)
}

func TestWXMini(t *testing.T) {
	gb.WXNewMiniMiniProgramService(gb.MiniProgram{}, func(UnionID string) (user any, exists bool, err error) {
		return nil, false, err
	}, func(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error) {
		return nil, err
	}, func(user any) (data any, err error) {
		return nil, err
	})

	gb.PublicRoutes = append(gb.PublicRoutes, gb.WX.WXMini.WXMiniHttpGroup)
}

func TestOffia(t *testing.T) {
	gb.WXNewOfficialAccountAppService(gb.OfficialAccount{}, func(unionID, openID string) error {
		return nil
	}, func(unionID, openID string) error {
		return nil
	}, func(c *gin.Context) {

	})

	// 发送模版消息
	gb.WX.WXOfficial.PushTemplateMessage("", "", "")
}
