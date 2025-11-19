// Package mail genera códigos numéricos ultra rápidos y sin overhead.
package email

import (
	"math/rand"
	"sync"
	"time"
)

const digits = "0123456789"

type EmailService struct {
	codeLength int
	prngPool   sync.Pool
}

func NewEmailService(codeLength int) *EmailService {
	s := &EmailService{codeLength: codeLength}
	s.prngPool.New = func() interface{} { return rand.New(rand.NewSource(time.Now().UnixNano())) }
	return s
}

func (s *EmailService) GenerateCode() string {
	r := s.prngPool.Get().(*rand.Rand)
	var buf [64]byte
	for i := 0; i < s.codeLength; i++ {
		buf[i] = digits[r.Intn(10)]
	}
	s.prngPool.Put(r)
	return string(buf[:s.codeLength])
}
