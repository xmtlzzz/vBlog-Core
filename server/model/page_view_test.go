package model

import "testing"

func TestPageViewTableName(t *testing.T) {
	p := PageView{}
	if p.TableName() != "page_views" {
		t.Errorf("expected 'page_views', got '%s'", p.TableName())
	}
}

func TestPageViewDefaultValues(t *testing.T) {
	p := PageView{}
	if p.IP != "" {
		t.Errorf("expected empty IP, got '%s'", p.IP)
	}
	if p.Path != "" {
		t.Errorf("expected empty Path, got '%s'", p.Path)
	}
}
