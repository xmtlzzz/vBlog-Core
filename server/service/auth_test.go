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
