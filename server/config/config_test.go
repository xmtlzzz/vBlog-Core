package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	envs := map[string]string{
		"SERVER_PORT":  "9090",
		"DB_HOST":      "db.example.com",
		"DB_PORT":      "5433",
		"DB_NAME":      "testdb",
		"DB_USER":      "testuser",
		"DB_PASSWORD":  "testpass",
		"JWT_SECRET":   "test-jwt-secret",
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	cfg := Load()

	if cfg.Server.Port != "9090" {
		t.Errorf("Server.Port = %q, want %q", cfg.Server.Port, "9090")
	}
	if cfg.DB.Host != "db.example.com" {
		t.Errorf("DB.Host = %q, want %q", cfg.DB.Host, "db.example.com")
	}
	if cfg.DB.Port != "5433" {
		t.Errorf("DB.Port = %q, want %q", cfg.DB.Port, "5433")
	}
	if cfg.DB.Name != "testdb" {
		t.Errorf("DB.Name = %q, want %q", cfg.DB.Name, "testdb")
	}
	if cfg.DB.User != "testuser" {
		t.Errorf("DB.User = %q, want %q", cfg.DB.User, "testuser")
	}
	if cfg.DB.Password != "testpass" {
		t.Errorf("DB.Password = %q, want %q", cfg.DB.Password, "testpass")
	}
	if cfg.JWT.Secret != "test-jwt-secret" {
		t.Errorf("JWT.Secret = %q, want %q", cfg.JWT.Secret, "test-jwt-secret")
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear all relevant env vars to ensure defaults are used
	envVars := []string{
		"SERVER_PORT", "DB_HOST", "DB_PORT", "DB_NAME",
		"DB_USER", "DB_PASSWORD", "JWT_SECRET",
	}
	for _, k := range envVars {
		os.Unsetenv(k)
	}

	cfg := Load()

	if cfg.Server.Port != "8080" {
		t.Errorf("Server.Port = %q, want %q", cfg.Server.Port, "8080")
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("DB.Host = %q, want %q", cfg.DB.Host, "localhost")
	}
	if cfg.DB.Port != "5432" {
		t.Errorf("DB.Port = %q, want %q", cfg.DB.Port, "5432")
	}
	if cfg.DB.Name != "vblog" {
		t.Errorf("DB.Name = %q, want %q", cfg.DB.Name, "vblog")
	}
	if cfg.DB.User != "vblog_admin" {
		t.Errorf("DB.User = %q, want %q", cfg.DB.User, "vblog_admin")
	}
	if cfg.DB.Password != "" {
		t.Errorf("DB.Password = %q, want %q", cfg.DB.Password, "")
	}
	if cfg.JWT.Secret != "vblog-default-secret" {
		t.Errorf("JWT.Secret = %q, want %q", cfg.JWT.Secret, "vblog-default-secret")
	}
}
