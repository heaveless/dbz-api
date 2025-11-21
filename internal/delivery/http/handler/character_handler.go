package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/heaveless/dbz-api/internal/domain/character"
)

type CharacterService interface {
	GetByName(ctx context.Context, name string) (*domain.CharacterDTO, error)
}

type CharacterHandler struct {
	service CharacterService
}

type getCharacterRequest struct {
	Name string `json:"name" binding:"required"`
}

func NewCharacterHandler(s CharacterService) *CharacterHandler {
	return &CharacterHandler{service: s}
}

func (h *CharacterHandler) GetOne(c *gin.Context) {
	var req getCharacterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The data submitted is invalid."})
		return
	}

	chr, err := h.service.GetByName(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusFound, gin.H{"data": chr})
}
