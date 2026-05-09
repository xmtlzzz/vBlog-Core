package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestTagService_Create(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewTagService(db)

	tag := &model.Tag{Name: "TestTag", Description: "A test tag"}
	err := svc.Create(tag)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if tag.ID == 0 {
		t.Fatal("expected tag ID to be set")
	}
	db.Unscoped().Delete(tag)
}

func TestTagService_List(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewTagService(db)

	svc.Create(&model.Tag{Name: "ListTag1"})
	svc.Create(&model.Tag{Name: "ListTag2"})

	tags, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(tags) < 2 {
		t.Errorf("expected >= 2 tags, got %d", len(tags))
	}

	db.Unscoped().Where("name LIKE ?", "ListTag%").Delete(&model.Tag{})
}

func TestTagService_Update(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewTagService(db)

	tag := &model.Tag{Name: "UpdateTag", Description: "Original"}
	svc.Create(tag)

	tag.Description = "Updated description"
	err := svc.Update(tag)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	var found model.Tag
	db.First(&found, tag.ID)
	if found.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%s'", found.Description)
	}
	db.Unscoped().Delete(tag)
}

func TestTagService_Delete(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewTagService(db)

	tag := &model.Tag{Name: "DeleteTag"}
	svc.Create(tag)

	err := svc.Delete(tag.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	var count int64
	db.Model(&model.Tag{}).Where("id = ?", tag.ID).Count(&count)
	if count != 0 {
		t.Error("expected tag to be deleted")
	}
	db.Unscoped().Delete(tag)
}
