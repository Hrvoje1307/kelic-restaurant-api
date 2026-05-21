package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DSN              string
	GinMode          string
	JWTSecret        string
	AnthropicAPIKey  string
}

func Load() *Config {
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:            getEnv("PORT", "8080"),
		DSN:             getEnv("DATABASE_URL", ""),
		GinMode:         getEnv("GIN_MODE", "debug"),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
