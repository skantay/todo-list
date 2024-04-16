package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, code int) {
	c.JSON(code, gin.H{"error": http.StatusText(code)})
}
