package email

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	assets "github.com/Jhonatan-Code-dev/Jmlk-Rpass"
)

// RenderTemplate retorna el HTML listo para enviar.
func (s *EmailService) RenderTemplate(ctx context.Context, code string) (string, error) {
	tmpl, err := template.ParseFS(assets.Templates, "templates/reset_password.html")
	if err != nil {
		return "", fmt.Errorf("error cargando plantilla embed: %w", err)
	}

	data := ResetEmailData{
		AppName:     s.conf.AppName,
		Title:       s.conf.Title,
		Code:        code,
		Minutes:     s.conf.CodeValidMinutes,
		MaxAttempts: s.conf.MaxResetAttempts,
		Restriction: fmt.Sprintf("%.0f horas", s.conf.RestrictionPeriod.Hours()),
	}

	var html strings.Builder
	if err := tmpl.Execute(&html, data); err != nil {
		return "", fmt.Errorf("error ejecutando plantilla: %w", err)
	}

	return html.String(), nil
}
