package gb

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	alipayv2 "github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/alipay/v3"
	"net/http"
	"time"
)

// ZfbMiniImp 定义用户相关操作接口
type ZfbMiniImp interface {
	// IsExistsUser 检查用户是否存在
	IsExistsUser(unionID string) (user any, exists bool, err error)

	// CreateUser 创建新用户
	CreateUser(phoneNumber, unionID, openID, areaCodeByIP, clientIP string) (user any, err error)

	// GenerateToken 生成用户token
	GenerateToken(user any, sessionKey string) (data any, err error)

	// Pay 支付
	Pay(c *gin.Context) (*ZFBPayParam, error)

	Refund(c *gin.Context) (*ZFBRefundParam, error)

	// PayNotify 支付成功的异步通知
	PayNotify(*ZFBPay)

	// RefundNotify 支付退款的异步通知
	RefundNotify(*ZFBRefund)
}

type ZFBClient struct {
	appid         string // 应用的appid
	appPrivateKey string // 应用私钥
	aesKey        string // 接口加密aes密钥
	appPublicKey  []byte // 应用公钥
	aliPublicKey  []byte // 阿里公钥
	aliRootKey    []byte // 阿里根证书
	client        *alipay.ClientV3
	notifyUrl     string
	zfbMiniImp    ZfbMiniImp
	// 是否保存请求日志
	IsSaveHandlerLog bool
}

var ZFB = new(ZFBClient)

// InitAliClient 其中appid,appPrivateKey,aesKey是内容本身,appPublicKey, aliPublicKey, aliRootKey是证书路径,notifyUrl为支付成功和退款异步通知地址,isSaveHandlerLog 是否保存请求日志
func InitAliClient(appid, appPrivateKey, aesKey, appPublicKeyFilePath, aliPublicKeyFilePath, aliRootKeyFilePath, notifyUrl string, isSaveHandlerLog bool, zfbMiniImp ZfbMiniImp) error {
	appPublicKey, err := ReadFileContent(appPublicKeyFilePath)
	if err != nil {
		return err
	}
	aliPublicKey, err := ReadFileContent(aliPublicKeyFilePath)
	if err != nil {
		return err
	}
	aliRootKey, err := ReadFileContent(aliRootKeyFilePath)
	if err != nil {
		return err
	}

	ZFB.appid = appid
	ZFB.appPrivateKey = appPrivateKey
	ZFB.aesKey = aesKey
	ZFB.appPublicKey = appPublicKey
	ZFB.aliPublicKey = aliPublicKey
	ZFB.aliRootKey = aliRootKey
	ZFB.notifyUrl = notifyUrl
	ZFB.IsSaveHandlerLog = isSaveHandlerLog
	ZFB.zfbMiniImp = zfbMiniImp

	clientV3, err := alipay.NewClientV3(appid, appPrivateKey, true)
	if err != nil {
		return err
	}
	clientV3.SetAESKey(aesKey)
	if err = clientV3.SetCert(appPublicKey, aliPublicKey, aliRootKey); err != nil {
		return err
	}

	ZFB.client = clientV3
	return nil
}

type ZFBPayParam struct {
	OutTradeNo  string
	BuyerOpenID string
	Subject     string
	TotalAmount float64
}

// TradeCreate 支付
func (a *ZFBClient) TradeCreate(param *ZFBPayParam) (aliRsp *alipay.TradeCreateRsp, err error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("total_amount", param.TotalAmount).
		Set("product_code", "JSAPI_PAY").
		Set("op_app_id", a.appid).
		Set("buyer_open_id", param.BuyerOpenID).
		Set("notify_url", a.notifyUrl).
		Set("subject", param.Subject)

	return a.client.TradeCreate(context.Background(), bm)
}

// TradeQuery 支付结果查询
func (a *ZFBClient) TradeQuery(outTradeNo string) (*alipay.TradeQueryRsp, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo)
	return a.client.TradeQuery(context.Background(), bm)
}

