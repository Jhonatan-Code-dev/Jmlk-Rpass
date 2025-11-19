package email

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
	"gopkg.in/gomail.v2"
)

// EmailConfig agrupa todas las opciones del servicio.
type EmailConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	AppName           string
	Title             string
	CodeLength        int
	CodeValidMinutes  int
	MaxResetAttempts  int
	RestrictionPeriod time.Duration
	AllowOverride     bool
}

// NewDefaultConfig devuelve un EmailConfig con valores por defecto.
func NewDefaultConfig() EmailConfig {
	return EmailConfig{
		Host:              "smtp.gmail.com",
		Port:              587,
		AppName:           "MiApp",
		Title:             "Restablecimiento de contraseña",
		CodeLength:        6,
		CodeValidMinutes:  15,
		MaxResetAttempts:  3,
		RestrictionPeriod: 24 * time.Hour,
		AllowOverride:     true,
	}
}

// applyDefaults completa cfg con valores por defecto si están vacíos.
func applyDefaults(cfg *EmailConfig) {
	def := NewDefaultConfig()

	if cfg.Host == "" {
		cfg.Host = def.Host
	}
	if cfg.Port == 0 {
		cfg.Port = def.Port
	}
	if cfg.AppName == "" {
		cfg.AppName = def.AppName
	}
	if cfg.Title == "" {
		cfg.Title = def.Title
	}
	if cfg.CodeLength == 0 {
		cfg.CodeLength = def.CodeLength
	}
	if cfg.CodeValidMinutes == 0 {
		cfg.CodeValidMinutes = def.CodeValidMinutes
	}
	if cfg.MaxResetAttempts == 0 {
		cfg.MaxResetAttempts = def.MaxResetAttempts
	}
	if cfg.RestrictionPeriod == 0 {
		cfg.RestrictionPeriod = def.RestrictionPeriod
	}
	// AllowOverride puede ser false intencionalmente, no lo sobreescribimos si usuario lo puso.
}

// InitBoltDBPath crea directorio y abre la db bolt
func InitBoltDBPath(dbPath string) (*bbolt.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, fmt.Errorf("mkdir storage: %w", err)
	}
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("open bolt db: %w", err)
	}
	return db, nil
}

// NewGomailDialer crea un dialer gomail desde config (helper).
func NewGomailDialer(cfg EmailConfig) *gomail.Dialer {
	return gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
}

// Helper de inicialización "rápida" que crea el repositorio Bolt y el sender SMTP real.
// Úsalo desde main si quieres la configuración por defecto basada en filesystem.
func NewServiceWithBoltAndSMTP(cfg EmailConfig, dbPath string) (*EmailService, error) {
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("'Username' y 'Password' son obligatorios")
	}
	applyDefaults(&cfg)

	db, err := InitBoltDBPath(dbPath)
	if err != nil {
		return nil, err
	}

	// Inicializar bucket si falta
	if err := initBucketIfMissing(db); err != nil {
		db.Close()
		return nil, err
	}

	repo := &BoltRepository{db: db}
	dialer := NewGomailDialer(cfg)
	sender := &SMTPSender{dialer: dialer, senderEmail: cfg.Username}

	svc := NewEmailService(cfg, repo, sender)
	log.Printf("✅ Servicio '%s' listo | Código: %d dígitos | Validez: %d min | Intentos: %d",
		cfg.AppName, cfg.CodeLength, cfg.CodeValidMinutes, cfg.MaxResetAttempts)

	return svc, nil
}
