package config

import "os"

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
