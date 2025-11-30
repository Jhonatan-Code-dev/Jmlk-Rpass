// internal/domain/policy.go
package domain

import (
	"fmt"
	"time"
)

// Ahora recibe restrictionWindow time.Duration
func CanSendReset(entry *CodeEntry, maxAttempts int, restrictionWindow time.Duration, allowOverride bool) (bool, string) {
	now := time.Now()

	if entry == nil {
		return true, "Primer envío permitido."
	}

	// Si excedió intentos
	if entry.Attempts >= maxAttempts {
		return false, fmt.Sprintf(
			"Máximo de intentos (%d). Espera %s.",
			maxAttempts,
			restrictionWindow.String(),
		)
	}

	// Override: siempre permite
	if allowOverride {
		return true, "Override activo → se generará nuevo código."
	}

	// Ya usó el anterior
	if entry.Used {
		return false, fmt.Sprintf(
			"Ya usaste tu último código. Espera %s.",
			restrictionWindow.String(),
		)
	}

	// Aún tiene un código activo
	if now.Before(entry.ExpireAt) {
		return false, fmt.Sprintf("Aún tienes un código activo hasta %s.",
			entry.ExpireAt.Format("15:04:05"),
		)
	}

	return true, "Cumple políticas, se enviará nuevo código."
}
