package model

type Favorite struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PlaceID   string `json:"place_id"`
	PlaceName string `json:"place_name"`
	CreatedAt string `json:"created_at"`
}

type SaveFavoriteRequest struct {
	PlaceID   string `json:"place_id" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}
