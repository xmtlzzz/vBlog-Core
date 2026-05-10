package model

import "time"

type DailyStats struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StatDate     time.Time `gorm:"type:date;uniqueIndex;not null" json:"stat_date"`
	PV           int64     `gorm:"default:0" json:"pv"`
	UV           int64     `gorm:"default:0" json:"uv"`
	PostCount    int       `gorm:"default:0" json:"post_count"`
	ViewTotal    int64     `gorm:"default:0" json:"view_total"`
	CommentCount int       `gorm:"default:0" json:"comment_count"`
	TagCount     int       `gorm:"default:0" json:"tag_count"`
	CreatedAt    time.Time `json:"created_at"`
}

func (DailyStats) TableName() string {
	return "daily_stats"
}
