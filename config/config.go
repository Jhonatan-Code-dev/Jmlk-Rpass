package config

import (
	"time"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/helpers"
)

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

	// NUEVO: carpeta y nombre de la base de datos
	DatabaseFolder string
	DatabaseName   string
}

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

		// Valores por defecto
		DatabaseFolder: "yona",
		DatabaseName:   "resetpass.db",
	}
}

func ApplyDefaults(cfg *EmailConfig) {
	def := NewDefaultConfig()

	cfg.Host = helpers.OrString(cfg.Host, def.Host)
	cfg.Port = helpers.OrInt(cfg.Port, def.Port)
	cfg.AppName = helpers.OrString(cfg.AppName, def.AppName)
	cfg.Title = helpers.OrString(cfg.Title, def.Title)
	cfg.CodeLength = helpers.OrInt(cfg.CodeLength, def.CodeLength)
	cfg.CodeValidMinutes = helpers.OrInt(cfg.CodeValidMinutes, def.CodeValidMinutes)
	cfg.MaxResetAttempts = helpers.OrInt(cfg.MaxResetAttempts, def.MaxResetAttempts)
	cfg.RestrictionPeriod = helpers.OrDuration(cfg.RestrictionPeriod, def.RestrictionPeriod)

	// Aplicar defaults nuevos
	cfg.DatabaseFolder = helpers.OrString(cfg.DatabaseFolder, def.DatabaseFolder)
	cfg.DatabaseName = helpers.OrString(cfg.DatabaseName, def.DatabaseName)
}
