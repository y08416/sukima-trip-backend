package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sukima-trip-backend/internal/model"
)

const (
	placesAPIURL  = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	searchRadius  = 5000
	CoinPerArrive = 10
)

type SpotRepository struct {
	apiKey string
}

func NewSpotRepository(apiKey string) *SpotRepository {
	return &SpotRepository{apiKey: apiKey}
}

func (r *SpotRepository) GetNearbySpots(lat, lng float64) ([]model.Spot, error) {
	url := fmt.Sprintf("%s?location=%f,%f&radius=%d&type=tourist_attraction&language=ja&key=%s",
		placesAPIURL, lat, lng, searchRadius, r.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Places API呼び出し失敗: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}

	var result struct {
		Results []struct {
			PlaceID  string `json:"place_id"`
			Name     string `json:"name"`
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}

	spots := make([]model.Spot, 0, len(result.Results))
	for _, r := range result.Results {
		spots = append(spots, model.Spot{
			PlaceID:    r.PlaceID,
			Name:       r.Name,
			Lat:        r.Geometry.Location.Lat,
			Lng:        r.Geometry.Location.Lng,
			DistanceKm: calcDistance(lat, lng, r.Geometry.Location.Lat, r.Geometry.Location.Lng),
		})
	}
	return spots, nil
}

func calcDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
