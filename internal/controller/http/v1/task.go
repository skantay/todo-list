// Пакет v1 предоставляет реализацию HTTP API для взаимодействия с задачами.
package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/skantay/todo-list/internal/entity"

	"github.com/gin-gonic/gin"
)

// taskUsecase определяет методы бизнес-логики для работы с задачами.
type taskUsecase interface {
	Create(ctx context.Context, title string, activeAt entity.TaskDate) (string, error)
	List(ctx context.Context, status string) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task entity.Task) error
	MarkTaskDone(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// taskRoutes определяет маршруты и их обработчики для задач.
type taskRoutes struct {
	taskUsecase taskUsecase  // Использование usecase-ов
	log         *slog.Logger // Логгер
}

// newTaskRoutes регистрирует эндпоинты для задач.
func newTaskRoutes(router *gin.RouterGroup, taskUsecase taskUsecase, log *slog.Logger) {
	taskRoutes := taskRoutes{
		taskUsecase: taskUsecase,
		log:         log,
	}

	router.GET("/tasks", taskRoutes.list) // Получение списка задач

	router.POST("/tasks", taskRoutes.create) // Создание задачи

	router.PUT("/tasks/:id", taskRoutes.update) // Обновление задачи

	router.DELETE("/tasks/:id", taskRoutes.delete) // Удаление задачи

	router.PUT("/tasks/:id/done", taskRoutes.markDone) // Пометить задачу как выполненную
}

// requestTask определяет структуру тела запроса для создания или обновления задачи.
type requestTask struct {
	Title    string          `json:"title" binding:"required"`
	ActiveAt entity.TaskDate `json:"activeAt" binding:"required"`
}

// resp определяет структуру ответа на успешное создание задачи.
type resp struct {
	ID string `json:"id"`
}

// list обрабатывает запрос на получение списка задач.

// @Summary List tasks
// @Description Get a list of tasks based on the provided status
// @Param status query string false "Status of the tasks (active, done)"
// @Produce json
// @Success 200 {array} entity.Task
// @Failure 400
// @Failure 500
// @Router /api/v1/todo-list/tasks [get]
func (t taskRoutes) list(c *gin.Context) {
	status := getStatus(c)

	tasks, err := t.taskUsecase.List(c.Request.Context(), status)
	if err != nil {
		t.log.Warn("", "error", err)
		if errors.Is(err, entity.ErrInvalidStatus) {
			t.respondStatus(c, http.StatusBadRequest, err)
		} else {
			t.respondStatus(c, http.StatusInternalServerError, err)
		}

		return
	}

	if len(tasks) == 0 {
		tasks = []entity.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

// getStatus извлекает статус из параметра запроса.
func getStatus(c *gin.Context) string {
	return c.Query("status")
}

// create обрабатывает запрос на создание новой задачи.

// @Summary Create task
// @Description Create a new task with the provided title and activeAt date
// @Accept json
// @Produce json
// @Param requestTask body requestTask true "Task details"
// @Success 201 {object} resp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/todo-list/tasks [post]
func (t taskRoutes) create(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {
		t.respondStatus(c, http.StatusInternalServerError, err)
		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), req.Title, req.ActiveAt)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidTitle) || errors.Is(err, entity.ErrInvalidID) {
			t.respondStatus(c, http.StatusBadRequest, err)
		} else if errors.Is(err, entity.ErrAlreadyExists) {
			t.respondStatus(c, http.StatusNotFound, err)
		} else {
			t.respondStatus(c, http.StatusInternalServerError, err)
		}

		return
	}

	response := resp{
		ID: id,
	}
	t.log.Debug(response.ID)

	c.JSON(http.StatusCreated, response)
}

// update обрабатывает запрос на обновление существующей задачи.

// @Summary Update task
// @Description Update the details of an existing task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param requestTask body requestTask true "Task details"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/todo-list/tasks/{id} [put]
func (t taskRoutes) update(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {
		t.respondStatus(c, http.StatusInternalServerError, err)

		return
	}

	id := c.Param("id")
	t.log.Debug(id)
	task := entity.NewTask(req.Title, req.ActiveAt)
	task.ID = id

	if err := t.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) || errors.Is(err, entity.ErrInvalidID) || errors.Is(err, entity.ErrInvalidTitle) {
			t.respondStatus(c, http.StatusBadRequest, err)
		} else if errors.Is(err, entity.ErrTaskNotFound) {
			t.respondStatus(c, http.StatusNotFound, err)
		} else {
			t.respondStatus(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

// delete обрабатывает запрос на удаление существующей задачи.

// @Summary Delete task
// @Description Delete an existing task based on its ID
// @Param id path string true "Task ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/todo-list/tasks/{id} [delete]
func (t taskRoutes) delete(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) || errors.Is(err, entity.ErrInvalidID) {
			t.respondStatus(c, http.StatusBadRequest, err)
		} else if errors.Is(err, entity.ErrTaskNotFound) {
			t.respondStatus(c, http.StatusNotFound, err)
		} else {
			t.respondStatus(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

// markDone обрабатывает запрос на пометку существующей задачи как выполненной.

// @Summary Mark task as done
// @Description Mark an existing task as done based on its ID
// @Param id path string true "Task ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/todo-list/tasks/{id}/done [put]
func (t taskRoutes) markDone(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.MarkTaskDone(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) || errors.Is(err, entity.ErrAlreadyExists) {
			t.respondStatus(c, http.StatusNotFound, err)
		} else if errors.Is(err, entity.ErrInvalidID) {
			t.respondStatus(c, http.StatusBadRequest, err)
		} else {
			t.respondStatus(c, http.StatusInternalServerError, err)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

func (t taskRoutes)respondStatus(c *gin.Context, code int, err error) {
	t.log.Warn(http.StatusText(code), "error", err)
	c.Status(code)
} 