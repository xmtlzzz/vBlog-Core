package service

import (
	"testing"

	"vblog-core/model"
	"vblog-core/testutil"
)

func TestPostService_Create(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPostService(db)

	post := &model.Post{
		Title:   "Test Post",
		Content: "This is test content for integration testing.",
		Status:  "published",
	}

	err := svc.Create(post)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if post.ID == 0 {
		t.Fatal("expected post ID to be set")
	}
	if post.ReadTime < 1 {
		t.Errorf("expected ReadTime >= 1, got %d", post.ReadTime)
	}
	if post.Excerpt == "" {
		t.Error("expected Excerpt to be auto-generated")
	}

	// Cleanup
	db.Unscoped().Delete(post)
}

func TestPostService_List(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPostService(db)

	// Create test posts
	for i := 0; i < 3; i++ {
		svc.Create(&model.Post{
			Title:   "List Test Post",
			Content: "Content",
			Status:  "published",
		})
	}

	posts, total, err := svc.List(1, 10, "", "published", "")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total < 3 {
		t.Errorf("expected total >= 3, got %d", total)
	}
	if len(posts) < 3 {
		t.Errorf("expected >= 3 posts, got %d", len(posts))
	}

	// Cleanup
	db.Unscoped().Where("title = ?", "List Test Post").Delete(&model.Post{})
}

func TestPostService_GetByID(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPostService(db)

	post := &model.Post{Title: "Get Test", Content: "Content", Status: "draft"}
	svc.Create(post)

	found, err := svc.GetByID(post.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.Title != "Get Test" {
		t.Errorf("expected title 'Get Test', got '%s'", found.Title)
	}

	db.Unscoped().Delete(post)
}

func TestPostService_Update(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPostService(db)

	post := &model.Post{Title: "Original", Content: "Short", Status: "draft"}
	svc.Create(post)

	post.Title = "Updated"
	post.Content = "This is a much longer content that should change the read time calculation significantly."
	err := svc.Update(post)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := svc.GetByID(post.ID)
	if found.Title != "Updated" {
		t.Errorf("expected title 'Updated', got '%s'", found.Title)
	}

	db.Unscoped().Delete(post)
}

func TestPostService_Delete(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPostService(db)

	post := &model.Post{Title: "To Delete", Content: "Content", Status: "draft"}
	svc.Create(post)

	err := svc.Delete(post.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Soft deleted -- should not be found by normal query
	_, err = svc.GetByID(post.ID)
	if err == nil {
		t.Error("expected error for soft-deleted post")
	}

	db.Unscoped().Delete(post)
}
