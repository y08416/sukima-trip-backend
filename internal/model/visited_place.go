package model

type VisitedPlace struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	PlaceID   string `json:"place_id"`
	PlaceName string `json:"place_name"`
	VisitedAt string `json:"visited_at"`
}

