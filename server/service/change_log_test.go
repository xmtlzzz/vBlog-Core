package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestChangeLogService_Write(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	err := svc.Write("new_post", nil, "Test Post", `{"id":1}`)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var log model.ChangeLog
	err = db.Order("id DESC").First(&log).Error
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if log.ChangeType != "new_post" {
		t.Errorf("expected 'new_post', got '%s'", log.ChangeType)
	}
	if log.Title != "Test Post" {
		t.Errorf("expected 'Test Post', got '%s'", log.Title)
	}

	db.Where("change_type = ?", "new_post").Delete(&model.ChangeLog{})
}

func TestChangeLogService_GetAfterID(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	svc.Write("new_post", nil, "Post 1", "")
	svc.Write("new_comment", nil, "Comment 1", "")

	var last model.ChangeLog
	db.Order("id DESC").First(&last)

	logs, err := svc.GetAfterID(int64(last.ID) - 1)
	if err != nil {
		t.Fatalf("GetAfterID failed: %v", err)
	}
	if len(logs) < 1 {
		t.Errorf("expected >= 1 log, got %d", len(logs))
	}

	db.Where("1 = 1").Delete(&model.ChangeLog{})
}

func TestChangeLogService_GetLatestID(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	svc.Write("test_type", nil, "test", "")

	id, err := svc.GetLatestID()
	if err != nil {
		t.Fatalf("GetLatestID failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero ID")
	}

	db.Where("change_type = ?", "test_type").Delete(&model.ChangeLog{})
}
