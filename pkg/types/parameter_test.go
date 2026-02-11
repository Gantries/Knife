package types

import (
	"testing"
	"time"
)

func TestParseTimeInLocation(t *testing.T) {
	location, _ = time.LoadLocation("Local")
	tests := []struct {
		name      string
		timeStr   string
		expected  time.Time
		expectErr bool
	}{
		{
			name:      "normal",
			timeStr:   "2024-05-23 12:00:00",
			expected:  time.Date(2024, 5, 23, 12, 0, 0, 0, location),
			expectErr: false,
		},
		{
			name:      "empty_string",
			timeStr:   "",
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name:      "invalid_format",
			timeStr:   "not-a-time",
			expected:  time.Time{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTime, err := ParseTimeInLocation(tt.timeStr)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for test %q, but got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for test %q: %v", tt.name, err)
					return
				}
				if !parsedTime.Equal(tt.expected) {
					t.Errorf("Expected %v, got %v for test %q", tt.expected, parsedTime, tt.name)
				}
			}
		})
	}
}
