package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	AppEnv     string
	AppKey     string
	AppURL     string
	
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
	
	YandexClientID     string
	YandexClientSecret string
	YandexOAuthToken   string
	YandexDefaultCurrency string
	DefaultTimezone    string
	
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppKey:     getEnv("APP_KEY", ""),
		AppURL:     getEnv("APP_URL", "http://localhost:8080"),
		
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBDatabase: getEnv("DB_DATABASE", "reports"),
		DBUsername: getEnv("DB_USERNAME", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		
		YandexClientID:       getEnv("YANDEX_CLIENT_ID", ""),
		YandexClientSecret:   getEnv("YANDEX_CLIENT_SECRET", ""),
		YandexOAuthToken:     getEnv("YANDEX_OAUTH_TOKEN", ""),
		YandexDefaultCurrency: getEnv("YANDEX_DEFAULT_CURRENCY", "RUB"),
		DefaultTimezone:      getEnv("DEFAULT_TIMEZONE", "Europe/Moscow"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

