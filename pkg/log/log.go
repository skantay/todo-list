package log

import (
	"log/slog"
	"os"
)

func InitSlog(opts *slog.HandlerOptions) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
