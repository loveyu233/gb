package pay

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
	"github.com/gin-gonic/gin"
)

type WXPay struct {
	PaymentApp          *payment.Payment
	payNotifyHandler    func(orderId string, attach string) error
	refundNotifyHandler func(orderId string) error
	payHandler          func(c *gin.Context) (*PayRequest, error)
	refundHandler       func(c *gin.Context) (*RefundRequest, error)
	// 是否保存请求日志
	IsSaveHandlerLog bool
}
type WXPayImp interface {
	PayNotify(orderId string, attach string) error
	RefundNotify(orderId string) error

	Pay(c *gin.Context) (*PayRequest, error)
	Refund(c *gin.Context) (*RefundRequest, error)
}
type Payment struct {
	AppID              string                `json:"appID,omitempty"`              // 小程序、公众号或者企业微信的appid
	MchID              string                `json:"mchID,omitempty"`              // 商户号 appID
	MchApiV3Key        string                `json:"mchApiV3Key,omitempty"`        // 微信V3接口调用必填
	Key                string                `json:"key,omitempty"`                // 微信V2接口调用必填
	CertPath           string                `json:"certPath,omitempty"`           // 商户后台支付的Cert证书路径
	KeyPath            string                `json:"keyPath,omitempty"`            // 商户后台支付的Key证书路径
	SerialNo           string                `json:"serialNo,omitempty"`           // 商户支付证书序列号
	WechatPaySerial    string                `json:"wechatPaySerial,omitempty"`    // 微信支付平台证书序列号[选填]
	CertificateKeyPath string                `json:"certificateKeyPath,omitempty"` // 微信支付平台证书路径，[选填]
	RSAPublicKeyPath   string                `json:"RSAPublicKeyPath,omitempty"`   // 微信支付平台证书的Key证书路径[选填]
	SubMchID           string                `json:"subMchID,omitempty"`           // 服务商平台下的子商户号Id，[选填]
	SubAppID           string                `json:"subAppID,omitempty"`           // 服务商平台下的子AppId，[选填]
	NotifyURL          string                `json:"notifyURL,omitempty"`          // 微信支付完成后的通知回调地址
	HttpDebug          bool                  `json:"httpDebug,omitempty"`          // 是否开启打印 SDK 调用微信 API 接口时候的日志
	Log                payment.Log           `json:"log"`                          // 可以重定向到你的目录下，如果设置File和Error，默认会在当前目录下的wechat文件夹下生成日志
	Http               payment.Http          `json:"http"`                         // 设置微信支付地址，比如想要设置成沙盒地址，把里面的值改成https://api.mch.weixin.qq.com/sandboxnew
	Cache              kernel.CacheInterface `json:"cache,omitempty"`              // 可选，不传默认走程序内存
	// 是否保存请求日志
	IsSaveHandlerLog bool
}

type WXPaymentAppConfig struct {
	Payment  Payment
	WXPayImp WXPayImp
}

func InitWXWXPaymentApp(paymentConfig WXPaymentAppConfig) (*WXPay, error) {
	paymentApp, err := payment.NewPayment(&payment.UserConfig{
		AppID:              paymentConfig.Payment.AppID,
		MchID:              paymentConfig.Payment.MchID,
		MchApiV3Key:        paymentConfig.Payment.MchApiV3Key,
		Key:                paymentConfig.Payment.Key,
		CertPath:           paymentConfig.Payment.CertPath,
		KeyPath:            paymentConfig.Payment.KeyPath,
		SerialNo:           paymentConfig.Payment.SerialNo,
		CertificateKeyPath: paymentConfig.Payment.CertificateKeyPath,
		WechatPaySerial:    paymentConfig.Payment.WechatPaySerial,
		RSAPublicKeyPath:   paymentConfig.Payment.RSAPublicKeyPath,
		SubAppID:           paymentConfig.Payment.SubAppID,
		SubMchID:           paymentConfig.Payment.SubMchID,
		Http:               paymentConfig.Payment.Http,
		ResponseType:       response.TYPE_MAP,
		Log:                paymentConfig.Payment.Log,
		Cache:              paymentConfig.Payment.Cache,
		HttpDebug:          paymentConfig.Payment.HttpDebug,
		NotifyURL:          paymentConfig.Payment.NotifyURL,
	})
	if err != nil {
		return nil, err
	}
	return &WXPay{
		PaymentApp:          paymentApp,
		payNotifyHandler:    paymentConfig.WXPayImp.PayNotify,
		refundNotifyHandler: paymentConfig.WXPayImp.RefundNotify,
		payHandler:          paymentConfig.WXPayImp.Pay,
		refundHandler:       paymentConfig.WXPayImp.Refund,
		IsSaveHandlerLog:    paymentConfig.Payment.IsSaveHandlerLog,
	}, nil
}