type ZFBRefundParam struct {
	OutTradeNo   string
	RefundReason string
	RefundAmount float64
}

// TradeRefund 发起退款
func (a *ZFBClient) TradeRefund(param *ZFBRefundParam) (*alipay.TradeRefundRsp, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", param.OutTradeNo).
		Set("refund_amount", param.RefundAmount).
		Set("refund_reason", param.RefundReason)

	return a.client.TradeRefund(context.Background(), bm)
}

// TradeFastPayRefundQuery 退款结果查询
func (a *ZFBClient) TradeFastPayRefundQuery(outTradeNo, outRequestNo string) (*alipay.TradeFastPayRefundQueryRsp, error) {
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", outTradeNo).
		Set("out_request_no", outRequestNo)

	return a.client.TradeFastPayRefundQuery(context.Background(), bm)
}

// SystemOauthToken 获取用户code
func (a *ZFBClient) SystemOauthToken(code string) (*alipay.SystemOauthTokenRsp, error) {
	bodyMap := make(gopay.BodyMap)
	body := bodyMap.Set("grant_type", "authorization_code").
		Set("code", code)

	return a.client.SystemOauthToken(context.Background(), body)
}

// UserInfoShare 用code换取用户信息
func (a *ZFBClient) UserInfoShare(authToken string) (*alipay.UserInfoShareRsp, error) {
	bodyMap := make(gopay.BodyMap)
	body := bodyMap.Set("auth_token", authToken)

	return a.client.UserInfoShare(context.Background(), body)
}

