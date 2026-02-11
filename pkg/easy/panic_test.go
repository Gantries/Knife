package easy

import (
	"errors"
	"testing"
)

func TestPanic(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		cleanup   func(error)
		wantPanic bool
	}{
		{
			name:      "no error",
			err:       nil,
			wantPanic: false,
		},
		{
			name:      "error without cleanup",
			err:       errors.New("test error"),
			wantPanic: true,
		},
		{
			name: "error with cleanup",
			err:  errors.New("test error"),
			cleanup: func(error) {
				// cleanup called
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Panic() expected panic but got none")
					}
				}()
			}
			Panic(tt.err, tt.cleanup)
		})
	}
}

func TestPanicE(t *testing.T) {
	tests := []struct {
		name      string
		v         *string
		err       error
		wantPanic bool
	}{
		{
			name:      "no error",
			v:         ptr("test"),
			err:       nil,
			wantPanic: false,
		},
		{
			name:      "error with nil value",
			v:         nil,
			err:       errors.New("test error"),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("PanicE() expected panic but got none")
					}
				}()
			}
			PanicE(tt.v, tt.err)
		})
	}
}

func TestPanicN(t *testing.T) {
	tests := []struct {
		name      string
		v         *string
		cleanup   func(error)
		wantPanic bool
	}{
		{
			name:      "non-nil value",
			v:         ptr("test"),
			wantPanic: false,
		},
		{
			name:      "nil value",
			v:         nil,
			wantPanic: true,
		},
		{
			name: "nil value with cleanup",
			v:    nil,
			cleanup: func(error) {
				// cleanup called
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("PanicN() expected panic but got none")
					}
				}()
			}
			PanicN(tt.v, tt.cleanup)
		})
	}
}

func ptr(s string) *string {
	return &s
}
