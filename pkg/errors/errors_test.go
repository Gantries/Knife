package errors

import (
	"errors"
	"testing"
)

func TestYes(t *testing.T) {
	called := false
	Yes(func() { called = true })
	if !called {
		t.Error("Yes() did not call the provided function")
	}
}

func TestYesMultiple(t *testing.T) {
	count := 0
	Yes(
		func() { count++ },
		func() { count++ },
		func() { count++ },
	)
	if count != 3 {
		t.Errorf("Yes() called %d functions, want 3", count)
	}
}

func TestNo(t *testing.T) {
	err := errors.New("test error")
	msg := No(err)

	if msg == nil {
		t.Error("No() returned nil message")
	}
}

func TestNoWithArgs(t *testing.T) {
	err := errors.New("test error with {{.value}}")
	msg := No(err, "value", "test")

	if msg == nil {
		t.Error("No() with args returned nil message")
	}
}
