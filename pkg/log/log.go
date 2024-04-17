package log

import (
	"log/slog"
	"os"
)

func InitSlog() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

