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
	resolved := make([]model.Tag, 0, len(tags))
	for _, t := range tags {
		if t.Name == "" {
			continue
		}
		var existing model.Tag
		err := s.DB.Where("name = ?", t.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.DB.Create(&t).Error; err != nil {
				return nil, err
			}
			resolved = append(resolved, t)
		} else if err != nil {
			return nil, err
		} else {
			resolved = append(resolved, existing)
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
