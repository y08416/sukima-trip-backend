package handler

import (
	"net/http"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteHandler(repo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: repo}
}

func (h *FavoriteHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")

	favorites, err := h.repo.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *FavoriteHandler) Save(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.SaveFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	if err := h.repo.Save(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの保存に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "お気に入りに追加しました"})
}

func (h *FavoriteHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.repo.Delete(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "お気に入りから削除しました"})
}
