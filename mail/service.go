package email

import (
	"context"
	"fmt"
	"log"
	"time"
)

// EmailService es la fachada principal del módulo (orquestación).
type EmailService struct {
	smtp SMTPClient
	repo Repository
	conf EmailConfig
}

// NewEmailService construye el servicio con inyección de dependencias.
func NewEmailService(cfg EmailConfig, repo Repository, smtp SMTPClient) *EmailService {
	applyDefaults(&cfg)
	return &EmailService{
		smtp: smtp,
		repo: repo,
		conf: cfg,
	}
}

// SendResetPassword hace todo el flujo: política, generación, persistencia, render y envío.
func (s *EmailService) SendResetPassword(ctx context.Context, to string) error {
	canSend, msg, err := s.CanSendReset(ctx, to)
	if err != nil {
		return fmt.Errorf("error evaluando políticas: %w", err)
	}
	log.Println(msg)
	if !canSend {
		return fmt.Errorf("no se enviará: %s", msg)
	}

	// Generar código y preparar entry
	code := s.GenerateCode()
	entry := CodeEntry{
		Email:    to,
		Code:     code,
		ExpireAt: time.Now().Add(time.Duration(s.conf.CodeValidMinutes) * time.Minute),
		Used:     false,
	}

	// Calcular attempts
	oldEntry, err := s.repo.GetCodeEntry(ctx, to)
	if err == nil && oldEntry != nil {
		entry.Attempts = oldEntry.Attempts + 1
	} else {
		entry.Attempts = 1
	}

	// Guardar código
	if err := s.repo.SaveCode(ctx, entry); err != nil {
		return fmt.Errorf("error guardando en repo: %w", err)
	}

	// Renderizar plantilla
	html, err := s.RenderTemplate(ctx, code)
	if err != nil {
		return fmt.Errorf("error renderizando plantilla: %w", err)
	}

	// Enviar correo (con contexto para cancelación/timeout)
	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.smtp.Send(sendCtx, s.conf.Username, to, s.conf.Title, html); err != nil {
		return fmt.Errorf("error enviando correo: %w", err)
	}

	log.Printf("📧 Enviado a %s (intento %d/%d)\n", to, entry.Attempts, s.conf.MaxResetAttempts)
	return nil
}
