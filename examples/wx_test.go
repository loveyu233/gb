package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

type PayC struct {
}

func (p PayC) PaySuccess(orderId string, attach string) error {
	//TODO implement me
	panic("implement me")
}

func (p PayC) RefundSuccess(orderId string) error {
	//TODO implement me
	panic("implement me")
}

func TestWXPay(t *testing.T) {
	gb.WXNewWXPaymentApp(gb.WXPaymentAppConfig{
		Payment:  gb.Payment{},
		WXPayImp: PayC{},
	})
}

func TestWXMini(t *testing.T) {

}

func TestOffia(t *testing.T) {

}
