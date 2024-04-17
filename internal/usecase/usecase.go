package usecase

import (
	"log/slog"

	"github.com/skantay/todo-list/internal/repository"
)

type Usecase struct {
	TaskUsecase taskUsecase
}

func New(repository repository.Repository, log *slog.Logger) Usecase {
	return Usecase{
		TaskUsecase: newTaskUsecase(repository.TaskRepository, log),
	}
}
