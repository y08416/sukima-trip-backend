package repository

import (
	"fmt"

	supa "github.com/supabase-community/supabase-go"
)

type LikeRepository struct {
	client *supa.Client
}

func NewLikeRepository(client *supa.Client) *LikeRepository {
	return &LikeRepository{client: client}
}

func (r *LikeRepository) Save(userID, placeID, placeName string) error {
	_, _, err := r.client.From("spot_likes").
		Insert(map[string]interface{}{
			"user_id":    userID,
			"place_id":   placeID,
			"place_name": placeName,
		}, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("いいね保存失敗: %w", err)
	}
	return nil
}

func (r *LikeRepository) Delete(userID, placeID string) error {
	_, _, err := r.client.From("spot_likes").
		Delete("", "").
		Eq("user_id", userID).
		Eq("place_id", placeID).
		Execute()
	if err != nil {
		return fmt.Errorf("いいね削除失敗: %w", err)
	}
	return nil
}
