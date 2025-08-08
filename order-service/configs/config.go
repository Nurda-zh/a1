package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI            string
	Database            string
	ServerPort          string
	InventoryServiceURL string
}

func LoadConfig() *Config {
	cfg := &Config{
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:            getEnv("MONGO_DB", "orders_db"),
		ServerPort:          getEnv("SERVER_PORT", "8002"),
		InventoryServiceURL: getEnv("INVENTORY_URL", "http://localhost:8001/api"),
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
