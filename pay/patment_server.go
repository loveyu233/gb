package pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	nRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/notify/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/response"
	rRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/refund/request"
	rResponse "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/refund/response"
	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

func (wx *WXPay) RegisterHandlers(r *gin.RouterGroup) {
	r.Use(gb.GinLogSetModuleName("微信支付"))
	r.POST("/wx/notify/payment", gb.GinLogSetOptionName("支付异步回调", wx.IsSaveHandlerLog), wx.wxPayCallback)
	r.POST("/wx/notify/refund", gb.GinLogSetOptionName("退款异步回调", wx.IsSaveHandlerLog), wx.wxRefundCallback)
	r.POST("/wx/pay", gb.GinLogSetOptionName("支付请求", wx.IsSaveHandlerLog), wx.pay)
	r.POST("/wx/refund", gb.GinLogSetOptionName("退款请求", wx.IsSaveHandlerLog), wx.refund)
}

func (wx *WXPay) pay(c *gin.Context) {
	payRequest, err := wx.payHandler(c)
	if err != nil {
		gb.ResponseError(c, gb.ErrRequestWechatPay.WithMessage(err.Error()))
		return
	}
	pay, err := wx.Pay(payRequest)
	if err != nil {
		gb.ResponseError(c, gb.ErrRequestWechatPay.WithMessage(err.Error()))
		return
	}

	gb.ResponseSuccess(c, pay)
}

func (wx *WXPay) refund(c *gin.Context) {
	refundRequest, err := wx.refundHandler(c)
	if err != nil {
		gb.ResponseError(c, gb.ErrRequestWechatPay.WithMessage(err.Error()))
		return
	}

	refund, err := wx.Refund(refundRequest)
	if err != nil {
		gb.ResponseError(c, gb.ErrRequestWechatPay.WithMessage(err.Error()))
	}

	gb.ResponseSuccess(c, refund)
}

func (wx *WXPay) wxPayCallback(c *gin.Context) {
	res, err := wx.payNotify(c.Request, wx.payNotifyHandler)
	if err != nil {
		c.XML(500, err.Error())
		return
	}

	err = res.Write(c.Writer)
	if err != nil {
		c.XML(500, err.Error())
		return
	}
}

func (wx *WXPay) wxRefundCallback(c *gin.Context) {
	res, err := wx.refundNotify(c.Request, wx.refundNotifyHandler)
	if err != nil {
		c.XML(500, err.Error())
		return
	}

	err = res.Write(c.Writer)
	if err != nil {
		c.XML(500, err.Error())
		return
	}
}

type PayRequest struct {
	Price       int64  `json:"price"`
	Description string `json:"description"`
	Ip          string `json:"ip,omitempty"`
	Openid      string `json:"openid"`
	Attach      string `json:"attach"`
	NotifyUrl   string `json:"notify_url"`
	OutTradeNo  string `json:"out_trade_no"` // 可以使用自带的snowflake.GetId()
}

type WxPayResp struct {
	AppId      string `json:"appId"`
	NonceStr   string `json:"nonceStr"`
	Package    string `json:"package"`
	PaySign    string `json:"paySign"`
	SignType   string `json:"signType"`
	TimeStamp  string `json:"timeStamp"`
	OutTradeNo string `json:"outTradeNo"`
	BizOrder   string `json:"bizOrder"`
}

// Pay 支付
func (wx *WXPay) Pay(req *PayRequest) (*WxPayResp, error) {
	options := &request.RequestJSAPIPrepay{
		Amount: &request.JSAPIAmount{
			Total:    int(req.Price),
			Currency: "CNY",
		},
		Attach:      req.Attach,
		Description: req.Description,
		OutTradeNo:  req.OutTradeNo,
		Payer: &request.JSAPIPayer{
			OpenID: req.Openid,
		},
	}

	if req.NotifyUrl != "" {
		options.NotifyUrl = req.NotifyUrl
	}
	// 下单
	resp, err := wx.PaymentApp.Order.JSAPITransaction(context.Background(), options)
	if err != nil {
		return nil, err
	}
	if resp.PrepayID == "" {
		return nil, errors.New("get prepayId err")
	}

	payConf, err := wx.PaymentApp.JSSDK.BridgeConfig(resp.PrepayID, true)
	if err != nil {
		return nil, err
	}
	base64Str, _ := payConf.([]byte)
	var data WxPayResp
	err = json.Unmarshal(base64Str, &data)
	if err != nil {
		return nil, err
	}
	data.OutTradeNo = req.OutTradeNo
	return &data, nil
}

