package server

import (
	"github.com/mojocn/base64Captcha"
	"github.com/zhangxa/gfcore/core"
)

type sCaptcha struct{}

var (
	captchaStore  = base64Captcha.DefaultMemStore
	captchaDriver = newDriver()
)

func init() {
	core.RegisterCaptcha(New())
}

// New 验证码管理服务
func New() core.ICaptcha {
	return &sCaptcha{}
}

func newDriver() *base64Captcha.DriverString {
	driver := &base64Captcha.DriverString{
		Height:          44,
		Width:           126,
		NoiseCount:      1,
		ShowLineOptions: base64Captcha.OptionShowSineLine,
		Length:          4,
		Source:          "1234567890",
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	return driver.ConvertFonts()
}

// NewAndStore 创建验证码，直接输出验证码图片内容到HTTP Response.
func (s *sCaptcha) NewAndStore() (id string, img string, err error) {
	captcha := base64Captcha.NewCaptcha(captchaDriver, captchaStore)
	id, img, _, err = captcha.Generate()
	return
}

// VerifyAndClear 校验验证码，并清空缓存的验证码信息
func (s *sCaptcha) VerifyAndClear(id string, value string) bool {
	return captchaStore.Verify(id, value, true)
}
