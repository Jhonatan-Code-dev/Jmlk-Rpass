// internal/adapters/smtp/gomail.go
package smtpadapter

import (
	"context"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/ports"
	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	dialer *gomail.Dialer
	from   string
}

func NewGomailSender(host string, port int, user, pass string) ports.SMTPClient {
	return &GomailSender{
		dialer: gomail.NewDialer(host, port, user, pass),
		from:   user,
	}
}

func (s *GomailSender) Send(ctx context.Context, from, to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)
	return s.dialer.DialAndSend(msg)
}
