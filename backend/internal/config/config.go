package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	AppEnv      string
	AppKey      string
	AppURL      string // Backend URL
	FrontendURL string // Frontend URL for OAuth redirects

	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string

	YandexClientID        string
	YandexClientSecret    string
	YandexOAuthToken      string
	YandexOAuthScopes     string
	YandexDefaultCurrency string
	YandexDirectSandbox   bool // Use sandbox environment for Yandex Direct API
	DefaultTimezone       string

	JWTSecret string // Secret key for JWT tokens
	JWTExpiry int    // JWT token expiry in hours (default 24)
	OllamaAPIKey  string // Ollama API key for metrics analysis
	OllamaAPIURL  string // Ollama API URL (default: https://api.ollama.com/v1)
	OllamaModel   string // Ollama model name (default: llama3.2)

	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int


	LogLevel string
}

// Load loads configuration from environment variables
// It tries to load .env file if it exists (for development)
func Load() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		AppKey:      getEnv("APP_KEY", ""),
		AppURL:      getEnv("APP_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"), // Default frontend port

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBDatabase: getEnv("DB_DATABASE", "reports"),
		DBUsername: getEnv("DB_USERNAME", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),

		YandexClientID:        getEnv("YANDEX_CLIENT_ID", ""),
		YandexClientSecret:    getEnv("YANDEX_CLIENT_SECRET", ""),
		YandexOAuthToken:      getEnv("YANDEX_OAUTH_TOKEN", ""),
		YandexOAuthScopes:     getEnv("YANDEX_OAUTH_SCOPES", ""), // If empty, scopes from app registration will be used
		YandexDefaultCurrency: getEnv("YANDEX_DEFAULT_CURRENCY", "RUB"),
		YandexDirectSandbox:   getEnv("YANDEX_DIRECT_SANDBOX", "false") == "true", // Use sandbox for testing
		DefaultTimezone:       getEnv("DEFAULT_TIMEZONE", "Europe/Moscow"),
		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:             getEnvInt("JWT_EXPIRY", 24), // Default 24 hours
		OllamaAPIKey:          getEnv("OLLAMA_API_KEY", ""),
		OllamaAPIURL:          getEnv("OLLAMA_API_URL", "https://ollama.com/api"),
		OllamaModel:           getEnv("OLLAMA_MODEL", "glm-4.6"),
		LogLevel:              getEnv("LOG_LEVEL", "info"),


		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:      getEnvInt("REDIS_DB", 0),

		LogLevel: getEnv("LOG_LEVEL", "info"),

	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// UpdateEnvFile updates or adds a key-value pair in .env file
// For MVP: saves OAuth token to .env file
func UpdateEnvFile(key, value string) error {
	// Find .env file (usually in backend directory)
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try backend/.env
		envPath = "backend/.env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			// Create new .env file if it doesn't exist
			return createEnvFile(envPath, key, value)
		}
	}

	// Read existing .env file
	file, err := os.OpenFile(envPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	var lines []string
	keyFound := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Check if this line contains the key
		if strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			// Update existing line
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			keyFound = true
		} else {
			// Keep original line
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	// If key not found, add it at the end
	if !keyFound {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	// Write back to file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate .env file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek .env file: %w", err)
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to .env file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush .env file: %w", err)
	}

	return nil
}

// createEnvFile creates a new .env file with the given key-value pair
func createEnvFile(path, key, value string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	return err
}
