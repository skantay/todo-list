package v1

import (
	"log/slog"

	"github.com/skantay/todo-list/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/skantay/todo-list/docs/api/v1"
)

// @title Todo List API
// @version 1
// @description API for managing todo list tasks
func Set(router *gin.Engine, usecase usecase.Usecase, log *slog.Logger) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	{
		newTaskRoutes(apiV1.Group("/todo-list"), usecase.TaskUsecase, log)
	}
}
