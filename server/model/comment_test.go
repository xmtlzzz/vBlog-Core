package model

import "testing"

func TestCommentTableName(t *testing.T) {
	c := Comment{}
	if c.TableName() != "comments" {
		t.Errorf("expected 'comments', got '%s'", c.TableName())
	}
}

func TestCommentDefaultValues(t *testing.T) {
	c := Comment{}
	if c.Status != "" {
		t.Errorf("expected default status '', got '%s'", c.Status)
	}
	if c.PostID != 0 {
		t.Errorf("expected default post_id 0, got %d", c.PostID)
	}
}
