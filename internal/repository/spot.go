package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	neturl "net/url"
	"sukima-trip-backend/internal/model"
	"time"
)

const (
	placesAPIURL        = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	placesDetailsAPIURL = "https://maps.googleapis.com/maps/api/place/details/json"
	searchRadius        = 5000
	ArriveRadiusKm      = 0.1
	CoinPerArrive       = 10
)

type SpotRepository struct {
	apiKey     string
	httpClient *http.Client
	wikiClient *http.Client
}

func NewSpotRepository(apiKey string) *SpotRepository {
	return &SpotRepository{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		wikiClient: &http.Client{Timeout: 3 * time.Second},
	}
}

func (r *SpotRepository) GetNearbySpots(lat, lng float64) ([]model.Spot, error) {
	url := fmt.Sprintf("%s?location=%f,%f&radius=%d&type=tourist_attraction&language=ja&key=%s",
		placesAPIURL, lat, lng, searchRadius, r.apiKey)

	resp, err := r.httpClient.Get(url)
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
			DistanceKm: CalcDistance(lat, lng, r.Geometry.Location.Lat, r.Geometry.Location.Lng),
		})
	}
	return spots, nil
}

func (r *SpotRepository) GetPlaceLocation(placeID string) (float64, float64, error) {
	endpoint := fmt.Sprintf("%s?place_id=%s&fields=geometry&key=%s",
		placesDetailsAPIURL, placeID, r.apiKey)

	resp, err := r.httpClient.Get(endpoint)
	if err != nil {
		return 0, 0, fmt.Errorf("Places Details API呼び出し失敗: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}

	var result struct {
		Result struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("データ変換失敗: %w", err)
	}

	lat := result.Result.Geometry.Location.Lat
	lng := result.Result.Geometry.Location.Lng
	if lat == 0 && lng == 0 {
		return 0, 0, fmt.Errorf("スポット座標の取得に失敗しました: place_id=%s", placeID)
	}

	return lat, lng, nil
}

func (r *SpotRepository) GetWikiInfo(name string) (string, string) {
	endpoint := fmt.Sprintf("https://ja.wikipedia.org/api/rest_v1/page/summary/%s",
		neturl.PathEscape(name))

	resp, err := r.wikiClient.Get(endpoint)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ""
	}

	var result struct {
		Extract   string `json:"extract"`
		Thumbnail struct {
			Source string `json:"source"`
		} `json:"thumbnail"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", ""
	}

	return result.Extract, result.Thumbnail.Source
}

func (r *SpotRepository) GetNearestSpot(lat, lng float64) (*model.NearestSpotResponse, error) {
	spots, err := r.GetNearbySpots(lat, lng)
	if err != nil {
		return nil, err
	}
	if len(spots) == 0 {
		return nil, nil
	}

	nearest := spots[0]
	for _, s := range spots[1:] {
		if s.DistanceKm < nearest.DistanceKm {
			nearest = s
		}
	}

	return &model.NearestSpotResponse{
		PlaceID:    nearest.PlaceID,
		Name:       nearest.Name,
		DistanceKm: nearest.DistanceKm,
		Bearing:    CalcBearing(lat, lng, nearest.Lat, nearest.Lng),
	}, nil
}

func CalcBearing(lat1, lng1, lat2, lng2 float64) float64 {
	dLng := (lng2 - lng1) * math.Pi / 180
	lat1R := lat1 * math.Pi / 180
	lat2R := lat2 * math.Pi / 180
	y := math.Sin(dLng) * math.Cos(lat2R)
	x := math.Cos(lat1R)*math.Sin(lat2R) - math.Sin(lat1R)*math.Cos(lat2R)*math.Cos(dLng)
	return math.Mod(math.Atan2(y, x)*180/math.Pi+360, 360)
}

func CalcDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
