package model

import "testing"

func TestComponentTableName(t *testing.T) {
	c := Component{}
	if c.TableName() != "components" {
		t.Errorf("expected 'components', got '%s'", c.TableName())
	}
}

func TestComponentDefaultValues(t *testing.T) {
	c := Component{}
	if c.Status != "" {
		t.Errorf("expected default status '', got '%s'", c.Status)
	}
	if c.Origin != "" {
		t.Errorf("expected default origin '', got '%s'", c.Origin)
	}
	if c.Version != "" {
		t.Errorf("expected default version '', got '%s'", c.Version)
	}
}
