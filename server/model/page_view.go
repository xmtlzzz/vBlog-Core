package model

import "time"

type PageView struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IP        string    `gorm:"size:45" json:"ip"`
	Path      string    `gorm:"size:500" json:"path"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

func (PageView) TableName() string {
	return "page_views"
}
