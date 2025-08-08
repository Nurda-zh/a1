package config

import "os"

type Config struct {
	MongoURI            string
	MongoDBName         string
	ListenPort          string
	InventoryServiceURL string
}

func Load() *Config {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017"
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "orders"
	}
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8002"
	}
	inv := os.Getenv("INVENTORY_URL")
	if inv == "" {
		inv = "http://inventory:8001/api"
	}
	return &Config{
		MongoURI:            uri,
		MongoDBName:         dbName,
		ListenPort:          port,
		InventoryServiceURL: inv,
	}
}
