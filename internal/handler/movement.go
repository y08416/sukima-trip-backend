package handler

import (
	"net/http"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type MovementHandler struct {
	repo *repository.MovementRepository
}

func NewMovementHandler(repo *repository.MovementRepository) *MovementHandler {
	return &MovementHandler{repo: repo}
}

func (h *MovementHandler) GetToday(c *gin.Context) {
	userID := c.GetString("user_id")

	movement, err := h.repo.GetToday(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移動距離取得に失敗しました"})
		return
	}

	if movement == nil {
		c.JSON(http.StatusOK, model.MovementResponse{
			RealDistanceKm:        0,
			VirtualDistanceKm:     0,
			UsedVirtualDistanceKm: 0,
			RemainingDistanceKm:   0,
		})
		return
	}

	virtualDistanceKm := movement.RealDistanceKm * 10
	remainingDistanceKm := virtualDistanceKm - movement.UsedVirtualDistanceKm

	c.JSON(http.StatusOK, model.MovementResponse{
		Date:                  movement.Date,
		RealDistanceKm:        movement.RealDistanceKm,
		VirtualDistanceKm:     virtualDistanceKm,
		UsedVirtualDistanceKm: movement.UsedVirtualDistanceKm,
		RemainingDistanceKm:   remainingDistanceKm,
	})
}

func (h *MovementHandler) SaveToday(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.SaveMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Save(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移動距離の保存に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存しました"})
}

func (h *MovementHandler) GetTotal(c *gin.Context) {
	userID := c.GetString("user_id")

	total, err := h.repo.GetTotal(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "総移動距離取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, model.TotalMovementResponse{
		TotalRealDistanceKm: total,
	})
}
