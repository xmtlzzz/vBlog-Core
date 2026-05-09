package model

import "testing"

func TestUserTableName(t *testing.T) {
	u := User{}
	if u.TableName() != "users" {
		t.Errorf("expected 'users', got '%s'", u.TableName())
	}
}

func TestUserDefaultValues(t *testing.T) {
	u := User{}
	if u.Username != "" {
		t.Errorf("expected default username '', got '%s'", u.Username)
	}
	if u.Role != "" {
		t.Errorf("expected default role '', got '%s'", u.Role)
	}
}
