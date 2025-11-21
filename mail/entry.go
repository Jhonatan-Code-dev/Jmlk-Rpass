package mail

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/repository"
)

var Service *EmailService

func Init(cfg config.EmailConfig) (*EmailService, error) {
	config.ApplyDefaults(&cfg)

	// Construir ruta dinámica de la DB
	baseDir, _ := os.Getwd()
	dbPath := filepath.Join(baseDir, cfg.DatabaseFolder, cfg.DatabaseName)

	db, err := repository.InitBoltDBPath(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error creando DB Bolt: %w", err)
	}

	if err := repository.InitBucketIfMissing(db); err != nil {
		return nil, fmt.Errorf("error iniciando bucket Bolt: %w", err)
	}

	repo := &repository.BoltRepository{DB: db}

	smtpClient := &SMTPSender{
		Dialer:      NewGomailDialer(cfg),
		SenderEmail: cfg.Username,
	}

	svc := NewEmailService(cfg, repo, smtpClient)
	Service = svc

	return svc, nil
}

func SendReset(to string) error {
	if Service == nil {
		return fmt.Errorf("email service no inicializado — llama a Init() primero")
	}
	return Service.SendResetPassword(context.Background(), to)
}
