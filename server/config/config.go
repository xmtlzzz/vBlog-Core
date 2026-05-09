package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port string
}

// DBConfig holds PostgreSQL database configuration.
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret           string
	ExpireHours      int
	RefreshExpireHours int
}

func init() {
	// Find project root by looking for go.mod
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Try to load .env from project root (one level up from server/)
	envPath := filepath.Join(filepath.Dir(dir), ".env")
	godotenv.Load(envPath)
}

// Load reads configuration from environment variables with defaults.
func Load() Config {
	return Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "vblog"),
			User:     getEnv("DB_USER", "vblog_admin"),
			Password: getEnv("DB_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "vblog-default-secret"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
