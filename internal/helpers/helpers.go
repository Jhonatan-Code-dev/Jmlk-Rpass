package helpers

import (
	"strings"
	"time"
)

func OrString(value, def string) string {
	if strings.TrimSpace(value) == "" {
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
