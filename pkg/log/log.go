// Package log provides a simple logging facade built on top of the tel package.
//
// It creates structured loggers using the standard library's log/slog with
// OpenTelemetry integration provided by the tel package.
package log

import (
	"github.com/gantries/knife/pkg/tel"

	"log/slog"
)

func New(name string) *slog.Logger {
	return tel.Logger(name)
}
