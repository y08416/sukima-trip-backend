package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	supa "github.com/supabase-community/supabase-go"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	repo *repository.ProfileRepository
	db   *supa.Client
}

func NewProfileHandler(repo *repository.ProfileRepository, db *supa.Client) *ProfileHandler {
	return &ProfileHandler{repo: repo, db: db}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := h.repo.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロフィール取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateProfile(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロフィール更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新しました"})
}

func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルが必要です"})
		return
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", userID, ext)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの読み込みに失敗しました"})
		return
	}
	defer src.Close()

	_, err = h.db.Storage.UploadFile("avatars", fileName, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像のアップロードに失敗しました"})
		return
	}

	avatarURL := h.db.Storage.GetPublicUrl("avatars", fileName).SignedURL

	if err := h.repo.UpdateAvatarURL(userID, avatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "画像URLの保存に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
