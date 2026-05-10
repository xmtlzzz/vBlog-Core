package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestPageViewService_Record(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPageViewService(db)

	err := svc.Record("192.168.1.1", "/posts/1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	var count int64
	db.Model(&model.PageView{}).Where("ip = ? AND path = ?", "192.168.1.1", "/posts/1").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 page view, got %d", count)
	}

	db.Where("ip = ?", "192.168.1.1").Delete(&model.PageView{})
}

func TestPageViewService_GetPVUV(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPageViewService(db)

	svc.Record("10.0.0.1", "/posts/1", "ua")
	svc.Record("10.0.0.2", "/posts/1", "ua")
	svc.Record("10.0.0.1", "/posts/2", "ua")

	pv, uv, err := svc.GetPVUVToday()
	if err != nil {
		t.Fatalf("GetPVUVToday failed: %v", err)
	}
	if pv < 3 {
		t.Errorf("expected PV >= 3, got %d", pv)
	}
	if uv < 2 {
		t.Errorf("expected UV >= 2, got %d", uv)
	}

	db.Where("ip LIKE ?", "10.0.0.%").Delete(&model.PageView{})
}
