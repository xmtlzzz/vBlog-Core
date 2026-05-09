package middleware

import (
	"testing"
	"time"
)

const testSecret = "test-secret-key-for-jwt"

func TestGenerateAndValidateToken(t *testing.T) {
	tokenStr, err := GenerateToken(42, "admin", testSecret, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}

	claims, err := ValidateToken(tokenStr, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("expected Username 'admin', got %s", claims.Username)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected non-nil ExpiresAt")
	}
	if claims.IssuedAt == nil {
		t.Fatal("expected non-nil IssuedAt")
	}
}

func TestValidateExpiredToken(t *testing.T) {
	tokenStr, err := GenerateToken(1, "user", testSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(tokenStr, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateInvalidToken(t *testing.T) {
	_, err := ValidateToken("this-is-not-a-valid-jwt", testSecret)
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
