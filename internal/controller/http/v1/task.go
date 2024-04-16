package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skantay/todo-list/internal/entity"

	"github.com/skantay/todo-list/internal/usecase"
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
}

func newTaskRoutes(router *gin.RouterGroup, taskUsecase taskUsecase) {
	taskRoutes := taskRoutes{
		taskUsecase: taskUsecase,
	}

	router.GET("/tasks", taskRoutes.list)

	router.POST("/tasks", taskRoutes.create)

	router.PUT("/tasks/:id", taskRoutes.update)

	router.DELETE("/tasks/:id", taskRoutes.delete)

	router.PUT("/tasks/:id/done", taskRoutes.markDone)
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

	c.JSON(http.StatusOK, tasks)
}

func getStatus(c *gin.Context) string {
	return c.Query("status")
}

// Create(ctx context.Context, title string, activeAt entity.TaskDate) (string, error)
func (t taskRoutes) create(c *gin.Context) {
	type request entity.Task

	if err := c.BindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest)

		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), request.Title, request.ActiveAt)
	if err != nil {
		respondError(c, http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, id)
}

func (t taskRoutes) update(c *gin.Context) {
	panic("implement me")
}

func (t taskRoutes) delete(c *gin.Context) {
	panic("implement me")
}

func (t taskRoutes) markDone(c *gin.Context) {
	panic("implement me")
}
