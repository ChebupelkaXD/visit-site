package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	From     string
	Password string
	To       string
}

var SMTP SMTPConfig

func init() {
	_ = godotenv.Load()

	SMTP = SMTPConfig{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getEnvInt("SMTP_PORT", 587),
		From:     getEnv("FROM_EMAIL", ""),
		Password: getEnv("FROM_PASSWORD", ""),
		To:       getEnv("TO_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