type MobilePhoneNumberDecryptionResp struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"subCode"`
	SubMsg  string `json:"subMsg"`
	Mobile  string `json:"mobile"`
}

// MobilePhoneNumberDecryption 解密用户手机号
func (a *ZFBClient) MobilePhoneNumberDecryption(response string) (*MobilePhoneNumberDecryptionResp, error) {
	if a.aesKey == "" {
		return nil, fmt.Errorf("aes key is empty")
	}
	var data, err = base64.StdEncoding.DecodeString(a.aesKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(data)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, 16)
	for i := 0; i < 16; i++ {
		iv[i] = 0
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(response)
	if err != nil {
		return nil, err
	}

	if len(encryptedBytes)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptedBytes, encryptedBytes)

	decryptedBytes, err := a.pkcs5Unpad(encryptedBytes)
	if err != nil {
		return nil, err
	}

	var resp = new(MobilePhoneNumberDecryptionResp)
	if err := json.Unmarshal(decryptedBytes, resp); err != nil {
		return nil, err
	}

	if resp.Code != "10000" {
		return resp, fmt.Errorf("%s:%s", resp.Msg, resp.SubMsg)
	}

	return resp, nil
}

func (a *ZFBClient) pkcs5Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpadding error")
	}

	return src[:(length - unpadding)], nil
}

func (a *ZFBClient) RegisterHandlers(r *gin.RouterGroup) {
	r.Use(SetModuleName("支付宝"))
	r.POST("/zfb/login", SetOptionName("小程序登录", a.IsSaveHandlerLog), a.login)
	r.POST("/zfb/notify", SetOptionName("支付异步回调", a.IsSaveHandlerLog), a.notify)
	r.POST("/zfb/pay", SetOptionName("支付请求", a.IsSaveHandlerLog), a.pay)
	r.POST("/zfb/refund", SetOptionName("退款请求", a.IsSaveHandlerLog), a.refund)
}

func (a *ZFBClient) login(c *gin.Context) {
	var params struct {
		Code          string `binding:"required" json:"code"`
		EncryptedData string `json:"encrypted_data"`
		IvStr         string `json:"iv_str"`
	}
	if err := c.BindJSON(&params); err != nil {
		ResponseError(c, ErrInvalidParam)
		return
	}

	session, err := a.SystemOauthToken(params.Code)
	if err != nil || session.OpenId == "" {
		ResponseError(c, ErrRequestAli.WithMessage("获取支付宝小程序用户会话代码失败"))
		return
	}

	var (
		user   any
		exists bool
	)

	//检测用户是否注册
	user, exists, err = a.zfbMiniImp.IsExistsUser(session.OpenId)
	if err != nil {
		ResponseError(c, ErrDatabase.WithMessage("查询用户信息失败:%s", err.Error()))
		return
	}
	if !exists {
		if params.EncryptedData == "" {
			ResponseSuccess(c, map[string]interface{}{
				"open_id": session.OpenId,
			})
			return
		}
		decryption, err := a.MobilePhoneNumberDecryption(params.EncryptedData)
		if err != nil {
			ResponseError(c, ErrRequestAli.WithMessage("获取支付宝小程序用户数据失败"))
		}

		if user, err = a.zfbMiniImp.CreateUser(decryption.Mobile, session.UnionId, session.OpenId, getAreaCodeByIp(c.ClientIP()), c.ClientIP()); err != nil {
			ResponseError(c, ErrDatabase.WithMessage("创建用户信息失败:%s", err.Error()))
			return
		}
	}

	data, err := a.zfbMiniImp.GenerateToken(user, session.OpenId)
	if err != nil {
		ResponseError(c, ErrServerBusy.WithMessage("token生成失败:%s", err.Error()))
		return
	}
	ResponseSuccess(c, data)
}

type ZFBSyncNotify struct {
	AppId          string    `json:"app_id"`
	NotifyTime     time.Time `json:"notify_time"`      // 发送通知时间
	TradeNo        string    `json:"trade_no"`         // 支付宝交易号
	OutTradeNo     string    `json:"out_trade_no"`     // 商户订单号。
	OutBizNo       string    `json:"out_biz_no"`       // 商家业务号。商家业务 ID，主要是退款通知中返回退款申请的流水号。
	BuyerOpenId    string    `json:"buyer_open_id"`    // 买家支付宝用户号
	BuyerLogonId   string    `json:"buyer_logon_id"`   // 买家支付宝账号
	TradeStatus    string    `json:"trade_status"`     // 交易状态。咨询目前所处的状态。
	TotalAmount    float64   `json:"total_amount"`     // 订单金额。本次交易支付的订单金额，单位为人民币（元）。支持小数点后两位。
	ReceiptAmount  float64   `json:"receipt_amount"`   // 实收金额。商家在交易中实际收到的款项，单位为人民币（元）。支持小数点后两位。
	BuyerPayAmount float64   `json:"buyer_pay_amount"` // 付款金额。用户在交易中支付的金额。支持小数点后两位。
	PointAmount    float64   `json:"point_amount"`     // 集分宝金额。使用集分宝支付的金额。支持小数点后两位。
	RefundFee      float64   `json:"refund_fee"`       // 总退款金额。退款通知中，返回总退款金额，单位为元，支持小数点后两位。
	SendBackFee    float64   `json:"send_back_fee"`    // 实际退款金额。商家实际退款给用户的金额，单位为元，支持小数点后两位。
	Subject        string    `json:"subject"`          // 订单标题。商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来。
	PassbackParams string    `json:"passback_params"`  // 公共回传参数，如果请求时传递了该参数，则返回给商家时会在异步通知时将该参数原样返回。本参数必须进行 UrlEncode 之后才可以发送给支付宝。
	GmtCreate      time.Time `json:"gmt_create"`       // 交易创建时间。该笔交易创建的时间。格式 为 yyyy-MM-dd HH:mm:ss。
	GmtPayment     time.Time `json:"gmt_payment"`      // 交易 付款时间。该笔交易的买家付款时间。格式为 yyyy-MM-dd HH:mm:ss。
	GmtRefund      time.Time `json:"gmt_refund"`       // 交易退款时间。该笔交易的退款时间。格式 为 yyyy-MM-dd HH:mm:ss.SS。
	GmtClose       time.Time `json:"gmt_close"`        // 交易结束时间。该笔交易结束时间。格式为 yyyy-MM-dd HH:mm:ss
}

type ZFBRefund struct {
	RefundFee   float64 `json:"refund_fee"`    // 总退款金额。退款通知中，返回总退款金额，单位为元，支持小数点后两位。
	SendBackFee float64 `json:"send_back_fee"` // 实际退款金额。商家实际退款给用户的金额，单位为元，支持小数点后两位。
}
type ZFBPay struct {
	TotalAmount    float64 `json:"total_amount"`     // 订单金额。本次交易支付的订单金额，单位为人民币（元）。支持小数点后两位。
	ReceiptAmount  float64 `json:"receipt_amount"`   // 实收金额。商家在交易中实际收到的款项，单位为人民币（元）。支持小数点后两位。
	BuyerPayAmount float64 `json:"buyer_pay_amount"` // 付款金额。用户在交易中支付的金额。支持小数点后两位。
}

func (a *ZFBClient) zfbPayNotify(req *http.Request) (*ZFBSyncNotify, error) {
	bodyMap, err := alipayv2.ParseNotifyToBodyMap(req)
	if err != nil {
		return nil, err
	}
	ok, err := alipayv2.VerifySign(string(a.appPublicKey), bodyMap)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("sign err")
	}

	var resp = new(ZFBSyncNotify)
	if err := bodyMap.Unmarshal(resp); err != nil {
		return nil, err
	}

	if resp.RefundFee != 0 {
		// 退款
		a.zfbMiniImp.RefundNotify(&ZFBRefund{
			RefundFee:   resp.RefundFee,
			SendBackFee: resp.SendBackFee,
		})
	} else if resp.ReceiptAmount != 0 {
		// 支付
		a.zfbMiniImp.PayNotify(&ZFBPay{
			TotalAmount:    resp.TotalAmount,
			ReceiptAmount:  resp.ReceiptAmount,
			BuyerPayAmount: resp.BuyerPayAmount,
		})
	}

	return resp, nil
}

func (a *ZFBClient) notify(c *gin.Context) {
	_, err := a.zfbPayNotify(c.Request)
	if err != nil {
		c.Writer.WriteString("fail")
		return
	}
	c.Writer.WriteString("success")
}

func (a *ZFBClient) pay(c *gin.Context) {
	payParam, err := a.zfbMiniImp.Pay(c)
	if err != nil {
		ResponseError(c, ErrRequestAliPay.WithMessage(err.Error()))
		return
	}
	aliRsp, err := a.TradeCreate(payParam)
	if err != nil {
		ResponseError(c, ErrRequestAliPay.WithMessage(err.Error()))
		return
	}
	if aliRsp.StatusCode != 10000 {
		ResponseError(c, ErrRequestAliPay.WithMessage(aliRsp.ErrResponse.Message))
		return
	}
	ResponseSuccess(c, gin.H{
		"trade_no":     aliRsp.TradeNo,
		"out_trade_no": aliRsp.OutTradeNo,
	})
}

func (a *ZFBClient) refund(c *gin.Context) {
	param, err := a.zfbMiniImp.Refund(c)
	if err != nil {
		ResponseError(c, ErrRequestAliPay.WithMessage(err.Error()))
		return
	}
	refund, err := a.TradeRefund(param)
	if err != nil {
		ResponseError(c, ErrRequestAliPay.WithMessage(err.Error()))
		return
	}
	if refund.StatusCode != 10000 {
		ResponseError(c, ErrRequestAliPay.WithMessage(refund.ErrResponse.Message))
		return
	}
	ResponseSuccess(c, gin.H{
		"trade_no":     refund.TradeNo,
		"out_trade_no": refund.OutTradeNo,
		"refund_fee":   refund.RefundFee,
	})
}
