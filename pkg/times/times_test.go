package times

import (
	"testing"
	"time"
)

func TestFormatTs(t *testing.T) {
	// Test that the function runs without error and returns a non-empty string
	tests := []struct {
		name string
		ts   int64
	}{
		{
			name: "timestamp zero",
			ts:   0,
		},
		{
			name: "timestamp 1 billion",
			ts:   1000000000,
		},
		{
			name: "current time",
			ts:   time.Now().UnixMilli(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTs(tt.ts)
			if result == "" {
				t.Error("FormatTs() returned empty string")
			}
			// Verify it matches expected format length (19 chars for "2006-01-02 15:04:05")
			if len(result) != 19 {
				t.Errorf("FormatTs() returned length %d, expected 19", len(result))
			}
		})
	}
}

func TestFormatTsByLayout(t *testing.T) {
	tests := []struct {
		name   string
		ts     int64
		layout string
		check  func(string) bool
	}{
		{
			name:   "default layout",
			ts:     1000000000,
			layout: "2006-01-02 15:04:05",
			check: func(s string) bool {
				return len(s) == 19
			},
		},
		{
			name:   "year only",
			ts:     1000000000,
			layout: "2006",
			check: func(s string) bool {
				return len(s) == 4
			},
		},
		{
			name:   "custom layout",
			ts:     1000000000,
			layout: "2006/01/02",
			check: func(s string) bool {
				return len(s) == 10
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTsByLayout(tt.ts, tt.layout)
			if result == "" {
				t.Error("FormatTsByLayout() returned empty string")
			}
			if !tt.check(result) {
				t.Errorf("FormatTsByLayout() = %v, check failed", result)
			}
		})
	}
}

func TestFormatTsConsistency(t *testing.T) {
	ts := int64(1000000000)
	result1 := FormatTs(ts)
	result2 := FormatTs(ts)

	if result1 != result2 {
		t.Errorf("FormatTs() inconsistent results: %v != %v", result1, result2)
	}
}
