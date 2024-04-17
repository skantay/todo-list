package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skantay/todo-list/internal/entity"
	"github.com/skantay/todo-list/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t taskRoutes) respondError(c *gin.Context, code int) {
	c.JSON(code, gin.H{"error": http.StatusText(code)})
}

func (t taskRoutes) handleCreateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrAlreadyExists) || errors.Is(err, entity.ErrTaskNotFound):
		t.respondError(c, http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidTitle) || errors.Is(err, primitive.ErrInvalidHex):
		t.respondError(c, http.StatusBadRequest)
	default:
		t.respondError(c, http.StatusInternalServerError)
	}
}
