package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/skantay/todo-list/internal/entity"
)

// Константы для usecase
const (
	maxTitleLen        = 200
	weekendTitlePrefix = "ВЫХОДНОЙ - "
	defaultStatus      = entity.Active
)

// taskRepo определяет интерфейс для repository
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

// Create создает новую задачу
func (t taskUsecase) Create(ctx context.Context, title string, activeAt entity.TaskDate) (string, error) {
	// Проверка максимальной длины заголовка
	if utf8.RuneCountInString(title) > maxTitleLen {
		return "", entity.ErrInvalidTitle
	}

	// Создание новой задачи
	task := entity.NewTask(title, activeAt)

	// Вызов метода репозитория для создания задачи
	id, err := t.repo.Create(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	return id, nil
}

// List возвращает список задач на основе указанного статуса
func (t taskUsecase) List(ctx context.Context, status string) ([]entity.Task, error) {
	// Проверка валидности статуса
	if status != entity.Active && status != entity.Done && status != "" {
		return nil, entity.ErrInvalidStatus
	}

	// Установка статуса по умолчанию, если не указан
	if status == "" {
		status = defaultStatus
	}

	// Получение списка задач из репозитория
	tasks, err := t.repo.List(ctx, status, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Добавление префикса к заголовку задачи, если она выпадает на выходные
	for i := range tasks {
		if tasks[i].ActiveAt.Time().Weekday() == time.Saturday || tasks[i].ActiveAt.Time().Weekday() == time.Sunday {
			tasks[i].Title = weekendTitlePrefix + tasks[i].Title
		}
	}

	return tasks, nil
}

// UpdateTask обновляет информацию о задаче
func (t taskUsecase) UpdateTask(ctx context.Context, task entity.Task) error {
	// Проверка максимальной длины заголовка
	if utf8.RuneCountInString(task.Title) > maxTitleLen {
		return entity.ErrInvalidTitle
	}

	// Вызов метода репозитория для обновления задачи
	if err := t.repo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// MarkTaskDone помечает задачу как выполненную
func (t taskUsecase) MarkTaskDone(ctx context.Context, id string) error {
	// Вызов метода репозитория для пометки задачи как выполненной
	if err := t.repo.MarkDone(ctx, id); err != nil {
		return fmt.Errorf("failed to mark task done: %w", err)
	}

	return nil
}

// Delete удаляет задачу
func (t taskUsecase) Delete(ctx context.Context, id string) error {
	// Вызов метода репозитория для удаления задачи
	if err := t.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
