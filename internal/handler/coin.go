package handler

import (
	"net/http"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type CoinHandler struct {
	repo *repository.CoinRepository
}

func NewCoinHandler(repo *repository.CoinRepository) *CoinHandler {
	return &CoinHandler{repo: repo}
}

func (h *CoinHandler) GetBalance(c *gin.Context) {
	userID := c.GetString("user_id")

	balance, err := h.repo.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "コイン残高取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, model.CoinResponse{Balance: balance})
}
