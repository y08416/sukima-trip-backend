package main

import (
	"log"
	"sukima-trip-backend/config"
	"sukima-trip-backend/internal/handler"
	"sukima-trip-backend/internal/middleware"
	"sukima-trip-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	supa "github.com/supabase-community/supabase-go"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	db, err := supa.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey, nil)
	if err != nil {
		log.Fatalf("Supabaseクライアントの初期化に失敗しました: %v", err)
	}

	authRepo := repository.NewAuthRepository(db)
	authHandler := handler.NewAuthHandler(authRepo, db)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// 認証が必要なルートグループ
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	{
		api.GET("/health-auth", func(c *gin.Context) {
			userID := c.GetString("user_id")
			c.JSON(200, gin.H{"status": "ok", "user_id": userID})
		})
	}

	r.Run(":8080")
}
