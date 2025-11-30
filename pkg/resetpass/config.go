// pkg/resetpass/config.go
package resetpass

import "time"

type Config struct {
	Host              string
	Port              int
	Username          string
	Password          string
	AppName           string
	Title             string
	CodeLength        int
	CodeValidMinutes  int
	MaxResetAttempts  int
	RestrictionWindow time.Duration
	AllowOverride     bool
	DatabaseFolder    string
	DatabaseName      string
	EmailTimeout      time.Duration
}
