package env

import (
	"os"
	"strings"
)

func Get(env, defaultValue string) string {
	if value := os.Getenv(env); strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}
