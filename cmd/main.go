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

	profileRepo := repository.NewProfileRepository(db)
	profileHandler := handler.NewProfileHandler(profileRepo, db)

	movementRepo := repository.NewMovementRepository(db)
	movementHandler := handler.NewMovementHandler(movementRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	{
		api.GET("/profile", profileHandler.GetProfile)
		api.PUT("/profile", profileHandler.UpdateProfile)
		api.POST("/profile/avatar", profileHandler.UploadAvatar)

		api.GET("/movements/today", movementHandler.GetToday)
		api.POST("/movements/today", movementHandler.SaveToday)
		api.GET("/movements/total", movementHandler.GetTotal)
	}

	r.Run(":8080")
}
