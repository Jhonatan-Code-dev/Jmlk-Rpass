package domain

import (
	"fmt"
	"time"
)

func CanSendReset(entry *CodeEntry, maxAttempts int, restrictionHours int, allowOverride bool) (bool, string) {
	now := time.Now()
	restriction := time.Duration(restrictionHours) * time.Hour

	if entry == nil {
		return true, "Primer envío permitido."
	}

	if entry.Attempts >= maxAttempts {
		return false, fmt.Sprintf("Máximo de intentos (%d). Espera %.0f horas.", maxAttempts, restriction.Hours())
	}

	if allowOverride {
		return true, "Override activo → se generará nuevo código."
	}

	if entry.Used {
		return false, fmt.Sprintf("Ya usaste tu último código. Espera %.0f horas.", restriction.Hours())
	}

	if now.Before(entry.ExpireAt) {
		return false, fmt.Sprintf("Aún tienes un código activo hasta %s.", entry.ExpireAt.Format("15:04:05"))
	}

	return true, "Cumple políticas, se enviará nuevo código."
}
