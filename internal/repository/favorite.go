package repository

import (
	"encoding/json"
	"fmt"
	"sukima-trip-backend/internal/model"

	"github.com/supabase-community/postgrest-go"
	supa "github.com/supabase-community/supabase-go"
)

type FavoriteRepository struct {
	client *supa.Client
}

func NewFavoriteRepository(client *supa.Client) *FavoriteRepository {
	return &FavoriteRepository{client: client}
}

func (r *FavoriteRepository) GetAll(userID string) ([]model.Favorite, error) {
	data, _, err := r.client.From("favorites").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("お気に入り取得失敗: %w", err)
	}

	var favorites []model.Favorite
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	return favorites, nil
}

func (r *FavoriteRepository) Save(userID string, req model.SaveFavoriteRequest) error {
	_, _, err := r.client.From("favorites").
		Insert(map[string]interface{}{
			"user_id":    userID,
			"place_id":   req.PlaceID,
			"place_name": req.PlaceName,
		}, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("お気に入り保存失敗: %w", err)
	}
	return nil
}

func (r *FavoriteRepository) Delete(userID, id string) error {
	_, _, err := r.client.From("favorites").
		Delete("", "").
		Eq("user_id", userID).
		Eq("id", id).
		Execute()
	if err != nil {
		return fmt.Errorf("お気に入り削除失敗: %w", err)
	}
	return nil
}
