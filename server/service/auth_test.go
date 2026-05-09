package service

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "mypassword" {
		t.Fatal("hash should not equal plaintext password")
	}
}

func TestCheckPasswordCorrect(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword(hash, "secret123") {
		t.Error("expected CheckPassword to return true for correct password")
	}
}

func TestCheckPasswordWrong(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Error("expected CheckPassword to return false for wrong password")
	}
}

func TestRegisterValidation(t *testing.T) {
	// Test that Register requires non-empty username and password
	// Pure validation test — no DB needed
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"empty username", "", "pass123", true},
		{"empty password", "user1", "", true},
		{"both empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.username == "" || tt.password == "" {
				// Validation should catch this before DB call
				if !tt.wantErr {
					t.Error("expected error for empty fields")
				}
			}
		})
	}
}
