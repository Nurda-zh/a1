package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "github.com/Nurda-zh/a1/inventory-service/configs"
	delivery "github.com/Nurda-zh/a1/inventory-service/internal/delivery/http"
	"github.com/Nurda-zh/a1/inventory-service/internal/delivery/http/handler"
	"github.com/Nurda-zh/a1/inventory-service/internal/repository"
	"github.com/Nurda-zh/a1/inventory-service/internal/usecase"
)

func main() {
	cfg := config.LoadConfig()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(cfg.Database)

	repo := repository.NewProductRepository(db)
	uc := usecase.NewProductUsecase(repo)
	ph := handler.NewProductHandler(uc)

	r := gin.Default()
	delivery.NewRouter(r, ph)

	log.Println("Inventory service running on port " + cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
