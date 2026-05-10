package service

import (
	"testing"
	"time"

	"vblog-core/model"
	"vblog-core/testutil"
)

func TestDailyStatsService_Snapshot(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewDailyStatsService(db)

	err := svc.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	var stats model.DailyStats
	err = db.Where("stat_date = ?", today).First(&stats).Error
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	db.Where("stat_date = ?", today).Delete(&model.DailyStats{})
}

func TestDailyStatsService_GetTrends(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewDailyStatsService(db)

	now := time.Now()
	for i := 0; i < 5; i++ {
		date := now.AddDate(0, 0, -i)
		db.Create(&model.DailyStats{
			StatDate:  date,
			PV:        int64(100 + i),
			UV:        int64(50 + i),
			ViewTotal: int64(1000 + i*10),
		})
	}

	points, err := svc.GetTrends("day", 5)
	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}
	if len(points) < 1 {
		t.Errorf("expected >= 1 point, got %d", len(points))
	}

	db.Where("pv >= ?", 100).Delete(&model.DailyStats{})
}
