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

// SendMail はメールを送信します
func SendMail(config MailConfig, to string, subject string, body string) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", config.From, to, subject, body)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return smtp.SendMail(addr, auth, config.From, []string{to}, []byte(msg))
}

// GeneratePasswordResetEmail はパスワードリセットメールの本文を生成します
func GeneratePasswordResetEmail(resetURL string) string {
	return fmt.Sprintf(`
		<h2>パスワードリセットのリクエスト</h2>
		<p>パスワードリセットのリクエストを受け付けました。</p>
		<p>以下のリンクをクリックしてパスワードを再設定してください：</p>
		<p><a href="%s">パスワードを再設定する</a></p>
		<p>このリンクは24時間有効です。</p>
		<p>このメールに心当たりがない場合は、無視してください。</p>
	`, resetURL)
}
