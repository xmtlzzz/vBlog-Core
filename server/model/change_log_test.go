package model

import "testing"

func TestChangeLogTableName(t *testing.T) {
	c := ChangeLog{}
	if c.TableName() != "change_log" {
		t.Errorf("expected 'change_log', got '%s'", c.TableName())
	}
}

func TestChangeLogDefaultValues(t *testing.T) {
	c := ChangeLog{}
	if c.ChangeType != "" {
		t.Errorf("expected empty change_type, got '%s'", c.ChangeType)
	}
}
