package handler

import (
	"net/http"
	"sukima-trip-backend/internal/model"
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

func (h *VisitedPlaceHandler) Save(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.SaveVisitedPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	if err := h.repo.Save(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "訪問地の保存に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "訪問地を記録しました"})
}
