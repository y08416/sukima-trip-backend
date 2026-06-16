package repository

import (
	"encoding/json"
	"fmt"
	"sukima-trip-backend/internal/model"

	supa "github.com/supabase-community/supabase-go"
)

type ProfileRepository struct {
	client *supa.Client
}

func NewProfileRepository(client *supa.Client) *ProfileRepository {
	return &ProfileRepository{client: client}
}

func (r *ProfileRepository) GetProfile(userID string) (*model.Profile, error) {
	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("id", userID).
		Single().
		Execute()
	if err != nil {
		return nil, fmt.Errorf("プロフィール取得失敗: %w", err)
	}

	var profile model.Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(userID string, req model.UpdateProfileRequest) error {
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}

	_, _, err := r.client.From("users").
		Update(updates, "", "").
		Eq("id", userID).
		Execute()
	return err
}

func (r *ProfileRepository) UpdateAvatarURL(userID, avatarURL string) error {
	_, _, err := r.client.From("users").
		Update(map[string]interface{}{"avatar_url": avatarURL}, "", "").
		Eq("id", userID).
		Execute()
	return err
}
