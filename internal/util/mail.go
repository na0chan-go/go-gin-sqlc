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

// Mailer はメール送信のインターフェースです
type Mailer interface {
	SendMail(to, subject, body string) error
}

// SMTPMailer はSMTPを使用したメール送信の実装です
type SMTPMailer struct {
	config MailConfig
}

// NewSMTPMailer は新しいSMTPMailerを作成します
func NewSMTPMailer(config MailConfig) *SMTPMailer {
	return &SMTPMailer{
		config: config,
	}
}

// SendMail はメールを送信します
func (m *SMTPMailer) SendMail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", m.config.From, to, subject, body))

	return smtp.SendMail(addr, auth, m.config.From, []string{to}, msg)
}

// GeneratePasswordResetEmail はパスワードリセットメールの本文を生成します
func GeneratePasswordResetEmail(resetURL string) string {
	return fmt.Sprintf(`パスワードリセットのリクエストを受け付けました。

以下のURLをクリックしてパスワードをリセットしてください：
%s

このリンクは24時間有効です。
心当たりがない場合は、このメールを無視してください。`, resetURL)
}
