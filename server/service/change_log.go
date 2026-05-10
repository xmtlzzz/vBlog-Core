package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

// ChangeLogService handles change log recording and retrieval.
type ChangeLogService struct {
	DB *gorm.DB
}

// NewChangeLogService creates a new ChangeLogService.
func NewChangeLogService(db *gorm.DB) *ChangeLogService {
	return &ChangeLogService{DB: db}
}

// Write records a new change log entry.
func (s *ChangeLogService) Write(changeType string, targetID *uint, title, detail string) error {
	return s.DB.Create(&model.ChangeLog{
		ChangeType: changeType,
		TargetID:   targetID,
		Title:      title,
		Detail:     detail,
	}).Error
}

// GetAfterID returns all change log entries with ID greater than afterID.
func (s *ChangeLogService) GetAfterID(afterID int64) ([]model.ChangeLog, error) {
	var logs []model.ChangeLog
	err := s.DB.Where("id > ?", afterID).Order("id ASC").Find(&logs).Error
	return logs, err
}

// GetLatestID returns the ID of the most recent change log entry, or 0 if none exist.
func (s *ChangeLogService) GetLatestID() (int64, error) {
	var log model.ChangeLog
	err := s.DB.Order("id DESC").First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return int64(log.ID), nil
}

// Backfill populates change_log with existing posts and comments.
// Safe to call multiple times — skips if change_log already has entries.
func (s *ChangeLogService) Backfill() error {
	var count int64
	s.DB.Model(&model.ChangeLog{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Backfill published posts
	var posts []model.Post
	s.DB.Where("status = ?", "published").Order("created_at ASC").Find(&posts)
	for _, p := range posts {
		s.DB.Create(&model.ChangeLog{
			ChangeType: "new_post",
			TargetID:   &p.ID,
			Title:      p.Title,
			CreatedAt:  p.CreatedAt,
		})
	}

	// Backfill comments
	var comments []model.Comment
	s.DB.Order("created_at ASC").Find(&comments)
	for _, c := range comments {
		s.DB.Create(&model.ChangeLog{
			ChangeType: "new_comment",
			TargetID:   &c.ID,
			Title:      c.Body,
			CreatedAt:  c.CreatedAt,
		})
	}

	return nil
}
