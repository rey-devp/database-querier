package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Log = slog.New(handler)
}

// Info logs an informational message with a specific component tag.
func Info(tag string, msg string, args ...any) {
	allArgs := append([]any{"tag", tag}, args...)
	Log.Info(msg, allArgs...)
}

// Warn logs a warning message with a specific component tag.
func Warn(tag string, msg string, args ...any) {
	allArgs := append([]any{"tag", tag}, args...)
	Log.Warn(msg, allArgs...)
}

// Error logs an error message with a specific component tag.
func Error(tag string, msg string, args ...any) {
	allArgs := append([]any{"tag", tag}, args...)
	Log.Error(msg, allArgs...)
}
