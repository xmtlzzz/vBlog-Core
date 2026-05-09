package model

import "testing"

func TestPostTableName(t *testing.T) {
	p := Post{}
	if p.TableName() != "posts" {
		t.Errorf("expected 'posts', got '%s'", p.TableName())
	}
}

func TestPostDefaultValues(t *testing.T) {
	p := Post{}
	if p.Status != "" {
		t.Errorf("expected default status '', got '%s'", p.Status)
	}
	if p.Pinned != false {
		t.Errorf("expected default pinned false, got %v", p.Pinned)
	}
	if p.Views != 0 {
		t.Errorf("expected default views 0, got %d", p.Views)
	}
	if p.ReadTime != 0 {
		t.Errorf("expected default read_time 0, got %d", p.ReadTime)
	}
}
