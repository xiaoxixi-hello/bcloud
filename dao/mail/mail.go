package mail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

var e = &email.Email{
	From:    "2025907338@qq.com",
	Subject: "Bcloud任务提交进程通知",
}

var auth = smtp.PlainAuth("", "2025907338@qq.com", "gwuwqyzgcwyzdbdj", "smtp.qq.com")

func SendToMe(text string) error {
	e.To = []string{"2025907338@qq.com"}
	e.Text = []byte(text)
	//return nil
	return e.Send("smtp.qq.com:25", auth)
}
