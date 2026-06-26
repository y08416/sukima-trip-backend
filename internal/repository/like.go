package repository

import (
	"errors"
	"fmt"
	"strings"

	supa "github.com/supabase-community/supabase-go"
)

var ErrAlreadyLiked = errors.New("already liked")
var ErrLikeNotFound = errors.New("like not found")

type LikeRepository struct {
	client *supa.Client
}

func NewLikeRepository(client *supa.Client) *LikeRepository {
	return &LikeRepository{client: client}
}

func (r *LikeRepository) Save(userID, placeID, placeName, photoURL, description string) error {
	_, _, err := r.client.From("spot_likes").
		Insert(map[string]interface{}{
			"user_id":     userID,
			"place_id":    placeID,
			"place_name":  placeName,
			"photo_url":   photoURL,
			"description": description,
		}, false, "", "", "").
		Execute()
	if err != nil {
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			return ErrAlreadyLiked
		}
		return fmt.Errorf("いいね保存失敗: %w", err)
	}
	return nil
}

func (r *LikeRepository) Delete(userID, placeID string) error {
	data, _, err := r.client.From("spot_likes").
		Delete("", "").
		Eq("user_id", userID).
		Eq("place_id", placeID).
		Execute()
	if err != nil {
		return fmt.Errorf("いいね削除失敗: %w", err)
	}
	if string(data) == "[]" {
		return ErrLikeNotFound
	}
	return nil
}
