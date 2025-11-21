package ports

import (
	"context"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/domain"
)

type Repository interface {
	SaveCode(ctx context.Context, entry domain.CodeEntry) error
	GetCodeEntry(ctx context.Context, email string) (*domain.CodeEntry, error)
	Close() error
}
