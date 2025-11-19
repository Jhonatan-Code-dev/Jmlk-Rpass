package mail

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/models"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/rules"
)

type EmailService struct {
	smtp SMTPClient
	repo Repository
	conf config.EmailConfig
	gen  *rules.CodeGenerator
}

func NewEmailService(cfg config.EmailConfig, repo Repository, smtp SMTPClient) *EmailService {
	config.ApplyDefaults(&cfg)
	return &EmailService{
		smtp: smtp,
		repo: repo,
		conf: cfg,
		gen:  rules.NewCodeGenerator(cfg.CodeLength),
	}
}

func (s *EmailService) SendResetPassword(ctx context.Context, to string) error {
	entry, _ := s.repo.GetCodeEntry(ctx, to)
	canSend, msg := rules.CanSendReset(entry, s.conf.MaxResetAttempts, s.conf.RestrictionPeriod, s.conf.AllowOverride)
	log.Println(msg)
	if !canSend {
		return fmt.Errorf("no se enviará: %s", msg)
	}

	code := s.gen.Generate()
	newEntry := models.CodeEntry{
		Email:    to,
		Code:     code,
		ExpireAt: time.Now().Add(time.Duration(s.conf.CodeValidMinutes) * time.Minute),
		Used:     false,
	}
	if entry != nil {
		newEntry.Attempts = entry.Attempts + 1
	} else {
		newEntry.Attempts = 1
	}

	if err := s.repo.SaveCode(ctx, newEntry); err != nil {
		return fmt.Errorf("error guardando código: %w", err)
	}

	html, err := s.RenderTemplate(ctx, code)
	if err != nil {
		return fmt.Errorf("error renderizando plantilla: %w", err)
	}

	sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := s.smtp.Send(sendCtx, s.conf.Username, to, s.conf.Title, html); err != nil {
		return fmt.Errorf("error enviando correo: %w", err)
	}

	log.Printf("📧 Enviado a %s (intento %d/%d)", to, newEntry.Attempts, s.conf.MaxResetAttempts)
	return nil
}
