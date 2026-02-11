package tel

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestMultiHandler_Enabled(t *testing.T) {
	h := MultiHandler{
		name:     "test",
		handlers: []slog.Handler{slog.NewTextHandler(os.Stdout, nil)},
	}

	result := h.Enabled(context.Background(), slog.LevelInfo)
	if !result {
		t.Error("MultiHandler.Enabled() returned false")
	}
}

func TestMultiHandler_Handle(t *testing.T) {
	h := MultiHandler{
		name:     "test",
		handlers: []slog.Handler{slog.NewTextHandler(os.Stdout, nil)},
	}

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
	err := h.Handle(context.Background(), record)
	// Should not error even if handlers fail
	if err != nil {
		t.Errorf("MultiHandler.Handle() unexpected error: %v", err)
	}
}

func TestMultiHandler_WithAttrs(t *testing.T) {
	h := MultiHandler{
		name:     "test",
		handlers: []slog.Handler{slog.NewTextHandler(os.Stdout, nil)},
	}

	attrs := []slog.Attr{slog.String("key", "value")}
	newHandler := h.WithAttrs(attrs)
	if newHandler == nil {
		t.Error("MultiHandler.WithAttrs() returned nil")
	}
}

func TestMultiHandler_WithGroup(t *testing.T) {
	h := MultiHandler{
		name:     "test",
		handlers: []slog.Handler{slog.NewTextHandler(os.Stdout, nil)},
	}

	newHandler := h.WithGroup("group")
	if newHandler == nil {
		t.Error("MultiHandler.WithGroup() returned nil")
	}
}

func TestLogger(t *testing.T) {
	logger := Logger("test/logger")
	if logger == nil {
		t.Error("Logger() returned nil")
	}
}

func TestMultipleLoggers(t *testing.T) {
	loggers := []string{
		"test/logger1",
		"test/logger2",
		"test/logger3",
	}

	for _, name := range loggers {
		logger := Logger(name)
		if logger == nil {
			t.Errorf("Logger(%q) returned nil", name)
		}
	}
}
