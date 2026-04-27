package repository

import (
	"encoding/json"
	"fmt"
	"sukima-trip-backend/internal/model"

	"github.com/supabase-community/postgrest-go"
	supa "github.com/supabase-community/supabase-go"
)

type VisitedPlaceRepository struct {
	client *supa.Client
}

func NewVisitedPlaceRepository(client *supa.Client) *VisitedPlaceRepository {
	return &VisitedPlaceRepository{client: client}
}

func (r *VisitedPlaceRepository) GetAll(userID string) ([]model.VisitedPlace, error) {
	data, _, err := r.client.From("visited_places").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("visited_at", &postgrest.OrderOpts{Ascending: false}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("訪問地取得失敗: %w", err)
	}

	var places []model.VisitedPlace
	if err := json.Unmarshal(data, &places); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	return places, nil
}

func (r *VisitedPlaceRepository) Save(userID string, req model.SaveVisitedPlaceRequest) error {
	_, _, err := r.client.From("visited_places").
		Insert(map[string]interface{}{
			"user_id":    userID,
			"place_id":   req.PlaceID,
			"place_name": req.PlaceName,
		}, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("訪問地保存失敗: %w", err)
	}
	return nil
}
