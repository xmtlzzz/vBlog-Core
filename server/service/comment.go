package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

// CommentService handles comment CRUD operations.
type CommentService struct {
	DB     *gorm.DB
	LogSvc *ChangeLogService
}

// NewCommentService creates a new CommentService.
func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{DB: db, LogSvc: NewChangeLogService(db)}
}

// List returns paginated comments with optional status filter and search.
func (s *CommentService) List(page, perPage int, status, search string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	q := s.DB.Model(&model.Comment{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("body ILIKE ?", "%"+search+"%")
	}
	q.Count(&total)

	err := q.Order("created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&comments).Error
	return comments, total, err
}

// Create creates a new comment with auto-set status "pending".
func (s *CommentService) Create(c *model.Comment) error {
	c.Status = "pending"
	if err := s.DB.Create(c).Error; err != nil {
		return err
	}
	s.LogSvc.Write("new_comment", &c.ID, c.Body, "")
	return nil
}

// Approve sets a comment's status to "approved".
func (s *CommentService) Approve(id uint) error {
	return s.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", "approved").Error
}

// MarkSpam sets a comment's status to "spam".
func (s *CommentService) MarkSpam(id uint) error {
	return s.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", "spam").Error
}

// Delete deletes a comment by ID.
func (s *CommentService) Delete(id uint) error {
	return s.DB.Delete(&model.Comment{}, id).Error
}
