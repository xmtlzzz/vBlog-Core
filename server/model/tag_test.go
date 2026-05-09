package model

import "testing"

func TestTagTableName(t *testing.T) {
	tag := Tag{}
	if tag.TableName() != "tags" {
		t.Errorf("expected 'tags', got '%s'", tag.TableName())
	}
}

func TestTagDefaultValues(t *testing.T) {
	tag := Tag{}
	if tag.Name != "" {
		t.Errorf("expected default name '', got '%s'", tag.Name)
	}
	if tag.Description != "" {
		t.Errorf("expected default description '', got '%s'", tag.Description)
	}
}
