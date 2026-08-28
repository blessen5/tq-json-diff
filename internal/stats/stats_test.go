package stats

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{1024 * 1024 * 5, "5.0 MB"},
		{1024 * 1024 * 1024 * 2, "2.00 GB"},
	}

	for _, tt := range tests {
		res := FormatBytes(tt.bytes)
		if res != tt.expected {
			t.Errorf("FormatBytes(%d) = %q, expected %q", tt.bytes, res, tt.expected)
		}
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasErr   bool
	}{
		{"500B", 500, false},
		{"500", 500, false},
		{"1KB", 1024, false},
		{"2K", 2048, false},
		{"10MB", 10 * 1024 * 1024, false},
		{"1.5MB", int64(1.5 * 1024 * 1024), false},
		{"1GB", 1024 * 1024 * 1024, false},
		{"", 0, true},
		{"invalid", 0, true},
		{"-50MB", 0, true},
	}

	for _, tt := range tests {
		res, err := ParseSize(tt.input)
		if tt.hasErr {
			if err == nil {
				t.Errorf("ParseSize(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseSize(%q) unexpected error: %v", tt.input, err)
			}
			if res != tt.expected {
				t.Errorf("ParseSize(%q) = %d, expected %d", tt.input, res, tt.expected)
			}
		}
	}
}
