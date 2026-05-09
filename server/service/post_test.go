package service

import "testing"

func TestCalcReadTime(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"empty content returns 1", "", 1},
		{"short content returns 1", "hello", 1},
		{"exactly 500 chars returns 1", string(make([]byte, 500)), 1},
		{"501 chars returns 2", string(make([]byte, 501)), 2},
		{"1000 chars returns 2", string(make([]byte, 1000)), 2},
		{"1500 chars returns 3", string(make([]byte, 1500)), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcReadTime(tt.content)
			if got != tt.want {
				t.Errorf("CalcReadTime(len=%d) = %d, want %d", len(tt.content), got, tt.want)
			}
		})
	}
}

func TestBuildExcerpt(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short content returned as-is", "hello", 200, "hello"},
		{"exact maxLen returned as-is", "hello", 5, "hello"},
		{"long content truncated with ellipsis", "hello world this is a long string", 5, "hello..."},
		{"empty content", "", 10, ""},
		{"maxLen 0 always truncates", "abc", 0, "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildExcerpt(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("BuildExcerpt(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}
