package util

import (
	"fmt"
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// MailSender はメール送信のインターフェースです
type MailSender func(config MailConfig, to, subject, body string) error

// SendMail はデフォルトのメール送信関数です
var SendMail MailSender = func(config MailConfig, to, subject, body string) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", config.From, to, subject, body))

	return smtp.SendMail(addr, auth, config.From, []string{to}, msg)
}

// GeneratePasswordResetEmail はパスワードリセットメールの本文を生成します
func GeneratePasswordResetEmail(resetURL string) string {
	return fmt.Sprintf(`パスワードリセットのリクエストを受け付けました。

以下のURLをクリックしてパスワードをリセットしてください：
%s

このリンクは24時間有効です。
心当たりがない場合は、このメールを無視してください。`, resetURL)
}
