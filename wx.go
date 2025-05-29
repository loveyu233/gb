package gb

var WX = &WXClient{WXPay: &WXPay{}, WXMini: &WXMini{}, WXOfficial: &WXOfficial{}}

type WXClient struct {
	WXPay      *WXPay
	WXMini     *WXMini
	WXOfficial *WXOfficial
}
