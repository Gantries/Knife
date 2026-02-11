package maps

import (
	"testing"
)

func TestSetOf(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "empty set",
			args: []string{},
			want: []string{},
		},
		{
			name: "single element",
			args: []string{"a"},
			want: []string{"a"},
		},
		{
			name: "multiple elements",
			args: []string{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SetOf(tt.args...)
			for _, v := range tt.want {
				if !s.Has(v) {
					t.Errorf("SetOf().Has(%v) = false, want true", v)
				}
			}
		})
	}
}

func TestSet_Has(t *testing.T) {
	s := SetOf("a", "b", "c")

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "existing key",
			key:  "a",
			want: true,
		},
		{
			name: "non-existing key",
			key:  "z",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Has(tt.key); got != tt.want {
				t.Errorf("Set.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}
