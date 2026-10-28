package config

import (
	"os"
)

type Config struct {
	Port       string
	DBFile     string
	JWTSecret  string
	CORSOrigin string
}

func LoadConfig() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBFile:     getEnv("DB_FILE", "app.db"),
		JWTSecret:  getEnv("JWT_SECRET", "golang-super-secret-key-change-in-production"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:4200"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
