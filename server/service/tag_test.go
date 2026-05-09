package service

import "testing"

func TestNewTagService(t *testing.T) {
	s := NewTagService(nil)
	if s == nil {
		t.Fatal("NewTagService returned nil")
	}
	if s.DB != nil {
		t.Error("expected DB to be nil")
	}
}

func TestTagServiceStruct(t *testing.T) {
	s := &TagService{}
	if s.DB != nil {
		t.Error("expected zero-value DB to be nil")
	}
}
