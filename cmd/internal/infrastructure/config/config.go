package config

import "os"

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
	AppEnv    string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBUrl:     getEnv("DB_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),
		AppEnv:    getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
