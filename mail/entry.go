// Package mail gestiona la inicialización del servicio de email y envío de correos.
package mail

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/repository"
)

// Instancia global opcional del servicio
var Service *EmailService

// Init inicializa todo el servicio de email: config, DB, repositorio y SMTP.
func Init(cfg config.EmailConfig) (*EmailService, error) {
	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("'Username' y 'Password' son obligatorios")
	}

	config.ApplyDefaults(&cfg)

	baseDir, _ := os.Getwd()
	dbPath := filepath.Join(baseDir, "storage", "resetpassj.db")

	db, err := repository.InitBoltDBPath(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error creando DB Bolt: %w", err)
	}

	if err := repository.InitBucketIfMissing(db); err != nil {
		return nil, fmt.Errorf("error iniciando bucket Bolt: %w", err)
	}

	repo := &repository.BoltRepository{Db: db}

	smtpClient := &SMTPSender{
		Dialer:      NewGomailDialer(cfg),
		SenderEmail: cfg.Username,
	}

	svc := NewEmailService(cfg, repo, smtpClient)
	Service = svc
	return svc, nil
}

// SendReset permite enviar un correo de restablecimiento usando la instancia global.
func SendReset(to string) error {
	if Service == nil {
		return fmt.Errorf("email service no inicializado — llama a Init() primero")
	}
	return Service.SendResetPassword(context.Background(), to)
}
