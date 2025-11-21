package mail

import (
	"context"
	"fmt"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/config"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/db/repository"
)

var Service *EmailService

func Init(cfg config.EmailConfig) (*EmailService, error) {
	config.ApplyDefaults(&cfg)

	repo := repository.InitBoltRepository(&cfg)
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
