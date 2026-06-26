package repository

import (
	"encoding/json"
	"fmt"
	"strings"
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

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func (r *MovementRepository) GetToday(userID string) (*model.Movement, error) {
	today := time.Now().In(jst).Format("2006-01-02")

	data, _, err := r.client.From("movements").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("date", today).
		Single().
		Execute()
	if err != nil {
		if strings.Contains(err.Error(), "PGRST116") {
			return nil, nil
		}
		return nil, fmt.Errorf("移動距離取得失敗: %w", err)
	}

	var movement model.Movement
	if err := json.Unmarshal(data, &movement); err != nil {
		return nil, fmt.Errorf("データ変換失敗: %w", err)
	}
	return &movement, nil
}

func (r *MovementRepository) Save(userID string, req model.SaveMovementRequest) (*model.Movement, error) {
	result := r.client.Rpc("save_movement", "", map[string]interface{}{
		"p_user_id":                  userID,
		"p_real_distance_km":         req.RealDistanceKm,
		"p_used_virtual_distance_km": req.UsedVirtualDistanceKm,
	})

	trimmed := strings.TrimSpace(result)
	var rows []model.Movement
	if err := json.Unmarshal([]byte(trimmed), &rows); err != nil {
		var rpcErr rpcError
		if jsonErr := json.Unmarshal([]byte(trimmed), &rpcErr); jsonErr == nil && rpcErr.Code != "" {
			return nil, fmt.Errorf("移動距離保存失敗: %s", rpcErr.Message)
		}
		return nil, fmt.Errorf("移動距離保存失敗: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("移動距離保存失敗: レコードが返りませんでした")
	}
	return &rows[0], nil
}

func (r *MovementRepository) GetTotal(userID string) (float64, error) {
	result := r.client.Rpc("get_total_distance", "", map[string]interface{}{
		"p_user_id": userID,
	})

	trimmed := strings.TrimSpace(result)
	var total float64
	if err := json.Unmarshal([]byte(trimmed), &total); err == nil {
		return total, nil
	}

	var rpcErr rpcError
	if err := json.Unmarshal([]byte(trimmed), &rpcErr); err != nil {
		return 0, fmt.Errorf("合計距離取得失敗: レスポンス解析エラー: %w", err)
	}
	return 0, fmt.Errorf("合計距離取得失敗: %s", rpcErr.Message)
}
