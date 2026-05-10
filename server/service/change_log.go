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
