package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/skantay/todo-list/internal/entity"

	"github.com/gin-gonic/gin"
)

type taskUsecase interface {
	Create(ctx context.Context, title string, activeAt entity.TaskDate) (string, error)
	List(ctx context.Context, status string) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task entity.Task) error
	MarkTaskDone(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type taskRoutes struct {
	taskUsecase taskUsecase
	log         *slog.Logger
}

func newTaskRoutes(router *gin.RouterGroup, taskUsecase taskUsecase, log *slog.Logger) {
	taskRoutes := taskRoutes{
		taskUsecase: taskUsecase,
		log:         log,
	}

	router.GET("/tasks", taskRoutes.list)

	router.POST("/tasks", taskRoutes.create)

	router.PUT("/tasks/:id", taskRoutes.update)

	router.DELETE("/tasks/:id", taskRoutes.delete)

	router.PUT("/tasks/:id/done", taskRoutes.markDone)
}

type requestTask struct {
	Title    string          `json:"title" binding:"required"`
	ActiveAt entity.TaskDate `json:"activeAt" binding:"required"`
}

type resp struct {
	ID string `json:"id"`
}

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
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	if len(tasks) == 0 {
		tasks = []entity.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

func getStatus(c *gin.Context) string {
	return c.Query("status")
}

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
		t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
		c.Status(http.StatusInternalServerError)

		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), req.Title, req.ActiveAt)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidTitle) || errors.Is(err, entity.ErrInvalidID) {
			t.log.Warn(http.StatusText(http.StatusBadRequest), "error", err)
			c.Status(http.StatusBadRequest)
		} else if errors.Is(err, entity.ErrAlreadyExists) {
			t.log.Warn(http.StatusText(http.StatusNotFound), "error", err)
			c.Status(http.StatusNotFound)
		} else {
			t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	response := resp{
		ID: id,
	}
	t.log.Debug(response.ID)

	c.JSON(http.StatusCreated, response)
}

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
		t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
		c.Status(http.StatusInternalServerError)

		return
	}

	id := c.Param("id")
	t.log.Debug(id)
	task := entity.NewTask(req.Title, req.ActiveAt)
	task.ID = id

	if err := t.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) || errors.Is(err, entity.ErrInvalidID) {
			t.log.Warn(http.StatusText(http.StatusBadRequest), "error", err)
			c.Status(http.StatusBadRequest)
		} else if errors.Is(err, entity.ErrTaskNotFound) {
			t.log.Warn(http.StatusText(http.StatusNotFound), "error", err)
			c.Status(http.StatusNotFound)
		} else {
			t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

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
			t.log.Warn(http.StatusText(http.StatusBadRequest), "error", err)
			c.Status(http.StatusBadRequest)
		} else if errors.Is(err, entity.ErrTaskNotFound) {
			t.log.Warn(http.StatusText(http.StatusNotFound), "error", err)
			c.Status(http.StatusNotFound)
		} else {
			t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
}

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
			t.log.Warn(http.StatusText(http.StatusNotFound), "error", err)
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, entity.ErrInvalidID) {
			t.log.Warn(http.StatusText(http.StatusBadRequest), "error", err)
			c.Status(http.StatusBadRequest)
		} else {
			t.log.Warn(http.StatusText(http.StatusInternalServerError), "error", err)
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
}
