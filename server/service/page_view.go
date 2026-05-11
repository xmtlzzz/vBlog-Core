package service

import (
	"time"

	"vblog-core/model"

	"gorm.io/gorm"
)

// PageViewService handles page view recording and analytics.
type PageViewService struct {
	DB *gorm.DB
}

// NewPageViewService creates a new PageViewService.
func NewPageViewService(db *gorm.DB) *PageViewService {
	return &PageViewService{DB: db}
}

// Record creates a new page view entry.
func (s *PageViewService) Record(ip, path, userAgent string) error {
	return s.DB.Create(&model.PageView{
		IP:        ip,
		Path:      path,
		UserAgent: userAgent,
	}).Error
}

// GetPVUVToday returns today's page views and unique visitors.
func (s *PageViewService) GetPVUVToday() (pv int64, uv int64, err error) {
	today := time.Now().Format("2006-01-02")
	return s.GetPVUVByDate(today)
}

// GetPVUVByDate returns page views and unique visitors for a given date.
func (s *PageViewService) GetPVUVByDate(date string) (pv int64, uv int64, err error) {
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, 0, err
	}
	end := start.Add(24 * time.Hour)
	err = s.DB.Model(&model.PageView{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Select("COUNT(*), COUNT(DISTINCT ip)").
		Row().Scan(&pv, &uv)
	return
}
