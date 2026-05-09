package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Excerpt   string         `gorm:"size:500" json:"excerpt"`
	Status    string         `gorm:"size:20;default:draft" json:"status"`
	Pinned    bool           `gorm:"default:false" json:"pinned"`
	Views     int            `gorm:"default:0" json:"views"`
	ReadTime  int            `gorm:"default:0" json:"read_time"`
	AuthorID  uint           `json:"author_id"`
	Tags      []Tag          `gorm:"many2many:post_tags;" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}

// FormatDate returns a date string in local timezone for JSON serialization.
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("2006-01-02")
}

// FormatDateTime returns a datetime string in local timezone.
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("2006-01-02 15:04")
}
