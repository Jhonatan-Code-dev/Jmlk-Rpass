package email

import (
	"context"
	"fmt"
	"time"
)

// CanSendReset aplica las políticas de envío y devuelve (puedeEnviar, mensaje).
func (s *EmailService) CanSendReset(ctx context.Context, email string) (bool, string, error) {
	entry, err := s.repo.GetCodeEntry(ctx, email)
	if err != nil {
		// Si el repositorio devuelve "no encontrado", interpretamos como nil y permitimos envío.
		// Repository debe devolver error si hay problema de I/O. Aquí tratamos both.
		// Para distinguir, GetCodeEntry implementa error específico.
		entry = nil
	}

	now := time.Now()

	if entry == nil {
		return true, "Primer envío permitido.", nil
	}

	if entry.Attempts >= s.conf.MaxResetAttempts {
		return false, fmt.Sprintf(
			"Máximo de intentos (%d). Espera %.0f horas.",
			s.conf.MaxResetAttempts,
			s.conf.RestrictionPeriod.Hours(),
		), nil
	}

	if s.conf.AllowOverride {
		return true, "Override activo → se generará nuevo código.", nil
	}

	if entry.Used {
		return false, fmt.Sprintf(
			"Ya usaste tu último código. Espera %.0f horas.",
			s.conf.RestrictionPeriod.Hours(),
		), nil
	}

	if now.Before(entry.ExpireAt) {
		return false, fmt.Sprintf(
			"Aún tienes un código activo hasta %s.",
			entry.ExpireAt.Format("15:04:05"),
		), nil
	}

	return true, "Cumple políticas, se enviará nuevo código.", nil
}
