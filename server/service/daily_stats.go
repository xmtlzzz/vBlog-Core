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

	var postCount int64
	s.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)

	var viewTotal int64
	s.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)

	var commentCount int64
	s.DB.Model(&model.Comment{}).Count(&commentCount)

	var tagCount int64
	s.DB.Model(&model.Tag{}).Count(&tagCount)

	var pvToday, uvToday int64
	s.DB.Model(&model.PageView{}).Where("DATE(created_at) = ?", today).Count(&pvToday)
	s.DB.Model(&model.PageView{}).Where("DATE(created_at) = ?", today).Distinct("ip").Count(&uvToday)

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
