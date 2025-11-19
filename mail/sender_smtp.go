package mail

import (
	"context"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"gopkg.in/gomail.v2"
)

// SMTPSender implementa el envío de correos.
type SMTPSender struct {
	Dialer      *gomail.Dialer
	SenderEmail string
}

// NewGomailDialer construye un Dialer a partir de la configuración.
func NewGomailDialer(cfg config.EmailConfig) *gomail.Dialer {
	return gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
}

// Send envía un correo HTML.
func (s *SMTPSender) Send(ctx context.Context, from, to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)
	return s.Dialer.DialAndSend(msg)
}
