package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

// TagService handles tag CRUD operations.
type TagService struct {
	DB *gorm.DB
}

// NewTagService creates a new TagService.
func NewTagService(db *gorm.DB) *TagService {
	return &TagService{DB: db}
}

// List returns all tags with post_count calculated.
func (s *TagService) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := s.DB.Find(&tags).Error
	return tags, err
}

// Create creates a new tag.
func (s *TagService) Create(tag *model.Tag) error {
	return s.DB.Create(tag).Error
}

// Update updates an existing tag.
func (s *TagService) Update(tag *model.Tag) error {
	return s.DB.Save(tag).Error
}

// Delete deletes a tag by ID.
func (s *TagService) Delete(id uint) error {
	return s.DB.Delete(&model.Tag{}, id).Error
}
