package model

type Favorite struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"created_at"`
}

type SaveFavoriteRequest struct {
	PlaceID   string  `json:"place_id" binding:"required"`
	PlaceName string  `json:"place_name" binding:"required"`
	Lat       float64 `json:"latitude" binding:"required"`
	Lng       float64 `json:"longitude" binding:"required"`
}
