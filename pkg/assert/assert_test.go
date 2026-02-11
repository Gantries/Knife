package assert

import (
	"testing"
)

type mockTC struct {
	errors []string
}

func (m *mockTC) Errorf(format string, args ...any) {
	m.errors = append(m.errors, format)
}

func TestTrue(t *testing.T) {
	tests := []struct {
		name       string
		expression bool
		fmt        string
		args       []any
		wantError  bool
	}{
		{
			name:       "true expression",
			expression: true,
			fmt:        "should not error",
			wantError:  false,
		},
		{
			name:       "false expression",
			expression: false,
			fmt:        "expected error: %s",
			args:       []any{"test"},
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockTC{}
			True(m, tt.expression, tt.fmt, tt.args...)
			if tt.wantError && len(m.errors) == 0 {
				t.Errorf("True() expected error but got none")
			}
			if !tt.wantError && len(m.errors) > 0 {
				t.Errorf("True() unexpected error: %v", m.errors)
			}
		})
	}
}
