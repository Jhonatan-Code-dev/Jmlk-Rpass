package email

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Service es la instancia global opcional del módulo.
var Service *EmailService

// Init inicializa:
// ✔ Configuración
// ✔ Base de datos Bolt
// ✔ Repositorio
// ✔ Sender SMTP
// ✔ Servicio principal (EmailService)
// Y deja listo el módulo para uso externo.
//
// Es el ÚNICO punto de entrada público recomendado.
func Init(cfg EmailConfig) (*EmailService, error) {

	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("'Username' y 'Password' son obligatorios")
	}

	applyDefaults(&cfg)

	// Ubicación estándar del archivo Bolt
	baseDir, _ := os.Getwd()
	dbPath := filepath.Join(baseDir, "storage", "resetpassj.db")

	// Inicializar DB
	db, err := InitBoltDBPath(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error creando DB Bolt: %w", err)
	}

	// Crear bucket
	if err := initBucketIfMissing(db); err != nil {
		return nil, fmt.Errorf("error iniciando bucket Bolt: %w", err)
	}

	// Crear dependencias
	repo := &BoltRepository{db: db}
	smtpClient := &SMTPSender{
		dialer:      NewGomailDialer(cfg),
		senderEmail: cfg.Username,
	}

	// Crear servicio
	svc := NewEmailService(cfg, repo, smtpClient)

	// Guardar instancia global opcional
	Service = svc

	return svc, nil
}

// Helper global para enviar sin gestionar instancia.
// Permite: email.SendReset("correo")
func SendReset(to string) error {
	if Service == nil {
		return fmt.Errorf("email service no inicializado — llama a Init() primero")
	}
	return Service.SendResetPassword(context.Background(), to)
}
