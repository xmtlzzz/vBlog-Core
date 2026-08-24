package testutil

import (
	"testing"

	"vblog-core/config"
	"vblog-core/model"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func GetTestDB(t *testing.T) *gorm.DB {
	if testDB != nil {
		return testDB
	}
	// Connection info comes from environment variables or the local
	// config/config.toml (see config/config.example.toml).
	cfg := config.Load()
	db, err := cfg.DB.Connect()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate for test tables
	if err := model.AutoMigrate(db); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	testDB = db
	return db
}

func CleanupTables(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		db.Where("1 = 1").Delete(table)
	}
}
