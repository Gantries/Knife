package log

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test logger",
		},
		{
			name: "another/logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.name)
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}
