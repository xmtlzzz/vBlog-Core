package service

import (
	"time"

	"vblog-core/model"

	"gorm.io/gorm"
)

// DailyStatsService handles daily statistics aggregation.
type DailyStatsService struct {
	DB *gorm.DB
}

// NewDailyStatsService creates a new DailyStatsService.
func NewDailyStatsService(db *gorm.DB) *DailyStatsService {
	return &DailyStatsService{DB: db}
}

// Snapshot aggregates current stats from multiple tables and upserts into daily_stats.
func (s *DailyStatsService) Snapshot() error {
	today := time.Now().Format("2006-01-02")

	var postCount, viewTotal int64
	s.DB.Model(&model.Post{}).
		Select("COUNT(*) FILTER (WHERE status = 'published'), COALESCE(SUM(views), 0)").
		Row().Scan(&postCount, &viewTotal)

	var commentCount, tagCount int64
	s.DB.Raw(`SELECT
		(SELECT COUNT(*) FROM comments),
		(SELECT COUNT(*) FROM tags)`).Row().Scan(&commentCount, &tagCount)

	var pvToday, uvToday int64
	s.DB.Model(&model.PageView{}).
		Where("DATE(created_at) = ?", today).
		Select("COUNT(*), COUNT(DISTINCT ip)").
		Row().Scan(&pvToday, &uvToday)

	return s.DB.Where("stat_date = ?", today).
		Assign(model.DailyStats{
			PV:           pvToday,
			UV:           uvToday,
			PostCount:    int(postCount),
			ViewTotal:    viewTotal,
			CommentCount: int(commentCount),
			TagCount:     int(tagCount),
		}).FirstOrCreate(&model.DailyStats{StatDate: time.Now()}).Error
}

// GetTrends returns recent daily stats ordered by date descending.
func (s *DailyStatsService) GetTrends(granularity string, count int) ([]model.DailyStats, error) {
	var stats []model.DailyStats
	err := s.DB.Order("stat_date DESC").Limit(count).Find(&stats).Error
	return stats, err
}
