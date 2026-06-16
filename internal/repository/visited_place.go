package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sukima-trip-backend/internal/model"

	"github.com/supabase-community/postgrest-go"
	supa "github.com/supabase-community/supabase-go"
)

var ErrAlreadyVisited = errors.New("already visited")

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

func (r *VisitedPlaceRepository) Exists(userID, placeID string) (bool, error) {
	data, _, err := r.client.From("visited_places").
		Select("id", "", false).
		Eq("user_id", userID).
		Eq("place_id", placeID).
		Limit(1, "").
		Execute()
	if err != nil {
		return false, fmt.Errorf("訪問地確認失敗: %w", err)
	}

	var rows []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return false, fmt.Errorf("データ変換失敗: %w", err)
	}
	return len(rows) > 0, nil
}

func (r *VisitedPlaceRepository) SaveAndAddCoin(userID, placeID, placeName string, amount int) (int, error) {
	result := r.client.Rpc("arrive_spot", "", map[string]interface{}{
		"p_user_id":    userID,
		"p_place_id":   placeID,
		"p_place_name": placeName,
		"p_amount":     amount,
	})

	trimmed := strings.TrimSpace(result)
	if balance, err := strconv.Atoi(trimmed); err == nil {
		return balance, nil
	}

	var rpcErr rpcError
	if err := json.Unmarshal([]byte(result), &rpcErr); err != nil {
		return 0, fmt.Errorf("到着処理失敗: レスポンス解析エラー: %w", err)
	}
	if rpcErr.Code == "23505" || strings.Contains(rpcErr.Message, "duplicate key") {
		return 0, ErrAlreadyVisited
	}
	return 0, fmt.Errorf("到着処理失敗: %s", rpcErr.Message)
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
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			return ErrAlreadyVisited
		}
		return fmt.Errorf("訪問地保存失敗: %w", err)
	}
	return nil
}
