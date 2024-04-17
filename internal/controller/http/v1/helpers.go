package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t taskRoutes) respondError(c *gin.Context, code int) {
	c.JSON(code, gin.H{"error": http.StatusText(code)})
}

func (t taskRoutes) handleCreateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrAlreadyExists) || errors.Is(err, entity.ErrTaskNotFound):
		respondError(c, StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidTitle) || errors.Is(err, primitive.ErrInvalidHex):
		respondError(c, StatusBadRequest)
	default:
		respondError(c, StatusInternalServerError)
	}
}
