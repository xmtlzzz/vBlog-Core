package model

import "time"

type ChangeLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ChangeType string    `gorm:"size:50;not null" json:"change_type"`
	TargetID   *uint     `json:"target_id"`
	Title      string    `gorm:"size:200" json:"title"`
	Detail     string    `gorm:"type:text" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ChangeLog) TableName() string {
	return "change_log"
}
