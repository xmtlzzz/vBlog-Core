package model

import "gorm.io/gorm"

// AutoMigrate creates or updates all tables.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Tag{},
		&Post{},
		&Comment{},
		&Component{},
		&Setting{},
	)
}
