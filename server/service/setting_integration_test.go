package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestSettingService_GetAll(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewSettingService(db)

	// Save a setting first so GetAll has something to return
	svc.Save(map[string]string{"getall_test": "getall_value"})
	defer db.Unscoped().Where("key = ?", "getall_test").Delete(&model.Setting{})

	settings, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if settings["getall_test"] != "getall_value" {
		t.Errorf("expected getall_test='getall_value', got '%s'", settings["getall_test"])
	}
}

func TestSettingService_Save(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewSettingService(db)

	err := svc.Save(map[string]string{
		"test_key": "test_value",
	})
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var setting model.Setting
	db.Where("key = ?", "test_key").First(&setting)
	if setting.Value != "test_value" {
		t.Errorf("expected value 'test_value', got '%s'", setting.Value)
	}

	db.Unscoped().Where("key = ?", "test_key").Delete(&model.Setting{})
}

func TestSettingService_Reset(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewSettingService(db)

	// Save a custom setting
	svc.Save(map[string]string{"reset_test": "value"})

	err := svc.Reset()
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	// Reset overwrites with defaults via Save(DefaultSettings())
	defaults := DefaultSettings()
	settings, _ := svc.GetAll()
	for key, want := range defaults {
		got, ok := settings[key]
		if !ok {
			t.Errorf("missing default key %q after reset", key)
			continue
		}
		if got != want {
			t.Errorf("default %q: expected %q, got %q", key, want, got)
		}
	}

	// Cleanup
	db.Unscoped().Where("key = ?", "reset_test").Delete(&model.Setting{})
	for key := range defaults {
		db.Unscoped().Where("key = ?", key).Delete(&model.Setting{})
	}
}
