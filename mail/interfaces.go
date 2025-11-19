package mail

import (
	"context"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/models"
)

type SMTPClient interface {
	Send(ctx context.Context, from, to, subject, htmlBody string) error
}

type Repository interface {
	SaveCode(ctx context.Context, entry models.CodeEntry) error
	GetCodeEntry(ctx context.Context, email string) (*models.CodeEntry, error)
	Close() error
}
