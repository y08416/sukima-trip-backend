package handler

import (
	"net/http"
	"sukima-trip-backend/internal/model"
	"sukima-trip-backend/internal/repository"

	supa "github.com/supabase-community/supabase-go"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo *repository.AuthRepository
	db   *supa.Client
}

func NewAuthHandler(repo *repository.AuthRepository, db *supa.Client) *AuthHandler {
	return &AuthHandler{repo: repo, db: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.repo.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登録に失敗しました"})
		return
	}

	userID := session.User.ID.String()

	_, _, err = h.db.From("users").Insert(map[string]interface{}{
		"id":     userID,
		"name":   req.Name,
		"gender": req.Gender,
	}, false, "", "", "").Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "プロフィール作成に失敗しました"})
		return
	}

	_, _, err = h.db.From("coins").Insert(map[string]interface{}{
		"user_id": userID,
		"balance": 0,
	}, false, "", "", "").Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "コイン初期化に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken: session.AccessToken,
		UserID:      userID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.repo.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが間違っています"})
		return
	}

	c.JSON(http.StatusOK, model.AuthResponse{
		AccessToken: session.AccessToken,
		UserID:      session.User.ID.String(),
	})
}
