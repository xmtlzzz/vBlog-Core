package service

import "testing"

func TestNewSettingService(t *testing.T) {
	s := NewSettingService(nil)
	if s == nil {
		t.Fatal("NewSettingService returned nil")
	}
	if s.DB != nil {
		t.Error("expected DB to be nil")
	}
}

func TestSettingServiceStruct(t *testing.T) {
	s := &SettingService{}
	if s.DB != nil {
		t.Error("expected zero-value DB to be nil")
	}
}

func TestDefaultSettings(t *testing.T) {
	defaults := DefaultSettings()
	if defaults == nil {
		t.Fatal("DefaultSettings returned nil")
	}

	expected := map[string]string{
		"site_title":       "vBlog",
		"site_description": "A lightweight blog for geeks",
		"site_url":         "",
		"posts_per_page":   "10",
		"theme":            "default",
		"footer_text":      "Powered by vBlog Core",
	}

	for key, want := range expected {
		got, ok := defaults[key]
		if !ok {
			t.Errorf("missing default key %q", key)
			continue
		}
		if got != want {
			t.Errorf("default %q: expected %q, got %q", key, want, got)
		}
	}
}
