// Package logger is the single sanctioned logging shim (pkg/logger only —
// never fmt.Println/log.Printf in services).
package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const traceIDKey ctxKey = "trace_id"

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

// WithTrace attaches a trace id to the context for structured logs.
func WithTrace(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// Info logs at info with any trace id found in ctx.
func Info(ctx context.Context, msg string, args ...any) {
	if t, ok := ctx.Value(traceIDKey).(string); ok && t != "" {
		args = append(args, "trace_id", t)
	}
	defaultLogger.InfoContext(ctx, msg, args...)
}

// Error logs at error with any trace id found in ctx.
func Error(ctx context.Context, msg string, args ...any) {
	if t, ok := ctx.Value(traceIDKey).(string); ok && t != "" {
		args = append(args, "trace_id", t)
	}
	defaultLogger.ErrorContext(ctx, msg, args...)
}
