package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/skantay/todo-list/internal/entity"
)

var (
	ErrInvalidTitle  = errors.New("invalid title")
	ErrInvalidStatus = errors.New("invalid status")
)

const (
	maxTitleLen = 200
)

type taskRepo interface {
	Create(ctx context.Context, task entity.Task) (string, error)
	List(ctx context.Context, status string, now time.Time) ([]entity.Task, error)
	Update(ctx context.Context, task entity.Task) error
	MarkDone(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type taskUsecase struct {
	repo taskRepo
	log  *slog.Logger
}

func newTaskUsecase(taskRepo taskRepo, log *slog.Logger) taskUsecase {
	return taskUsecase{
		repo: taskRepo,
		log:  log,
	}
}

func (t taskUsecase) Create(ctx context.Context, title string, activeAt entity.TaskDate) (string, error) {
	if utf8.RuneCountInString(title) > maxTitleLen {
		t.log.Error(ErrInvalidTitle.Error())
		t.log.Debug(ErrInvalidTitle.Error(), "title", title)
		return "", ErrInvalidTitle
	}

	task := entity.NewTask(title, activeAt)

	id, err := t.repo.Create(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	return id, nil
}

func (t taskUsecase) List(ctx context.Context, status string) ([]entity.Task, error) {
	if (status != entity.Active && status != entity.Done) && status != "" {
		return nil, ErrInvalidStatus
	} else if status == "" {
		status = entity.Active
	}

	tasks, err := t.repo.List(ctx, status, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	for i := range tasks {
		if tasks[i].ActiveAt.Time().Weekday() == time.Saturday || tasks[i].ActiveAt.Time().Weekday() == time.Sunday {
			tasks[i].Title = "ВЫХОДНОЙ - " + tasks[i].Title
		}
	}

	return tasks, nil
}

func (t taskUsecase) UpdateTask(ctx context.Context, task entity.Task) error {
	if err := t.repo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (t taskUsecase) MarkTaskDone(ctx context.Context, id string) error {
	if err := t.repo.MarkDone(ctx, id); err != nil {
		return fmt.Errorf("failed to mark task done: %w", err)
	}

	return nil
}

func (t taskUsecase) Delete(ctx context.Context, id string) error {
	if err := t.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
