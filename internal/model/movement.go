package model

type Movement struct {
	ID                    string  `json:"id"`
	UserID                string  `json:"user_id"`
	Date                  string  `json:"date"`
	RealDistanceKm        float64 `json:"real_distance_km"`
	UsedVirtualDistanceKm float64 `json:"used_virtual_distance_km"`
}

type MovementResponse struct {
	Date                  string  `json:"date"`
	RealDistanceKm        float64 `json:"real_distance_km"`
	VirtualDistanceKm     float64 `json:"virtual_distance_km"`
	UsedVirtualDistanceKm float64 `json:"used_virtual_distance_km"`
	RemainingDistanceKm   float64 `json:"remaining_distance_km"`
}

type SaveMovementRequest struct {
	RealDistanceKm        float64 `json:"real_distance_km" binding:"required"`
	UsedVirtualDistanceKm float64 `json:"used_virtual_distance_km"`
}

type TotalMovementResponse struct {
	TotalRealDistanceKm float64 `json:"total_real_distance_km"`
}
