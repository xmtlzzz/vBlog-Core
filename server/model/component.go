package model

import "time"

type Component struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Version     string    `gorm:"size:20" json:"version"`
	Code        string    `gorm:"type:text" json:"code"`
	Category    string    `gorm:"size:50" json:"category"`
	Origin      string    `gorm:"size:20;default:built-in" json:"origin"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Component) TableName() string {
	return "components"
}
