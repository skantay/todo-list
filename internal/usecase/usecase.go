package usecase

import "github.com/skantay/todo-list/internal/repository"

type Usecase struct {
	TaskUsecase taskUsecase
}

func New(repository repository.Repository) Usecase {
	return Usecase{
		TaskUsecase: newTaskUsecase(repository.TaskRepository),
	}
}
