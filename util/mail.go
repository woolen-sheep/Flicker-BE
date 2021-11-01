package util

import (
	"github.com/labstack/gommon/random"
	"github.com/woolen-sheep/Flicker-BE/config"

	"gopkg.in/gomail.v2"
)

// GenerateVerifyCode 生成五位数验证码
func GenerateVerifyCode() string {
	return random.String(5, random.Numeric)
}

func SendMailVerifyCode(mail, code string) error {
	return SendMail(mail, "Flicker验证码", code)
}

// SendMail through qq SMTP service.
func SendMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader(`From`, config.C.Mail.Address)
	m.SetHeader(`To`, to)
	m.SetHeader(`Subject`, subject)
	m.SetBody(`text/html`, body)
	err := gomail.
		NewDialer(config.C.Mail.Host, config.C.Mail.Port, config.C.Mail.Address, config.C.Mail.Password).
		DialAndSend(m)
	return err
}
