package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sukima-trip-backend/internal/model"
	"time"
)

const (
	placesAPIURL        = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	placesDetailsAPIURL = "https://maps.googleapis.com/maps/api/place/details/json"
	searchRadius = 5000

	coinDefault     = 10
	coinPopular     = 20
	coinVeryPopular = 30

	ratingsThresholdPopular     = 1000
	ratingsThresholdVeryPopular = 5000
)

func CalcCoinFromRatings(userRatingsTotal int) int {
	switch {
	case userRatingsTotal >= ratingsThresholdVeryPopular:
		return coinVeryPopular
	case userRatingsTotal >= ratingsThresholdPopular:
		return coinPopular
	default:
		return coinDefault
	}
}

type SpotRepository struct {
	apiKey     string
	httpClient *http.Client
}

func NewSpotRepository(apiKey string) *SpotRepository {
	return &SpotRepository{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type PlaceDetails struct {
	UserRatingsTotal int
	Description      string
	PhotoURL         string
}

func (r *SpotRepository) fetchPlaceDetails(ctx context.Context, placeID, fields string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s?place_id=%s&fields=%s&key=%s",
		placesDetailsAPIURL, placeID, fields, r.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト生成失敗: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Places Details API呼び出し失敗: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込み失敗: %w", err)
	}
	return body, nil
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
		Status  string `json:"status"`
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
	if result.Status != "OK" && result.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Places API エラー: %s", result.Status)
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

func (r *SpotRepository) GetUserRatingsTotal(ctx context.Context, placeID string) (int, error) {
	body, err := r.fetchPlaceDetails(ctx, placeID, "user_ratings_total")
	if err != nil {
		return 0, err
	}

	var result struct {
		Status string `json:"status"`
		Result struct {
			UserRatingsTotal int `json:"user_ratings_total"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("データ変換失敗: %w", err)
	}
	if result.Status != "OK" {
		return 0, fmt.Errorf("Places Details API エラー: %s (place_id=%s)", result.Status, placeID)
	}

	return result.Result.UserRatingsTotal, nil
}

func (r *SpotRepository) GetPlaceDetails(ctx context.Context, placeID string) (*PlaceDetails, error) {
	body, err := r.fetchPlaceDetails(ctx, placeID, "user_ratings_total,editorial_summary,photos")
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Result struct {
			UserRatingsTotal int `json:"user_ratings_total"`
			EditorialSummary struct {
				Overview string `json:"overview"`
			} `json:"editorial_summary"`
			Photos []struct {
				PhotoReference string `json:"photo_reference"`
			} `json:"photos"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	if result.Status != "OK" {
		return nil, fmt.Errorf("Places Details API エラー: %s (place_id=%s)", result.Status, placeID)
	}

	details := &PlaceDetails{
		UserRatingsTotal: result.Result.UserRatingsTotal,
		Description:      result.Result.EditorialSummary.Overview,
	}

	if len(result.Result.Photos) > 0 {
		details.PhotoURL = r.resolvePhotoURL(result.Result.Photos[0].PhotoReference)
	}

	return details, nil
}

func (r *SpotRepository) resolvePhotoURL(photoRef string) string {
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=800&photo_reference=%s&key=%s",
		photoRef, r.apiKey)
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
