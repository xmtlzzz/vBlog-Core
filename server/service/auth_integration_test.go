package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestAuthService_Register(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewAuthService(db)

	user, err := svc.Register("testuser", "password123", "test@example.com")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user ID to be set")
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}
	if user.Password == "password123" {
		t.Error("password should be hashed, not plaintext")
	}

	db.Unscoped().Delete(user)
}

func TestAuthService_Login(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewAuthService(db)

	svc.Register("logintest", "mypassword", "login@test.com")

	user, err := svc.Login("logintest", "mypassword")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if user.Username != "logintest" {
		t.Errorf("expected username 'logintest', got '%s'", user.Username)
	}

	db.Unscoped().Where("username = ?", "logintest").Delete(&model.User{})
}

func TestAuthService_LoginWrongPassword(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewAuthService(db)

	svc.Register("wrongpass", "correct", "wrong@test.com")

	_, err := svc.Login("wrongpass", "incorrect")
	if err == nil {
		t.Error("expected error for wrong password")
	}

	db.Unscoped().Where("username = ?", "wrongpass").Delete(&model.User{})
}

func TestAuthService_LoginNonexistent(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewAuthService(db)

	_, err := svc.Login("nonexistent", "password")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}
