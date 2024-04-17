package log

import (
	"log/slog"
	"os"
)

func InitSlog() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
