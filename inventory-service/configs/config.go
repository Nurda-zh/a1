package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	Database   string
	ServerPort string
}

func LoadConfig() *Config {
	// Load .env file automatically
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	cfg := &Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:   getEnv("MONGO_DB", "inventory_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
	log.Println("Configuration loaded.")
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
