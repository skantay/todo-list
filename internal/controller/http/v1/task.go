package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/skantay/todo-list/internal/entity"
	"github.com/skantay/todo-list/internal/usecase"

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
	fmt.Println(err)
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
	type req struct {
		Title    string          `json:"title" binding:"required"`
		ActiveAt entity.TaskDate `json:"activeAt" binding:"required"`
	}

	var request req

	if err := c.BindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest)

		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), request.Title, request.ActiveAt)
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

// UpdateTask(ctx context.Context, task entity.Task) error
func (t taskRoutes) update(c *gin.Context) {
	type req struct {
		Title    string          `json:"title" binding:"required"`
		ActiveAt entity.TaskDate `json:"activeAt" binding:"required"`
	}

	var request req

	if err := c.BindJSON(&request); err != nil {
		respondError(c, http.StatusBadRequest)

		return
	}

	id := c.Param("id")

	task := entity.NewTask(request.Title, request.ActiveAt)
	task.ID = id

	if err := t.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, nil)
}

// Delete(ctx context.Context, id string) error
func (t taskRoutes) delete(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.Delete(c.Request.Context(), id); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, nil)
}

// MarkTaskDone(ctx context.Context, id string) error
func (t taskRoutes) markDone(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.MarkTaskDone(c.Request.Context(), id); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, nil)
}
