package config

import (
	"strings"
	"testing"
)

func TestDSN(t *testing.T) {
	db := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		Name:     "vblog",
		User:     "vblog_admin",
		Password: "secret",
	}

	dsn := db.DSN()

	if dsn == "" {
		t.Fatal("DSN() returned empty string")
	}

	// Verify DSN contains expected components
	expectedParts := []string{
		"host=localhost",
		"port=5432",
		"dbname=vblog",
		"user=vblog_admin",
		"password=secret",
	}
	for _, part := range expectedParts {
		if !strings.Contains(dsn, part) {
			t.Errorf("DSN() = %q, missing %q", dsn, part)
		}
	}
}

func TestDSNEmptyPassword(t *testing.T) {
	db := DBConfig{
		Host: "localhost",
		Port: "5432",
		Name: "vblog",
		User: "vblog_admin",
	}

	dsn := db.DSN()

	if dsn == "" {
		t.Fatal("DSN() returned empty string")
	}

	if !strings.Contains(dsn, "host=localhost") {
		t.Errorf("DSN() = %q, missing host", dsn)
	}
	if !strings.Contains(dsn, "password=") {
		t.Errorf("DSN() = %q, missing password field", dsn)
	}
}
