// Package assert provides simple assertion utilities for testing.
//
// It offers a minimal True function for test assertions with formatted error messages.
package assert

type tc interface {
	Errorf(format string, args ...any)
}

// True assert expression and display optional message.
func True(t tc, exp bool, fmt string, args ...interface{}) {
	if !exp {
		t.Errorf(fmt, args...)
	}
}
