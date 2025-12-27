package logging

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	once    sync.Once
	verbose bool
)

type CLIHandler struct {
	stdout *slog.TextHandler
	stderr *slog.TextHandler
}

func NewCLIHandler(verboseMode bool) *CLIHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			if a.Key == slog.LevelKey && a.Value.String() == "INFO" {
				return slog.Attr{}
			}
			return a
		},
	}

	if verboseMode {
		opts.Level = slog.LevelDebug
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		}
	}

	return &CLIHandler{
		stdout: slog.NewTextHandler(os.Stdout, opts),
		stderr: slog.NewTextHandler(os.Stderr, opts),
	}
}

func (h *CLIHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.stdout.Enabled(ctx, level)
}

func (h *CLIHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.stderr.Handle(ctx, r)
	}
	return h.stdout.Handle(ctx, r)
}

func (h *CLIHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newStdout, ok := h.stdout.WithAttrs(attrs).(*slog.TextHandler)
	if !ok {
		return h
	}
	newStderr, ok := h.stderr.WithAttrs(attrs).(*slog.TextHandler)
	if !ok {
		return h
	}
	return &CLIHandler{
		stdout: newStdout,
		stderr: newStderr,
	}
}

func (h *CLIHandler) WithGroup(name string) slog.Handler {
	newStdout, ok := h.stdout.WithGroup(name).(*slog.TextHandler)
	if !ok {
		return h
	}
	newStderr, ok := h.stderr.WithGroup(name).(*slog.TextHandler)
	if !ok {
		return h
	}
	return &CLIHandler{
		stdout: newStdout,
		stderr: newStderr,
	}
}

func Init(verboseMode bool) {
	once.Do(func() {
		verbose = verboseMode
		handler := NewCLIHandler(verboseMode)
		slog.SetDefault(slog.New(handler))
	})
}

func IsVerbose() bool {
	return verbose
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
