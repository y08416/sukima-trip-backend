package handler

import (
	"errors"
	"net/http"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"
	"sync"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repo     *repository.FavoriteRepository
	spotRepo *repository.SpotRepository
}

func NewFavoriteHandler(repo *repository.FavoriteRepository, spotRepo *repository.SpotRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: repo, spotRepo: spotRepo}
}

func (h *FavoriteHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")

	favorites, err := h.repo.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの取得に失敗しました"})
		return
	}

	ctx := c.Request.Context()
	var wg sync.WaitGroup
	for i := range favorites {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			total, err := h.spotRepo.GetUserRatingsTotal(ctx, favorites[i].PlaceID)
			if err != nil {
				return
			}
			favorites[i].CoinAmount = repository.CalcCoinFromRatings(total)
		}(i)
	}
	wg.Wait()

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
		if errors.Is(err, repository.ErrAlreadyFavorited) {
			c.JSON(http.StatusConflict, gin.H{"error": "すでにお気に入り済みです"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの保存に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "お気に入りに追加しました"})
}

func (h *FavoriteHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.repo.Delete(userID, id); err != nil {
		if errors.Is(err, repository.ErrFavoriteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "お気に入りが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "お気に入りの削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "お気に入りから削除しました"})
}
