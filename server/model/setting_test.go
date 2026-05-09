package model

import "testing"

func TestSettingTableName(t *testing.T) {
	s := Setting{}
	if s.TableName() != "settings" {
		t.Errorf("expected 'settings', got '%s'", s.TableName())
	}
}

func TestSettingDefaultValues(t *testing.T) {
	s := Setting{}
	if s.Key != "" {
		t.Errorf("expected default key '', got '%s'", s.Key)
	}
	if s.Value != "" {
		t.Errorf("expected default value '', got '%s'", s.Value)
	}
}
