package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DSN returns the PostgreSQL connection string.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name,
	)
}

// Connect opens a database connection using GORM.
func (d DBConfig) Connect() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(d.DSN()), &gorm.Config{})
}
