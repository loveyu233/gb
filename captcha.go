package gb

import "github.com/loveyu233/gb/captcha"

const (
	TypeClick  captcha.CaptchaType = "click"
	TypeRotate captcha.CaptchaType = "rotate"
	TypeSlide  captcha.CaptchaType = "slide"
)

func NewCaptchaManager(cache captcha.CacheImpl) *captcha.Manager {
	return captcha.NewCaptchaManager(cache)
}
