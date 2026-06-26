package model

type Spot struct {
	PlaceID    string  `json:"place_id"`
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	DistanceKm float64 `json:"distance_km"`
}

type ArriveRequest struct {
	PlaceName string `json:"place_name" binding:"required"`
}

type ArriveResponse struct {
	Message     string `json:"message"`
	CoinEarned  int    `json:"coin_earned"`
	Balance     int    `json:"balance"`
	WikiSummary string `json:"wiki_summary"`
	PhotoURL    string `json:"photo_url"`
}

type NearestSpotResponse struct {
	PlaceID    string  `json:"place_id"`
	Name       string  `json:"name"`
	DistanceKm float64 `json:"distance_km"`
	Bearing    float64 `json:"bearing"`
}
