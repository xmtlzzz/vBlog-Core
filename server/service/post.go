package service

import (
	"math"

	"gorm.io/gorm"
	"vblog-core/model"
)

// PostService handles blog post CRUD operations.
type PostService struct {
	DB *gorm.DB
}

// NewPostService creates a new PostService.
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

// CalcReadTime calculates reading time in minutes (~500 chars per minute).
func CalcReadTime(content string) int {
	minutes := int(math.Ceil(float64(len(content)) / 500))
	if minutes < 1 {
		return 1
	}
	return minutes
}

// BuildExcerpt truncates content to maxLen with "..." suffix.
func BuildExcerpt(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

// List returns paginated posts with optional filters.
func (s *PostService) List(page, perPage int, tag, status, search string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	q := s.DB.Model(&model.Post{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("title ILIKE ?", "%"+search+"%")
	}
	if tag != "" {
		q = q.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", tag)
	}
	q.Count(&total)

	err := q.Preload("Tags").Order("pinned DESC, created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&posts).Error
	return posts, total, err
}

// GetByID returns a single post by ID with tags preloaded.
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := s.DB.Preload("Tags").First(&post, id).Error
	return &post, err
}

// Create creates a new post, auto-calculating read time and excerpt.
func (s *PostService) Create(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	if post.Excerpt == "" {
		post.Excerpt = BuildExcerpt(post.Content, 200)
	}
	return s.DB.Create(post).Error
}

// Update updates an existing post, recalculating read time.
func (s *PostService) Update(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	return s.DB.Save(post).Error
}

// Delete soft-deletes a post by ID.
func (s *PostService) Delete(id uint) error {
	return s.DB.Delete(&model.Post{}, id).Error
}
