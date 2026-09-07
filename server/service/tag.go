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

// TagWithCount includes a post count.
type TagWithCount struct {
	model.Tag
	PostCount int `json:"post_count"`
}

// List returns all tags with post count.
func (s *TagService) List() ([]TagWithCount, error) {
	var tags []TagWithCount
	err := s.DB.Raw(`
		SELECT t.*, COALESCE(pc.cnt, 0) AS post_count
		FROM tags t
		LEFT JOIN (
			SELECT pt.tag_id, COUNT(*) AS cnt
			FROM post_tags pt
			JOIN posts p ON p.id = pt.post_id AND p.status = 'published' AND p.deleted_at IS NULL
			GROUP BY pt.tag_id
		) pc ON pc.tag_id = t.id
		ORDER BY t.name
	`).Scan(&tags).Error
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

// Delete deletes a tag by ID, clearing post associations first to avoid FK errors.
func (s *TagService) Delete(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM post_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Tag{}, id).Error
	})
}
