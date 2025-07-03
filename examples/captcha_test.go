package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestR(t *testing.T) {
	client := gb.InitCaptchaClient(gb.WithCaptchaRotateCapt(10))

	capt, err := client.RotateCapt("123")
	if err != nil {
		t.Fatal(err)
		return
	}

	client.RotateCaptVerify(capt.Key, 10)
}
