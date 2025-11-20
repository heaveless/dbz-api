package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heaveless/dbz-api/internal/delivery/http/handler"
)

func NewServer(handler *handler.CharacterHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/characters", handler.GetOne)

	return r
}
