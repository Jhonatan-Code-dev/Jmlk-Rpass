package email

import "time"

// CodeEntry representa un registro persistido para un correo.
type CodeEntry struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expire_at"`
	Attempts int       `json:"attempts"`
	Used     bool      `json:"used"`
}

// ResetEmailData se usa para renderizar la plantilla.
type ResetEmailData struct {
	AppName     string
	Title       string
	Code        string
	Minutes     int
	MaxAttempts int
	Restriction string
}
