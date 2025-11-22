package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/domain"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/ports"
)

type Config struct {
	Host             string
	Port             int
	Username         string
	Password         string
	AppName          string
	Title            string
	CodeLength       int
	CodeValidMinutes int
	MaxResetAttempts int
	RestrictionHours int
	AllowOverride    bool
	DatabaseFolder   string
	DatabaseName     string
	EmailTimeout     time.Duration
}

type Service struct {
	cfg      Config
	repo     ports.Repository
	smtp     ports.SMTPClient
	renderer ports.Renderer
	gen      *domain.CodeGenerator
}

func NewService(cfg Config, repo ports.Repository, smtp ports.SMTPClient, renderer ports.Renderer) *Service {
	if cfg.CodeLength == 0 {
		cfg.CodeLength = 6
	}
	if cfg.CodeValidMinutes == 0 {
		cfg.CodeValidMinutes = 15
	}
	if cfg.MaxResetAttempts == 0 {
		cfg.MaxResetAttempts = 3
	}
	if cfg.RestrictionHours == 0 {
		cfg.RestrictionHours = 24
	}
	if cfg.EmailTimeout == 0 {
		cfg.EmailTimeout = 30 * time.Second
	}
	return &Service{
		cfg:      cfg,
		repo:     repo,
		smtp:     smtp,
		renderer: renderer,
		gen:      domain.NewCodeGenerator(cfg.CodeLength),
	}
}

func (s *Service) SendReset(ctx context.Context, to string) error {
	entry, err := s.repo.GetCodeEntry(ctx, to)
	if err != nil {
		return fmt.Errorf("get entry: %w", err)
	}
	can, msg := domain.CanSendReset(entry, s.cfg.MaxResetAttempts, s.cfg.RestrictionHours, s.cfg.AllowOverride)
	log.Println(msg)
	if !can {
		return fmt.Errorf("policy: %s", msg)
	}
	code := s.gen.Generate()
	newEntry := domain.CodeEntry{
		Email:    to,
		Code:     code,
		ExpireAt: time.Now().Add(time.Duration(s.cfg.CodeValidMinutes) * time.Minute),
		Used:     false,
	}
	if entry != nil {
		newEntry.Attempts = entry.Attempts + 1
	} else {
		newEntry.Attempts = 1
	}
	if err := s.repo.SaveCode(ctx, newEntry); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	data := map[string]any{
		"AppName":     s.cfg.AppName,
		"Title":       s.cfg.Title,
		"Minutes":     s.cfg.CodeValidMinutes,
		"MaxAttempts": s.cfg.MaxResetAttempts,
		"Restriction": fmt.Sprintf("%d horas", s.cfg.RestrictionHours),
	}
	html, err := s.renderer.Render(code, data)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	timeout := s.cfg.EmailTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	sendCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := s.smtp.Send(sendCtx, s.cfg.Username, to, s.cfg.Title, html); err != nil {
		return fmt.Errorf("send: %w", err)
	}
	log.Printf("📧 Enviado a %s (intento %d/%d)", to, newEntry.Attempts, s.cfg.MaxResetAttempts)
	return nil
}
