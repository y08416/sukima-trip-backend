package model

type Profile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Gender    string `json:"gender"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}
