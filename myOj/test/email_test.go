package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Jordan Wright <996800442@qq.com>"
	e.To = []string{"test@example.com"}
	e.Subject = "验证码发送测试"
	e.Text = []byte("您的验证码是: <b>123456</b>")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))
}
