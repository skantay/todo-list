package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/skantay/todo-list/internal/entity"
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

// @BasePath /api/v1/todo-list

// ListExample godoc
// @Summary list example
// @Schemes
// @Description list tasks
// @Tags example
// @Produce json
// @Success 200 {string} Jelloworld
// @Router /tasks [get]
func (t taskRoutes) list(c *gin.Context) {
	panic("implement me")
}

func (t taskRoutes) create(c *gin.Context) {
	panic("implement me")
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
