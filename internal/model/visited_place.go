package model

type VisitedPlace struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PlaceID   string `json:"place_id"`
	PlaceName string `json:"place_name"`
	VisitedAt string `json:"visited_at"`
}

type SaveVisitedPlaceRequest struct {
	PlaceID   string `json:"place_id" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}
