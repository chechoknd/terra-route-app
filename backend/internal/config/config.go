package config

import "os"

type Config struct {
	AppEnv       string
	HTTPAddr     string
	LogLevel     string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn string
}

func Load() Config {
	return Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		HTTPAddr:     getEnv("HTTP_ADDR", ":8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://terraroute:terraroute_dev_password@localhost:5432/terraroute?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-for-local-development"),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "15m"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
