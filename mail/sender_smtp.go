package email

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

// SMTPSender implementa SMTPClient usando gomail.
type SMTPSender struct {
	dialer      *gomail.Dialer
	senderEmail string
}

func (s *SMTPSender) Send(ctx context.Context, from string, to string, subject string, htmlBody string) error {
	// gomail no soporta context nativo; en entornos reales se podría usar net/smtp o una librería con contexto.
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	// gomail.DialAndSend bloquea; si necesitas timeout/cancelación, ejecuta en goroutine y usa canal/select con ctx.Done().
	done := make(chan error, 1)
	go func() {
		done <- s.dialer.DialAndSend(msg)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("envío cancelado o timeout: %w", ctx.Err())
	}
}
