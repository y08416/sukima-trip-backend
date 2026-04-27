package repository

import (
	"encoding/json"
	"fmt"

	supa "github.com/supabase-community/supabase-go"
)

type CoinRepository struct {
	client *supa.Client
}

func NewCoinRepository(client *supa.Client) *CoinRepository {
	return &CoinRepository{client: client}
}

func (r *CoinRepository) GetBalance(userID string) (int, error) {
	data, _, err := r.client.From("coins").
		Select("balance", "", false).
		Eq("user_id", userID).
		Single().
		Execute()
	if err != nil {
		return 0, fmt.Errorf("コイン残高取得失敗: %w", err)
	}

	var result struct {
		Balance int `json:"balance"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("データ変換失敗: %w", err)
	}
	return result.Balance, nil
}

func (r *CoinRepository) AddCoin(userID string, amount int) error {
	balance, err := r.GetBalance(userID)
	if err != nil {
		return err
	}

	_, _, err = r.client.From("coins").
		Update(map[string]interface{}{
			"balance": balance + amount,
		}, "", "").
		Eq("user_id", userID).
		Execute()
	return err
}