type RefundRequest struct {
	OrderId    string `json:"order_id,omitempty"`
	TotalFee   int    `json:"total_fee,omitempty"`
	RefundFee  int    `json:"refund_fee,omitempty"`
	RefundDesc string `json:"refund_desc,omitempty"`
	NotifyUrl  string `json:"notify_url,omitempty"`
}

type RefundResp struct {
	Code        int    `json:"code"`
	OutRefundNo string `json:"out_refund_no"`
	Msg         string `json:"msg"`
}

// Refund 退款
func (wx *WXPay) Refund(req *RefundRequest) (*RefundResp, error) {
	outRefundNo := fmt.Sprintf("%s@%d", req.OrderId, gb.Now().Unix())
	options := &rRequest.RequestRefund{
		OutTradeNo:   req.OrderId,
		OutRefundNo:  outRefundNo,
		Reason:       req.RefundDesc,
		FundsAccount: "",
		Amount: &rRequest.RefundAmount{
			Refund:   req.RefundFee,                  // 退款金额，单位：分
			Total:    req.TotalFee,                   // 订单总金额，单位：分
			From:     []*rRequest.RefundAmountFrom{}, // 退款出资账户及金额。不传仍然需要这个空数组防止微信报错
			Currency: "CNY",
		},
		GoodsDetail: nil,
	}
	if req.NotifyUrl != "" {
		options.NotifyUrl = req.NotifyUrl
	}
	refund, err := wx.PaymentApp.Refund.Refund(context.Background(), options)
	if err != nil {
		return nil, err
	}
	if refund.TransactionID == "" {
		fmt.Println(refund, err)
		return nil, errors.New(fmt.Sprintf("get transactionID err:%s", refund.Message))
	}
	return &RefundResp{
		Code:        0,
		OutRefundNo: outRefundNo,
		Msg:         "SUCCESS",
	}, nil
}

func (wx *WXPay) payNotify(r *http.Request, f func(orderId, attach string) error) (*http.Response, error) {
	res, err := wx.PaymentApp.HandlePaidNotify(
		r,
		func(message *nRequest.RequestNotify, transaction *models.Transaction, fail func(message string)) interface{} {
			// 看下支付通知事件状态
			// 这里可能是微信支付失败的通知，所以可能需要在数据库做一些记录，然后告诉微信我处理完成了。
			if message.EventType != "TRANSACTION.SUCCESS" {
				return true
			}
			if transaction.OutTradeNo != "" {
				// 这里对照自有数据库里面的订单做查询以及支付状态改变
				log.Printf("订单号：%s 支付成功", transaction.OutTradeNo)
				if err := f(transaction.OutTradeNo, transaction.Attach); err != nil {
					fail("payment fail")
					return false
				}
			} else {
				// 因为微信这个回调不存在订单号，所以可以告诉微信我还没处理成功，等会它会重新发起通知
				// 如果不需要，直接返回true即可
				fail("payment fail")
				return true
			}
			return true
		},
	)

	// 这里可能是因为不是微信官方调用的，无法正常解析出transaction和message，所以直接抛错。
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (wx *WXPay) refundNotify(r *http.Request, f func(orderId string) error) (*http.Response, error) {
	res, err := wx.PaymentApp.HandleRefundedNotify(
		r,
		func(message *nRequest.RequestNotify, transaction *models.Refund, fail func(message string)) interface{} {
			// 看下支付通知事件状态
			// 这里可能是微信支付失败的通知，所以可能需要在数据库做一些记录，然后告诉微信我处理完成了。
			if message.EventType != "REFUND.SUCCESS" {
				return true
			}
			if transaction.OutTradeNo != "" {
				// 这里对照自有数据库里面的订单做查询以及支付状态改变
				log.Printf("订单号：%s 退款成功", transaction.OutRefundNo)
				if err := f(transaction.OutRefundNo); err != nil {
					return false
				}
			} else {
				// 因为微信这个回调不存在订单号，所以可以告诉微信我还没处理成功，等会它会重新发起通知
				// 如果不需要，直接返回true即可
				fail("refund fail")
				return true
			}
			return true
		},
	)

	// 这里可能是因为不是微信官方调用的，无法正常解析出transaction和message，所以直接抛错。
	if err != nil {
		return nil, err
	}
	return res, nil
}

// QueryOrder 查询支付订单
func (wx *WXPay) QueryOrder(orderId string) (*response.ResponseOrder, error) {
	order, err := wx.PaymentApp.Order.QueryByOutTradeNumber(context.Background(), orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// QueryRefundOrder 查询退款订单
func (wx *WXPay) QueryRefundOrder(orderId string) (*rResponse.ResponseRefund, error) {
	order, err := wx.PaymentApp.Refund.Query(context.Background(), orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
