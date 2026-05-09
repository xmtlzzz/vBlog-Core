package model

import "time"

type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PostID     uint      `gorm:"index;not null" json:"post_id"`
	AuthorName string    `gorm:"size:100;not null" json:"author_name"`
	AuthorEmail string   `gorm:"size:255" json:"author_email"`
	Body       string    `gorm:"type:text;not null" json:"body"`
	Status     string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Comment) TableName() string {
	return "comments"
}
