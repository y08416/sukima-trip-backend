package handler

import (
	"net/http"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type VisitedPlaceHandler struct {
	repo *repository.VisitedPlaceRepository
}

func NewVisitedPlaceHandler(repo *repository.VisitedPlaceRepository) *VisitedPlaceHandler {
	return &VisitedPlaceHandler{repo: repo}
}

func (h *VisitedPlaceHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")

	places, err := h.repo.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "訪問地の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, places)
}

