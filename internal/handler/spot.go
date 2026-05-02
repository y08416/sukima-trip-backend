package handler

import (
	"net/http"
	"strconv"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SpotHandler struct {
	spotRepo *repository.SpotRepository
	coinRepo *repository.CoinRepository
	visitedRepo *repository.VisitedPlaceRepository
}

func NewSpotHandler(
	spotRepo *repository.SpotRepository,
	coinRepo *repository.CoinRepository,
	visitedRepo *repository.VisitedPlaceRepository,
) *SpotHandler {
	return &SpotHandler{
		spotRepo:    spotRepo,
		coinRepo:    coinRepo,
		visitedRepo: visitedRepo,
	}
}

func (h *SpotHandler) GetSpots(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat と lng は必須です"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat の形式が不正です"})
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lng の形式が不正です"})
		return
	}

	spots, err := h.spotRepo.GetNearbySpots(lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "スポット取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, spots)
}

func (h *SpotHandler) Arrive(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")

	var req model.ArriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	if err := h.visitedRepo.Save(userID, model.SaveVisitedPlaceRequest{
		PlaceID:   placeID,
		PlaceName: req.PlaceName,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "訪問地の保存に失敗しました"})
		return
	}

	if err := h.coinRepo.AddCoin(userID, repository.CoinPerArrive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "コイン付与に失敗しました"})
		return
	}

	balance, err := h.coinRepo.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "コイン残高取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, model.ArriveResponse{
		Message:    "到着を記録しました",
		CoinEarned: repository.CoinPerArrive,
		Balance:    balance,
	})
}
