package helpers

import (
	"fmt"
	"strings"
	"time"
)

func OrString(value, def string) string {
	if value == "" {
		return def
	}
	return value
}

func OrInt(value, def int) int {
	if value == 0 {
		return def
	}
	return value
}

func OrDuration(value, def time.Duration) time.Duration {
	if value == 0 {
		return def
	}
	return value
}

func RequiredString(value string, fieldName string) string {
	if strings.TrimSpace(value) == "" {
		panic(fmt.Sprintf("el campo obligatorio '%s' no puede estar vacío", fieldName))
	}
	return value
}
