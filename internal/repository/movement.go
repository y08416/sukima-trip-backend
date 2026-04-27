package repository

import (
	"encoding/json"
	"fmt"
	"sukima-trip-backend/internal/model"
	"time"

	supa "github.com/supabase-community/supabase-go"
)

type MovementRepository struct {
	client *supa.Client
}

func NewMovementRepository(client *supa.Client) *MovementRepository {
	return &MovementRepository{client: client}
}

func (r *MovementRepository) GetToday(userID string) (*model.Movement, error) {
	today := time.Now().Format("2006-01-02")

	data, _, err := r.client.From("movements").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("date", today).
		Single().
		Execute()
	if err != nil {
		return nil, nil
	}

	var movement model.Movement
	if err := json.Unmarshal(data, &movement); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	return &movement, nil
}

func (r *MovementRepository) Save(userID string, req model.SaveMovementRequest) error {
	today := time.Now().Format("2006-01-02")

	existing, err := r.GetToday(userID)
	if err != nil {
		return err
	}

	if existing != nil {
		_, _, err = r.client.From("movements").
			Update(map[string]interface{}{
				"real_distance_km":         existing.RealDistanceKm + req.RealDistanceKm,
				"used_virtual_distance_km": existing.UsedVirtualDistanceKm + req.UsedVirtualDistanceKm,
			}, "", "").
			Eq("id", existing.ID).
			Execute()
	} else {
		_, _, err = r.client.From("movements").
			Insert(map[string]interface{}{
				"user_id":                 userID,
				"date":                    today,
				"real_distance_km":        req.RealDistanceKm,
				"used_virtual_distance_km": req.UsedVirtualDistanceKm,
			}, false, "", "", "").
			Execute()
	}
	return err
}

func (r *MovementRepository) GetTotal(userID string) (float64, error) {
	data, _, err := r.client.From("movements").
		Select("real_distance_km", "", false).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		return 0, err
	}

	var rows []struct {
		RealDistanceKm float64 `json:"real_distance_km"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return 0, err
	}

	var total float64
	for _, row := range rows {
		total += row.RealDistanceKm
	}
	return total, nil
}
