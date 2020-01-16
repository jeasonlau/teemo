package mail

import (
	"net/smtp"
)

func SendMail(email, content string) error { //from与password需要修改
	from := ""
	password := ""
	auth := smtp.PlainAuth(
		"",
		from,
		password,
		"smtp.qq.com",
	)
	to := []string{email}
	body := []byte("To: " + to[0] + "\r\nFrom: PUSH\r\nSubject: 绩点变化！\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + content)
	err := smtp.SendMail(
		"smtp.qq.com:25",
		auth,
		from,
		to,
		body,
	)
	return err
}
