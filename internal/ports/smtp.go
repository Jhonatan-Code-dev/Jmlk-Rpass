// internal/ports/smtp.go
package ports

import "context"

type SMTPClient interface {
	Send(ctx context.Context, from, to, subject, htmlBody string) error
}
