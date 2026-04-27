package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	supa "github.com/supabase-community/supabase-go"
)

func AuthMiddleware(client *supa.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが必要です"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := client.Auth.WithToken(tokenStr).GetUser()
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID.String())
		c.Next()
	}
}
