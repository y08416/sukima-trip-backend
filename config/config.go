package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
	SupabaseJWTSecret      string
	GooglePlacesAPIKey     string
}

func Load() *Config {
	return &Config{
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:        os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseJWTSecret:      os.Getenv("SUPABASE_JWT_SECRET"),
		GooglePlacesAPIKey:     os.Getenv("GOOGLE_PLACES_API_KEY"),
	}
}

func (c *Config) Validate() error {
	required := map[string]string{
		"SUPABASE_URL":              c.SupabaseURL,
		"SUPABASE_SERVICE_ROLE_KEY": c.SupabaseServiceRoleKey,
		"SUPABASE_JWT_SECRET":       c.SupabaseJWTSecret,
		"GOOGLE_PLACES_API_KEY":     c.GooglePlacesAPIKey,
	}
	var missing []string
	for key, val := range required {
		if val == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("必須環境変数が未設定です: %s", strings.Join(missing, ", "))
	}
	return nil
}
