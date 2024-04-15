package main

import (
	"log/slog"

	"github.com/skantay/todo-list/internal/app"
)

func main() {
	slog.Info("program started")
	if err := app.Run(); err != nil {
		slog.Error("failed to start the program", "error", err)
	}
	slog.Info("program finished")
}
