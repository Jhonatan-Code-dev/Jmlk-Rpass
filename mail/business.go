package email

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"text/template"
	"time"

	assets "github.com/Jhonatan-Code-dev/Jmlk-Rpass"
	"gopkg.in/gomail.v2"
)

// =====================================================
// 🔢 GENERAR CÓDIGO
// =====================================================
func (e *EmailService) generateCode() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	digits := "0123456789"
	code := make([]byte, e.conf.CodeLength)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}

// =====================================================
// ⚖️ POLÍTICAS DE ENVÍO
// =====================================================
func (e *EmailService) CanSendReset(email string) (bool, string) {
	entry, _ := GetCodeEntry(e.db, email)
	now := time.Now()

	// No existe registro → enviar
	if entry == nil {
		return true, "Primer envío permitido."
	}

	// 1. Regla que NUNCA se salta: máximo de intentos
	if entry.Attempts >= e.conf.MaxResetAttempts {
		return false, fmt.Sprintf(
			"Máximo de intentos (%d). Espera %.0f horas.",
			e.conf.MaxResetAttempts,
			e.conf.RestrictionPeriod.Hours(),
		)
	}

	// 2. Si AllowOverride está activo → ignorar expiración y Used
	//    pero NO ignorar el límite de intentos (ya validado arriba)
	if e.conf.AllowOverride {
		return true, "Override activo → se generará nuevo código."
	}

	// 3. Reglas normales (estrictas)
	if entry.Used {
		return false, fmt.Sprintf(
			"Ya usaste tu último código. Espera %.0f horas.",
			e.conf.RestrictionPeriod.Hours(),
		)
	}

	if now.Before(entry.ExpireAt) {
		return false, fmt.Sprintf(
			"Aún tienes un código activo hasta %s.",
			entry.ExpireAt.Format("15:04:05"),
		)
	}

	return true, "Cumple políticas, se enviará nuevo código."
}

// =====================================================
// 📤 ENVÍO DE CORREO
// =====================================================
type ResetEmailData struct {
	AppName     string
	Title       string
	Code        string
	Minutes     int
	MaxAttempts int
	Restriction string
}

func (e *EmailService) SendResetPassword(to string) error {
	canSend, msg := e.CanSendReset(to)
	log.Println(msg)
	if !canSend {
		return fmt.Errorf("❌ no se enviará: %s", msg)
	}

	code := e.generateCode()
	entry := CodeEntry{
		Email:    to,
		Code:     code,
		ExpireAt: time.Now().Add(time.Duration(e.conf.CodeValidMinutes) * time.Minute),
		Used:     false,
	}

	old, _ := GetCodeEntry(e.db, to)
	if old != nil {
		entry.Attempts = old.Attempts + 1
	} else {
		entry.Attempts = 1
	}
	if err := SaveCode(e.db, entry); err != nil {
		return fmt.Errorf("error guardando en BD: %w", err)
	}

	html, err := e.renderTemplate(code)
	if err != nil {
		return err
	}

	if err := e.send(to, e.conf.Title, html); err != nil {
		return fmt.Errorf("error enviando correo: %w", err)
	}

	log.Printf("📧 Enviado a %s (intento %d/%d)\n", to, entry.Attempts, e.conf.MaxResetAttempts)
	return nil
}

// =====================================================
// 🧱 RENDER HTML TEMPLATE
// =====================================================
func (e *EmailService) renderTemplate(code string) (string, error) {
	tmpl, err := template.ParseFS(assets.Templates, "templates/reset_password.html")

	if err != nil {
		return "", fmt.Errorf("error cargando plantilla embed: %w", err)
	}

	data := ResetEmailData{
		AppName:     e.conf.AppName,
		Title:       e.conf.Title,
		Code:        code,
		Minutes:     e.conf.CodeValidMinutes,
		MaxAttempts: e.conf.MaxResetAttempts,
		Restriction: fmt.Sprintf("%.0f horas", e.conf.RestrictionPeriod.Hours()),
	}

	var html strings.Builder

	if err := tmpl.Execute(&html, data); err != nil {
		return "", fmt.Errorf("error ejecutando plantilla: %w", err)
	}

	return html.String(), nil
}

// =====================================================
// ✉️ MÉTODO PRIVADO DE ENVÍO
// =====================================================
func (e *EmailService) send(to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)
	return e.dialer.DialAndSend(msg)
}
