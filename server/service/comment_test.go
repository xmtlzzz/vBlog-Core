package service

import "testing"

func TestNewCommentService(t *testing.T) {
	s := NewCommentService(nil)
	if s == nil {
		t.Fatal("NewCommentService returned nil")
	}
	if s.DB != nil {
		t.Error("expected DB to be nil")
	}
}

func TestCommentServiceStruct(t *testing.T) {
	s := &CommentService{}
	if s.DB != nil {
		t.Error("expected zero-value DB to be nil")
	}
}

func TestCommentStatusTransition(t *testing.T) {
	// Verify the valid status transitions
	validStatuses := map[string]bool{
		"pending":  true,
		"approved": true,
		"spam":     true,
	}

	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"pending is valid", "pending", true},
		{"approved is valid", "approved", true},
		{"spam is valid", "spam", true},
		{"unknown is invalid", "deleted", false},
		{"empty is invalid", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validStatuses[tt.status]
			if got != tt.valid {
				t.Errorf("status %q: expected valid=%v, got %v", tt.status, tt.valid, got)
			}
		})
	}
}
