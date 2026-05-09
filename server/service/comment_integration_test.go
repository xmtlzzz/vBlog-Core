package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestCommentService_Create(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewCommentService(db)

	// First create a post for the comment
	post := &model.Post{Title: "Comment Test Post", Content: "Content", Status: "published"}
	db.Create(post)

	comment := &model.Comment{
		PostID:     post.ID,
		AuthorName: "TestUser",
		Body:       "This is a test comment",
	}
	err := svc.Create(comment)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if comment.Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", comment.Status)
	}
	if comment.ID == 0 {
		t.Fatal("expected comment ID to be set")
	}

	db.Unscoped().Delete(comment)
	db.Unscoped().Delete(post)
}

func TestCommentService_List(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewCommentService(db)

	post := &model.Post{Title: "CommentList Post", Content: "Content", Status: "published"}
	db.Create(post)

	svc.Create(&model.Comment{PostID: post.ID, AuthorName: "User1", Body: "Comment 1"})
	svc.Create(&model.Comment{PostID: post.ID, AuthorName: "User2", Body: "Comment 2"})

	comments, total, err := svc.List(1, 10, "", "")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total < 2 {
		t.Errorf("expected total >= 2, got %d", total)
	}
	if len(comments) < 2 {
		t.Errorf("expected >= 2 comments, got %d", len(comments))
	}

	db.Unscoped().Where("post_id = ?", post.ID).Delete(&model.Comment{})
	db.Unscoped().Delete(post)
}

func TestCommentService_Approve(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewCommentService(db)

	post := &model.Post{Title: "Approve Post", Content: "Content", Status: "published"}
	db.Create(post)

	comment := &model.Comment{PostID: post.ID, AuthorName: "User", Body: "To approve"}
	svc.Create(comment)

	err := svc.Approve(comment.ID)
	if err != nil {
		t.Fatalf("Approve failed: %v", err)
	}

	var found model.Comment
	db.First(&found, comment.ID)
	if found.Status != "approved" {
		t.Errorf("expected status 'approved', got '%s'", found.Status)
	}

	db.Unscoped().Delete(comment)
	db.Unscoped().Delete(post)
}

func TestCommentService_MarkSpam(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewCommentService(db)

	post := &model.Post{Title: "Spam Post", Content: "Content", Status: "published"}
	db.Create(post)

	comment := &model.Comment{PostID: post.ID, AuthorName: "Spammer", Body: "Spam content"}
	svc.Create(comment)

	err := svc.MarkSpam(comment.ID)
	if err != nil {
		t.Fatalf("MarkSpam failed: %v", err)
	}

	var found model.Comment
	db.First(&found, comment.ID)
	if found.Status != "spam" {
		t.Errorf("expected status 'spam', got '%s'", found.Status)
	}

	db.Unscoped().Delete(comment)
	db.Unscoped().Delete(post)
}

func TestCommentService_Delete(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewCommentService(db)

	post := &model.Post{Title: "Delete Comment Post", Content: "Content", Status: "published"}
	db.Create(post)

	comment := &model.Comment{PostID: post.ID, AuthorName: "User", Body: "To delete"}
	svc.Create(comment)

	err := svc.Delete(comment.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	var count int64
	db.Model(&model.Comment{}).Where("id = ?", comment.ID).Count(&count)
	if count != 0 {
		t.Error("expected comment to be deleted")
	}

	db.Unscoped().Delete(comment)
	db.Unscoped().Delete(post)
}
