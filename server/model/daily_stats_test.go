package model

import "testing"

func TestDailyStatsTableName(t *testing.T) {
	d := DailyStats{}
	if d.TableName() != "daily_stats" {
		t.Errorf("expected 'daily_stats', got '%s'", d.TableName())
	}
}

func TestDailyStatsDefaultValues(t *testing.T) {
	d := DailyStats{}
	if d.PV != 0 {
		t.Errorf("expected PV 0, got %d", d.PV)
	}
	if d.UV != 0 {
		t.Errorf("expected UV 0, got %d", d.UV)
	}
}
