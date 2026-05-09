package testutil

import (
	"os"
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
	os.Setenv("DB_HOST", "192.168.81.101")
	os.Setenv("DB_USER", "vblog")
	os.Setenv("DB_PASSWORD", "Qwer1234")
	os.Setenv("DB_NAME", "vblog")
	os.Setenv("DB_PORT", "5432")

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
