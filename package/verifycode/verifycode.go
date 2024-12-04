package verifycode

import (
	"NothingBlog/settings"

	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore //base64Captcha.Store

func Init(cfg *settings.AuthConfig) {
	//store = base64Captcha.NewMemoryStore(cfg.CodeNum, time.Duration(cfg.ExpiredTime*int64(time.Second)))
}

// 使用数字验证码 还可以使用语音，数字等验证码
func GenerateVerifyCode() (id string, b64s string, ans string, err error) {
	//var driver base64Captcha.Driver

	//var ds = new(base64Captcha.DriverString)

	driver := base64Captcha.DefaultDriverDigit //ds.ConvertFonts()
	id, b64s, ans, err = base64Captcha.NewCaptcha(driver, store).Generate()
	return
}

// 检查验证码
func CheckVerifyCode(id string, code string) bool {
	return store.Verify(id, code, false)
}
