package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()

	if cfg.Server.Port == "" {
		t.Error("Server.Port should not be empty")
	}
	if cfg.DB.Host == "" {
		t.Error("DB.Host should not be empty")
	}
	if cfg.DB.Port == 0 {
		t.Error("DB.Port should not be zero")
	}
	if cfg.DB.Name == "" {
		t.Error("DB.Name should not be empty")
	}
	if cfg.DB.User == "" {
		t.Error("DB.User should not be empty")
	}
	if cfg.JWT.Secret == "" {
		t.Error("JWT.Secret should not be empty")
	}
}
