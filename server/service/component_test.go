package service

import "testing"

func TestNewComponentService(t *testing.T) {
	s := NewComponentService(nil)
	if s == nil {
		t.Fatal("NewComponentService returned nil")
	}
	if s.DB != nil {
		t.Error("expected DB to be nil")
	}
}

func TestComponentServiceStruct(t *testing.T) {
	s := &ComponentService{}
	if s.DB != nil {
		t.Error("expected zero-value DB to be nil")
	}
}

func TestComponentValidation(t *testing.T) {
	tests := []struct {
		name      string
		component struct {
			Name   string
			Status string
			Origin string
		}
		validName   bool
		validStatus bool
	}{
		{
			name:        "valid component",
			component:   struct{ Name, Status, Origin string }{"header", "active", "uploaded"},
			validName:   true,
			validStatus: true,
		},
		{
			name:        "empty name is invalid",
			component:   struct{ Name, Status, Origin string }{"", "active", "uploaded"},
			validName:   false,
			validStatus: true,
		},
		{
			name:        "valid inactive status",
			component:   struct{ Name, Status, Origin string }{"sidebar", "inactive", "uploaded"},
			validName:   true,
			validStatus: true,
		},
		{
			name:        "unknown status is invalid",
			component:   struct{ Name, Status, Origin string }{"footer", "deleted", "uploaded"},
			validName:   true,
			validStatus: false,
		},
	}

	validStatuses := map[string]bool{"active": true, "inactive": true}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameValid := tt.component.Name != ""
			if nameValid != tt.validName {
				t.Errorf("name validity: expected %v, got %v", tt.validName, nameValid)
			}
			statusValid := validStatuses[tt.component.Status]
			if statusValid != tt.validStatus {
				t.Errorf("status validity: expected %v, got %v", tt.validStatus, statusValid)
			}
		})
	}
}
