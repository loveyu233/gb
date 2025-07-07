package gb

import "github.com/loveyu233/gb/captcha"

const (
	TypeClick  captcha.CaptchaType = "click"
	TypeRotate captcha.CaptchaType = "rotate"
	TypeSlide  captcha.CaptchaType = "slide"
)

var CaptchaManager *captcha.Manager

func InitCaptchaManager(cache captcha.CacheImpl) {
	CaptchaManager = captcha.NewCaptchaManager(cache)
}
