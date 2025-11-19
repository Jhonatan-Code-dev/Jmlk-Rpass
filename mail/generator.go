package email

import (
	"math/rand"
	"time"
)

// GenerateCode produce un código numérico de longitud configurada.
// Es método del servicio para poder usar la configuración.
func (s *EmailService) GenerateCode() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	digits := "0123456789"
	code := make([]byte, s.conf.CodeLength)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}
