package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/skantay/todo-list/internal/entity"
	"github.com/skantay/todo-list/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (t taskRoutes) list(c *gin.Context) {
	status := getStatus(c)

	tasks, err := t.taskUsecase.List(c.Request.Context(), status)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidStatus) {
			respondError(c, http.StatusBadRequest)

			return
		}
		respondError(c, http.StatusInternalServerError)

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

func (t taskRoutes) create(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {

		respondError(c, http.StatusBadRequest)

		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), req.Title, req.ActiveAt)
	if err != nil {
		if errors.Is(err, entity.ErrAlreadyExists) {
			respondError(c, http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrInvalidTitle) {
			respondError(c, http.StatusBadRequest)
		} else {
			respondError(c, http.StatusInternalServerError)
		}

		return
	}

	type resp struct {
		ID string `json:"id"`
	}

	response := resp{
		ID: id,
	}

	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusCreated, response)
}

func (t taskRoutes) update(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {

		respondError(c, http.StatusBadRequest)

		return
	}

	id := c.Param("id")

	task := entity.NewTask(req.Title, req.ActiveAt)
	task.ID = id

	if err := t.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, primitive.ErrInvalidHex) {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
	return
}

func (t taskRoutes) delete(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, primitive.ErrInvalidHex) {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
	return
}

func (t taskRoutes) markDone(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.MarkTaskDone(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrTaskNotFound) {
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, primitive.ErrInvalidHex) {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusNoContent)
	return
}
