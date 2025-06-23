package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestR(t *testing.T) {
	gb.InitCaptchaClient(gb.WithCaptchaRotateCapt(10))

	capt, err := gb.Captcha.RotateCapt("123")
	if err != nil {
		t.Fatal(err)
		return
	}

	gb.Captcha.RotateCaptVerify(capt.Key, 10)
}
