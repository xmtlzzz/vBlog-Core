package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Addr     string
	Port     string
	GrpcPort string
}

// DBConfig holds PostgreSQL database configuration.
type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	Secret string
}

var projectRoot string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

// Load reads configuration from config.toml via Viper.
func Load() Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(filepath.Join(projectRoot, "config"))

	// Environment overrides (used by docker-compose and tests)
	v.BindEnv("http.addr", "HTTP_ADDR")
	v.BindEnv("http.port", "HTTP_PORT")
	v.BindEnv("http.grpc_port", "HTTP_GRPC_PORT")
	v.BindEnv("postgres.host", "DB_HOST", "PGHOST")
	v.BindEnv("postgres.port", "DB_PORT", "PGPORT")
	v.BindEnv("postgres.name", "DB_NAME", "PGDATABASE")
	v.BindEnv("postgres.user", "DB_USER", "PGUSER")
	v.BindEnv("postgres.password", "DB_PASSWORD", "PGPASSWORD")
	v.BindEnv("jwt.secret", "JWT_SECRET")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("failed to read config: %v", err)
		}
		log.Println("config.toml not found; using environment variables and defaults")
	}

	cfg := Config{
		Server: ServerConfig{
			Addr:     v.GetString("http.addr"),
			Port:     v.GetString("http.port"),
			GrpcPort: v.GetString("http.grpc_port"),
		},
		DB: DBConfig{
			Host:     v.GetString("postgres.host"),
			Port:     v.GetInt("postgres.port"),
			Name:     v.GetString("postgres.name"),
			User:     v.GetString("postgres.user"),
			Password: v.GetString("postgres.password"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("jwt.secret"),
		},
	}

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = "0.0.0.0"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.GrpcPort == "" {
		cfg.Server.GrpcPort = "50051"
	}
	if cfg.DB.Port == 0 {
		cfg.DB.Port = 5432
	}

	return cfg
}

