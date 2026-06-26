package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SpotHandler struct {
	spotRepo    *repository.SpotRepository
	visitedRepo *repository.VisitedPlaceRepository
}

func NewSpotHandler(
	spotRepo *repository.SpotRepository,
	visitedRepo *repository.VisitedPlaceRepository,
) *SpotHandler {
	return &SpotHandler{
		spotRepo:    spotRepo,
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

func (h *SpotHandler) GetNearestSpot(c *gin.Context) {
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

	spot, err := h.spotRepo.GetNearestSpot(lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "スポット取得に失敗しました"})
		return
	}
	if spot == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "近くにスポットが見つかりませんでした"})
		return
	}

	c.JSON(http.StatusOK, spot)
}

func (h *SpotHandler) Arrive(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")

	var req model.ArriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストが不正です"})
		return
	}

	details, err := h.spotRepo.GetPlaceDetails(c.Request.Context(), placeID)
	if err != nil {
		log.Printf("[Arrive] GetPlaceDetails失敗: placeID=%s err=%v", placeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "スポット情報の取得に失敗しました"})
		return
	}

	coinEarned := repository.CalcCoinFromRatings(details.UserRatingsTotal)

	balance, err := h.visitedRepo.SaveAndAddCoin(userID, placeID, req.PlaceName, coinEarned)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyVisited) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "このスポットにはすでに到着済みです"})
			return
		}
		log.Printf("[Arrive] SaveAndAddCoin失敗: userID=%s placeID=%s coinEarned=%d err=%v", userID, placeID, coinEarned, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "到着処理に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, model.ArriveResponse{
		Message:     "到着を記録しました",
		CoinEarned:  coinEarned,
		Balance:     balance,
		Description: details.Description,
		PhotoURL:    details.PhotoURL,
	})
}
