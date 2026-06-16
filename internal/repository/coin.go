package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	supa "github.com/supabase-community/supabase-go"
)

type rpcError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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
	result := r.client.Rpc("increment_balance", "", map[string]interface{}{
		"p_user_id": userID,
		"p_amount":  amount,
	})

	// void関数の成功時はレスポンスが空またはnull
	trimmed := strings.TrimSpace(result)
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	var rpcErr rpcError
	if err := json.Unmarshal([]byte(result), &rpcErr); err != nil {
		return fmt.Errorf("コイン加算失敗: レスポンス解析エラー: %w", err)
	}
	if rpcErr.Code != "" {
		return fmt.Errorf("コイン加算失敗: %s", rpcErr.Message)
	}
	return nil
}
