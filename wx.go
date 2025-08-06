package gb

var InsWX = &wxClient{
	WXPay:      &wxPay{},
	WXMini:     &wxMini{},
	WXOfficial: &wxOfficial{},
}

type wxClient struct {
	WXPay      *wxPay
	WXMini     *wxMini
	WXOfficial *wxOfficial
}
