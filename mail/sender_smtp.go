// Package email implementa utilidades para el envío de correos electrónicos.
package email

import (
	"context"

	"gopkg.in/gomail.v2"
)

type SMTPSender struct {
	dialer      *gomail.Dialer
	senderEmail string
}

func (s *SMTPSender) Send(ctx context.Context, from, to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	return s.dialer.DialAndSend(msg)
}
