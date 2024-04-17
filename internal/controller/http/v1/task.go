package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/skantay/todo-list/internal/entity"
	_ "github.com/skantay/todo-list/docs/api/v1"

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

// @Summary List tasks
// @Description Get a list of tasks based on the provided status
// @Param status query string false "Status of the tasks (active, done)"
// @Produce json
// @Success 200 {array} entity.Task
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks [get]
func (t taskRoutes) list(c *gin.Context) {
	status := getStatus(c)

	tasks, err := t.taskUsecase.List(c.Request.Context(), status)
	if err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

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
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks [post]
func (t taskRoutes) create(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	id, err := t.taskUsecase.Create(c.Request.Context(), req.Title, req.ActiveAt)
	if err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	type resp struct {
		ID string `json:"id"`
	}

	response := resp{
		ID: id,
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Update task
// @Description Update the details of an existing task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param requestTask body requestTask true "Task details"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/{id} [put]
func (t taskRoutes) update(c *gin.Context) {
	var req requestTask

	if err := c.BindJSON(&req); err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	id := c.Param("id")

	task := entity.NewTask(req.Title, req.ActiveAt)
	task.ID = id

	if err := t.taskUsecase.UpdateTask(c.Request.Context(), task); err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete task
// @Description Delete an existing task based on its ID
// @Param id path string true "Task ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/{id} [delete]
func (t taskRoutes) delete(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.Delete(c.Request.Context(), id); err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Mark task as done
// @Description Mark an existing task as done based on its ID
// @Param id path string true "Task ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/{id}/done [put]
func (t taskRoutes) markDone(c *gin.Context) {
	id := c.Param("id")

	if err := t.taskUsecase.MarkTaskDone(c.Request.Context(), id); err != nil {
		t.log.Error("", "error", err)
		t.handleCreateError(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
