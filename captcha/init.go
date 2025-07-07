package captcha

var CaptchaManager *Manager

func NewCaptchaManager(cache CacheImpl) *Manager {
	CaptchaManager = NewManager(cache)

	// 注册验证码生成器
	CaptchaManager.RegisterGenerator(NewClickCaptcha(CaptchaManager.cache))
	CaptchaManager.RegisterGenerator(NewRotateCaptcha(CaptchaManager.cache))
	CaptchaManager.RegisterGenerator(NewSlideCaptcha(CaptchaManager.cache))
	return CaptchaManager
}
