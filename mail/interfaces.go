package email

import "context"

// SMTPClient abstrae el envío de correos (permite mock en tests).
type SMTPClient interface {
	Send(ctx context.Context, from string, to string, subject string, htmlBody string) error
}

// Repository abstrae la persistencia de códigos (permite cambiar BoltDB por otra DB).
type Repository interface {
	SaveCode(ctx context.Context, entry CodeEntry) error
	GetCodeEntry(ctx context.Context, email string) (*CodeEntry, error)
	Close() error
}
