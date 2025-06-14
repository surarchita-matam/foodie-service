package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	MONGO_URI string
}

var config *Config

func init() {
	// Load .env file
	_ = godotenv.Load()

	config = &Config{
		JWTSecret: getEnvOrDefault("JWT_SECRET", "some-secret-key"),
		MONGO_URI: getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetConfig() *Config {
	return config
} 