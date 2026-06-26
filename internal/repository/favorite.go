package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sukima-trip-backend/internal/model"

	"github.com/supabase-community/postgrest-go"
	supa "github.com/supabase-community/supabase-go"
)

var ErrAlreadyFavorited = errors.New("already favorited")
var ErrFavoriteNotFound = errors.New("favorite not found")

type FavoriteRepository struct {
	client *supa.Client
}

func NewFavoriteRepository(client *supa.Client) *FavoriteRepository {
	return &FavoriteRepository{client: client}
}

func (r *FavoriteRepository) GetAll(userID string) ([]model.Favorite, error) {
	data, _, err := r.client.From("spot_likes").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("お気に入り取得失敗: %w", err)
	}

	var likes []struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		PlaceID   string `json:"place_id"`
		PlaceName string `json:"place_name"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(data, &likes); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}

	favorites := make([]model.Favorite, len(likes))
	for i, l := range likes {
		favorites[i] = model.Favorite{
			ID:        l.ID,
			UserID:    l.UserID,
			PlaceID:   l.PlaceID,
			Name:      l.PlaceName,
			CreatedAt: l.CreatedAt,
		}
	}
	return favorites, nil
}

func (r *FavoriteRepository) Save(userID string, req model.SaveFavoriteRequest) error {
	_, _, err := r.client.From("favorites").
		Insert(map[string]interface{}{
			"user_id":   userID,
			"place_id":  req.PlaceID,
			"name":      req.PlaceName,
			"latitude":  req.Lat,
			"longitude": req.Lng,
		}, false, "", "", "").
		Execute()
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			return ErrAlreadyFavorited
		}
		return fmt.Errorf("お気に入り保存失敗: %w", err)
	}
	return nil
}

func (r *FavoriteRepository) Delete(userID, id string) error {
	data, _, err := r.client.From("spot_likes").
		Delete("", "").
		Eq("user_id", userID).
		Eq("id", id).
		Execute()
	if err != nil {
		return fmt.Errorf("お気に入り削除失敗: %w", err)
	}
	if string(data) == "[]" {
		return ErrFavoriteNotFound
	}
	return nil
}
