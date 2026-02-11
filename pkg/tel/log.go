// Package tel provides OpenTelemetry integration for observability.
//
// This file implements logging functionality with multi-handler support,
// combining standard slog handlers with OpenTelemetry bridges.
package tel

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
)

type MultiHandler struct {
	name     string
	handlers []slog.Handler
}

func (h MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range h.handlers {
		h.Enabled(ctx, level)
	}
	return true
}

func (h MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var err error
	for _, h := range h.handlers {
		err = errors.Join(err, h.Handle(ctx, record))
	}
	return err
}

func (h MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var handler slog.Handler
	for _, h := range h.handlers {
		handler = h.WithAttrs(attrs)
	}
	return handler
}

func (h MultiHandler) WithGroup(name string) slog.Handler {
	var handler slog.Handler
	for _, h := range h.handlers {
		handler = h.WithGroup(name)
	}
	return handler
}

var logger = Logger("knife/tel/log")

var partialInitializedHandlers = make([]*MultiHandler, 0)

var defaultLogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})

func Logger(name string) *slog.Logger {
	var mh *MultiHandler
	if loggerProvider == nil {
		mh = &MultiHandler{name: name, handlers: []slog.Handler{defaultLogHandler}}
		partialInitializedHandlers = append(partialInitializedHandlers, mh)
	} else {
		mh = &MultiHandler{name: name, handlers: []slog.Handler{defaultLogHandler, otelslog.NewHandler(name, otelslog.WithLoggerProvider(loggerProvider))}}
	}
	return slog.New(mh)

}

func SetupLoggersCreated(p log.LoggerProvider) {
	if loggerProvider == nil {
		loggerProvider = p
	}
	for _, l := range partialInitializedHandlers {
		logger.Info("Logger has been fully initialized", "logger", l.name)
		l.handlers = append(l.handlers, otelslog.NewHandler(l.name, otelslog.WithLoggerProvider(loggerProvider)))
	}
	partialInitializedHandlers = partialInitializedHandlers[:0]
}
