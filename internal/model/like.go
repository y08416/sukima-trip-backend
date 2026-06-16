package model

type SaveLikeRequest struct {
	PlaceID   string `json:"place_id" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}
