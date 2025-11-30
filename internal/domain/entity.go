// internal/domain/entity.go
package domain

import "time"

type CodeEntry struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expire_at"`
	Attempts int       `json:"attempts"`
	Used     bool      `json:"used"`
}
