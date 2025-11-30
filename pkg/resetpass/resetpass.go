package resetpass

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	bolt "github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/adapters/bolt"
	render "github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/adapters/renderer"
	smtp "github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/adapters/smtp"
	app "github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/app"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/helpers"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/ports"
)

type Service struct {
	s *app.Service
}

func New(cfg Config) (*Service, error) {
	intCfg := app.Config{
		Host:             cfg.Host,
		Port:             cfg.Port,
		Username:         cfg.Username,
		Password:         cfg.Password,
		AppName:          cfg.AppName,
		Title:            cfg.Title,
		CodeLength:       cfg.CodeLength,
		CodeValidMinutes: cfg.CodeValidMinutes,
		MaxResetAttempts: cfg.MaxResetAttempts,
		RestrictionHours: cfg.RestrictionHours,
		AllowOverride:    cfg.AllowOverride,
		DatabaseFolder:   cfg.DatabaseFolder,
		DatabaseName:     cfg.DatabaseName,
		EmailTimeout:     cfg.EmailTimeout,
	}
	intCfg.CodeLength = helpers.OrInt(intCfg.CodeLength, 6)
	intCfg.CodeValidMinutes = helpers.OrInt(intCfg.CodeValidMinutes, 15)
	intCfg.MaxResetAttempts = helpers.OrInt(intCfg.MaxResetAttempts, 3)
	if intCfg.RestrictionHours == 0 {
		intCfg.RestrictionHours = 24
	}
	if intCfg.EmailTimeout == 0 {
		intCfg.EmailTimeout = 30 * time.Second
	}
	if intCfg.DatabaseFolder == "" || intCfg.DatabaseName == "" {
		return nil, fmt.Errorf("database folder and name required for default repository")
	}
	path := filepath.Join(intCfg.DatabaseFolder, intCfg.DatabaseName)
	repo, err := bolt.New(path)
	if err != nil {
		return nil, fmt.Errorf("bolt init: %w", err)
	}
	smtpClient := smtp.NewGomailSender(intCfg.Host, intCfg.Port, intCfg.Username, intCfg.Password)
	renderer := render.NewTemplateRenderer()
	internalSvc := app.NewService(intCfg, repo, smtpClient, renderer)
	return &Service{s: internalSvc}, nil
}

func NewWithAdapters(cfg Config, repo ports.Repository, smtpClient ports.SMTPClient, renderer ports.Renderer) (*Service, error) {
	intCfg := app.Config{
		Host:             cfg.Host,
		Port:             cfg.Port,
		Username:         cfg.Username,
		Password:         cfg.Password,
		AppName:          cfg.AppName,
		Title:            cfg.Title,
		CodeLength:       cfg.CodeLength,
		CodeValidMinutes: cfg.CodeValidMinutes,
		MaxResetAttempts: cfg.MaxResetAttempts,
		RestrictionHours: cfg.RestrictionHours,
		AllowOverride:    cfg.AllowOverride,
		DatabaseFolder:   cfg.DatabaseFolder,
		DatabaseName:     cfg.DatabaseName,
		EmailTimeout:     cfg.EmailTimeout,
	}
	internalSvc := app.NewService(intCfg, repo, smtpClient, renderer)
	return &Service{s: internalSvc}, nil
}

func (s *Service) SendReset(ctx context.Context, email string) error {
	return s.s.SendReset(ctx, email)
}
