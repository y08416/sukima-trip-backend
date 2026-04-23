package config

import "os"

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
