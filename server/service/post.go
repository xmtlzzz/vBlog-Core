package service

import (
	"math"

	"gorm.io/gorm"
	"vblog-core/model"
)

// PostService handles blog post CRUD operations.
type PostService struct {
	DB     *gorm.DB
	LogSvc *ChangeLogService
}

// NewPostService creates a new PostService.
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db, LogSvc: NewChangeLogService(db)}
}

// CalcReadTime calculates reading time in minutes (~500 chars per minute).
func CalcReadTime(content string) int {
	minutes := int(math.Ceil(float64(len(content)) / 500))
	if minutes < 1 {
		return 1
	}
	return minutes
}

// BuildExcerpt truncates content to maxLen runes with "..." suffix.
func BuildExcerpt(content string, maxLen int) string {
	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}
	return string(runes[:maxLen]) + "..."
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

// GetByID returns a single post by ID with tags preloaded and increments view count.
func (s *PostService) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := s.DB.Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	s.DB.Model(&post).UpdateColumn("views", gorm.Expr("views + 1"))
	post.Views++
	return &post, nil
}

// resolveTags looks up existing tags by name and creates missing ones.
func (s *PostService) resolveTags(tags []model.Tag) ([]model.Tag, error) {
	// Collect valid tag names
	names := make([]string, 0, len(tags))
	for _, t := range tags {
		if t.Name != "" {
			names = append(names, t.Name)
		}
	}
	if len(names) == 0 {
		return nil, nil
	}

	// Batch SELECT existing tags
	var existing []model.Tag
	if err := s.DB.Where("name IN ?", names).Find(&existing).Error; err != nil {
		return nil, err
	}

	// Find missing tags
	existingMap := make(map[string]bool, len(existing))
	for _, t := range existing {
		existingMap[t.Name] = true
	}
	var missing []model.Tag
	for _, t := range tags {
		if t.Name != "" && !existingMap[t.Name] {
			missing = append(missing, t)
		}
	}

	// Batch INSERT missing tags
	if len(missing) > 0 {
		if err := s.DB.CreateInBatches(missing, 100).Error; err != nil {
			return nil, err
		}
		existing = append(existing, missing...)
	}

	// Preserve original order
	tagMap := make(map[string]model.Tag, len(existing))
	for _, t := range existing {
		tagMap[t.Name] = t
	}
	resolved := make([]model.Tag, 0, len(tags))
	for _, t := range tags {
		if t.Name != "" {
			resolved = append(resolved, tagMap[t.Name])
		}
	}
	return resolved, nil
}

// Create creates a new post, auto-calculating read time and excerpt.
func (s *PostService) Create(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	if post.Excerpt == "" {
		post.Excerpt = BuildExcerpt(post.Content, 200)
	}
	if len(post.Tags) > 0 {
		resolved, err := s.resolveTags(post.Tags)
		if err != nil {
			return err
		}
		post.Tags = resolved
	}
	if err := s.DB.Create(post).Error; err != nil {
		return err
	}
	if post.Status == "published" {
		s.LogSvc.Write("new_post", &post.ID, post.Title, "")
	}
	return nil
}

// Update updates an existing post, recalculating read time and syncing tags.
func (s *PostService) Update(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	// Save tags separately to control join table sync.
	tags := post.Tags
	post.Tags = nil
	if err := s.DB.Save(post).Error; err != nil {
		return err
	}
	if len(tags) > 0 {
		resolved, err := s.resolveTags(tags)
		if err != nil {
			return err
		}
		return s.DB.Model(post).Association("Tags").Replace(resolved)
	}
	return nil
}

// Delete soft-deletes a post by ID.
func (s *PostService) Delete(id uint) error {
	var post model.Post
	if err := s.DB.First(&post, id).Error; err != nil {
		return err
	}
	if err := s.DB.Delete(&model.Post{}, id).Error; err != nil {
		return err
	}
	s.LogSvc.Write("delete_post", &id, post.Title, "")
	return nil
}

// ListTrash returns all soft-deleted posts.
func (s *PostService) ListTrash() ([]model.Post, error) {
	var posts []model.Post
	err := s.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Order("deleted_at DESC").Find(&posts).Error
	return posts, err
}

// Restore restores a soft-deleted post.
func (s *PostService) Restore(id uint) error {
	return s.DB.Unscoped().Model(&model.Post{}).
		Where("id = ?", id).Update("deleted_at", nil).Error
}

// PermanentDelete hard-deletes a post and its associations.
func (s *PostService) PermanentDelete(id uint) error {
	// Remove join table associations first (foreign key)
	if err := s.DB.Exec("DELETE FROM post_tags WHERE post_id = ?", id).Error; err != nil {
		return err
	}
	// Delete related comments
	s.DB.Unscoped().Where("post_id = ?", id).Delete(&model.Comment{})
	// Delete related change_log entries
	s.DB.Where("target_id = ?", id).Delete(&model.ChangeLog{})
	// Hard delete the post
	return s.DB.Unscoped().Delete(&model.Post{}, id).Error
}
