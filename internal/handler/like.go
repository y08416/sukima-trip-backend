package handler

import (
	"errors"
	"log"
	"net/http"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	repo *repository.LikeRepository
}

func NewLikeHandler(repo *repository.LikeRepository) *LikeHandler {
	return &LikeHandler{repo: repo}
}

func (h *LikeHandler) Save(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")

	var req struct {
		PlaceName   string `json:"place_name" binding:"required"`
		PhotoURL    string `json:"photo_url"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	if err := h.repo.Save(userID, placeID, req.PlaceName, req.PhotoURL, req.Description); err != nil {
		if errors.Is(err, repository.ErrAlreadyLiked) {
			c.JSON(http.StatusConflict, gin.H{"error": "すでにいいね済みです"})
			return
		}
		log.Printf("[Like] Save失敗: userID=%s placeID=%s err=%v", userID, placeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "いいねに失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "いいねしました"})
}

func (h *LikeHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")

	if err := h.repo.Delete(userID, placeID); err != nil {
		if errors.Is(err, repository.ErrLikeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "いいねが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "いいね削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "いいねを外しました"})
}
